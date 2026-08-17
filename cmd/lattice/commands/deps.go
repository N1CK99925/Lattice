package commands

import (
	"Lattice/internal/logger"
	"Lattice/internal/storage"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// depsCmd represents the deps command
var depsCmd = &cobra.Command{
	Use:   "deps [symbol]",
	Short: "Show dependencies of a symbol",
	Args:  cobra.ExactArgs(1),
	RunE:  getDeps,
}

func init() {
	rootCmd.AddCommand(depsCmd)

	// Here you will define your flags and configuration settings.

}

func getDeps(cmd *cobra.Command, args []string) error {
	logger.Init(logger.Config{
		Level: slog.LevelInfo,
		Json:  false,
	})

	store, err := storage.New("./Lattice.db")
	if err != nil {
		logger.Log.Error("Error opening database", "err", err)
	}
	defer store.Close()
	symbol := args[0]
	rows, err := store.GetDependencies(cmd.Context(), symbol)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if row.TargetSymbol.Valid {
			fmt.Println(row.TargetSymbol.String)
		}

		if row.TargetExternal.Valid {
			fmt.Println(row.TargetExternal.String)
		}
	}

	return nil
}
