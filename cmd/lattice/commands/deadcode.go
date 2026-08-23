package commands

import (
	"Lattice/internal/logger"
	"Lattice/internal/storage"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var deadcodeCmd = &cobra.Command{
	Use:   "deadcode",
	Short: "Find unused code in a Go repository",
	RunE:  CallDeadCode,
}

func init() {
	rootCmd.AddCommand(deadcodeCmd)
}
func CallDeadCode(cmd *cobra.Command, args []string) error {
	logger.Init(logger.Config{
		Level: slog.LevelInfo,
		Json:  false,
	})
	store, err := storage.New("./Lattice.db")
	if err != nil {
		return err
	}
	defer store.Close()

	rows, err := store.Deadcode(cmd.Context())
	if err != nil {
		return err
	}
	for _, row := range rows {
		fmt.Println(row.Name)
	}
	return nil
}
