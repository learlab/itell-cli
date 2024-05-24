package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "itell-cli",
	Short: "Create and Manage ITELL projects",
	Long:  `Create and Manage ITELL projects. "itell-cli create <dest> -t <template>" with create a new project with template. "itell-cli update" synchronize the project with upstream template.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
}
