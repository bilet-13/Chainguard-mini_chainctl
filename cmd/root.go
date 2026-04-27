package cmd

import (
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "mychainctl",
	Short: "mychainctl is a minimal chainctl-style CLI",
}

var setupOnce sync.Once

func setupRoot() {
	setupOnce.Do(func() {
		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true

		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)

		rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table or json")
	})
}

func Execute() error {
	setupRoot()
	return rootCmd.Execute()
}

func ExecuteWithArgs(args []string) error {
	setupRoot()
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
