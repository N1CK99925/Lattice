package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deadcodeCmd = &cobra.Command{
	Use:   "deadcode",
	Short: "Find unused code in a Go repository",
	Long:  `Scans a Go codebase and identifies functions, types, and variables that are defined but never referenced.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deadcode called")
	},
}

func init() {
	rootCmd.AddCommand(deadcodeCmd)
}
