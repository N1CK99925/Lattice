/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"Lattice/internal/storage"

	"github.com/spf13/cobra"
)

// callersCmd represents the callers command
var callersCmd = &cobra.Command{
	Use:   "callers [symbol]",
	Short: "Prints callers of a symbol",
	Args:  cobra.ExactArgs(1),
	RunE:  RunCallers,
}

func init() {
	rootCmd.AddCommand(callersCmd)

}

func RunCallers(cmd *cobra.Command, args []string) error {
	store, err := storage.New("./Lattice.db")
	if err != nil {
		return err
	}
	defer store.Close()
	symbol := args[0]
	rows, err := store.GetCallers(cmd.Context(), symbol)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if row.TargetExternal.Valid {
			println(row.TargetExternal.String)
		}
	}

	return nil
}
