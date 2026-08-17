package commands

import (
	"Lattice/internal/logger"
	"Lattice/internal/storage"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// depsCmd represents the deps command
var dependentsCmd = &cobra.Command{
	Use:   "dependents [symbol]",
	Short: "Show dependencies of a symbol",
	Args:  cobra.ExactArgs(1),
	RunE:  getDependents,
}

func init() {
	rootCmd.AddCommand(dependentsCmd)

	// Here you will define your flags and configuration settings.

}

func getDependents(cmd *cobra.Command, args []string) error {
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
	rows, err := store.GetDependents(cmd.Context(), symbol)
	if err != nil {
		return err
	}
	for _, row := range rows {
		fmt.Println(row.SourceSymbol)
	}

	return nil
}
