package refresh

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

const lformat = "=== %s ==="

// Logger struct
type Logger struct {
	log *log.Logger
}

// NewLogger func
func NewLogger(c *Configuration) *Logger {
	color.NoColor = !c.EnableColors
	if runtime.GOOS == "windows" {
		color.NoColor = true
	}
	if len(c.LogName) == 0 {
		c.LogName = "refresh"
	}
	var w io.Writer = c.Stdout
	if w == nil {
		w = os.Stdout
	}
	return &Logger{
		log: log.New(w, fmt.Sprintf("%s: ", c.LogName), log.LstdFlags),
	}
}

// Success func
func (l *Logger) Success(msg interface{}, args ...interface{}) {
	l.log.Print(color.GreenString(fmt.Sprintf(lformat, msg), args...))
}

// Error func
func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Print(color.RedString(fmt.Sprintf(lformat, msg), args...))
}

// Print func
func (l *Logger) Print(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(lformat, msg), args...)
}

// LogLocation func
var LogLocation = func() string {
	dir, _ := homedir.Dir()
	dir, _ = homedir.Expand(dir)
	dir = path.Join(dir, ".refresh")
	os.MkdirAll(dir, 0755)
	return dir
}

// ErrorLogPath func
var ErrorLogPath = func() string {
	return path.Join(LogLocation(), ID()+".err")
}
