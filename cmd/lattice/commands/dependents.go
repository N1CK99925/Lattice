package commands

import (
	"Lattice/internal/graph"
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
	dependentsCmd.Flags().IntP("depth", "d", 15, "Depth to traverse")
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
	g, err := graph.Load(cmd.Context(), store)
	if err != nil {
		return err
	}
	symbol := args[0]

	depth, err := cmd.Flags().GetInt("depth")
	fmt.Printf("depth: %d\n", depth)
	if err != nil {
		return err
	}
	nodes, err := g.Dependents(symbol, depth)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		fmt.Println(node.ID)
	}
	return nil
}
