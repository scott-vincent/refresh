package cmd

import (
	"fmt"
	"os"

	"github.com/scott-vincent/refresh/refresh"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generates a default configuration file for you.",
	Run: func(cmd *cobra.Command, args []string) {
		c := refresh.Configuration{
			AppRoot:            ".",
			IgnoredFolders:     []string{"vendor", "log", "logs", "tmp", "node_modules", "bin", "templates"},
			IncludedExtensions: []string{".go"},
			BuildTargetPath:    "",
			BuildPath:          os.TempDir(),
			BuildDelay:         200,
			BinaryName:         "refresh-build",
			GoOrVue:            "go",
			CommandFlags:       []string{},
			CommandEnv:         []string{},
			EnableColors:       true,
		}

		if cfgFile == "" {
			cfgFile = "refresh.yml"
		}
		err := c.Dump(cfgFile)
		if err != nil {
			fmt.Println("Failed to create file:", cfgFile, "error:", err)
		} else {
			fmt.Println("Created file:", cfgFile)
		}
	},
}
