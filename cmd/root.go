package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "mychainctl",
	Short: "mychainctl is a minimal chainctl-style CLI",
}

func Execute() error {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table or json")

	return rootCmd.Execute()
}
