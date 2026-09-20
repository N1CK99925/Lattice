package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "lattice update",
	Short: "Updates the lattice graph",
}

func UpdateGraph() {
	err := UpdateCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
