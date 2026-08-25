package commands

import (
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:  "export",
	RunE: Export,
}

func init() {
	rootCmd.AddCommand(exportCmd)

}

func Export(cmd *cobra.Command, args []string) error {
	return nil
}
