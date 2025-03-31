package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "memcached-checker",
	Long:              "A tool to verify whether it fulfills the functionality of memcached.",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	SilenceUsage:      true,
}

func init() {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		rootCmd.Version = bi.Main.Version
	}
}

func Execute() error {
	return rootCmd.Execute()
}
