package refresh

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Manager struct
type Manager struct {
	*Configuration
	ID         string
	Logger     *Logger
	Restart    chan bool
	cancelFunc context.CancelFunc
	context    context.Context
	gil        *sync.Once
}

// New func
func New(c *Configuration) *Manager {
	return NewWithContext(context.Background(), c)
}

// NewWithContext func
func NewWithContext(ctx context.Context, c *Configuration) *Manager {
	ctx, cancelFunc := context.WithCancel(ctx)
	m := &Manager{
		Configuration: c,
		ID:            ID(),
		Logger:        NewLogger(c),
		Restart:       make(chan bool),
		cancelFunc:    cancelFunc,
		context:       ctx,
		gil:           &sync.Once{},
	}
	return m
}

// Start func
func (r *Manager) Start() error {
	w := NewWatcher(r)
	w.Start()
	go r.build(fsnotify.Event{Name: ":start:"})
	if !r.Debug {
		go func() {
			for {
				select {
				case event := <-w.Events():
					if event.Op != fsnotify.Chmod {
						go r.build(event)
					}
					w.Remove(event.Name)
					w.Add(event.Name)
				case <-r.context.Done():
					break
				}
			}
		}()
	}
	go func() {
		for {
			select {
			case err := <-w.Errors():
				r.Logger.Error(err)
			case <-r.context.Done():
				break
			}
		}
	}()
	r.runner()
	return nil
}

func (r *Manager) build(event fsnotify.Event) {
	r.gil.Do(func() {
		defer func() {
			r.gil = &sync.Once{}
		}()
		r.buildTransaction(func() error {
			// time.Sleep(r.BuildDelay * time.Millisecond)

			now := time.Now()
			r.Logger.Print("Rebuild on: %s", event.Name)

			var cmd *exec.Cmd
			if r.GoOrVue == "vue" {
				cmd = exec.CommandContext(r.context, "npm", "run", "build")
			} else {
				args := []string{"build", "-v"}
				args = append(args, r.BuildFlags...)
				args = append(args, "-o", r.FullBuildPath(), r.BuildTargetPath)
				cmd = exec.CommandContext(r.context, "go", args...)
			}

			err := r.runAndListen(cmd)
			if err != nil {
				if strings.Contains(err.Error(), "no buildable Go source files") {
					r.cancelFunc()
					log.Fatal(err)
				}
				return err
			}

			tt := time.Since(now)
			r.Logger.Success("Building Completed (PID: %d) (Time: %s)", cmd.Process.Pid, tt)
			r.Restart <- true
			return nil
		})
	})
}

func (r *Manager) buildTransaction(fn func() error) {
	lpath := ErrorLogPath()
	err := fn()
	if err != nil {
		f, _ := os.Create(lpath)
		fmt.Fprint(f, err)
		r.Logger.Error("Error!")
		r.Logger.Error(err)
	} else {
		os.Remove(lpath)
	}
}
