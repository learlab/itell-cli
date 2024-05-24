package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Synchronize with template",
	Long:  `Synchronize the project with upstream template, leaving config files untouched`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update is not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
