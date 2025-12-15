package handlers

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "casd",
	Short: "casd CLI",
	Long:  "CLI tool for deduplicated content addressing system.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return rootCmd.Execute()
}
