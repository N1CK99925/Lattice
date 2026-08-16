package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lattice",
	Short: "Code analysis tool for Go repositories",
	Long:  `Lattice parses Go source files, builds a symbol table, and resolves cross-file references for dead code detection and dependency analysis.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
