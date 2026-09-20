package graph

import (
	"Lattice/internal/logger"
	"Lattice/internal/storage"
	"context"
)

// This creates a graph from GetAllSymbols and GetAllEdges and makes a in memory Graph
func Load(ctx context.Context, store *storage.Store) (*Graph, error) {
	g := &Graph{
		Nodes:    make(map[string]Node),
		outgoing: make(map[string][]Edge),
		incoming: make(map[string][]Edge),
	}
	symbols, err := store.Queries.GetAllSymbols(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	for _, symbol := range symbols {
		g.Nodes[symbol] = Node{
			ID: symbol,
		}
	}
	edges, err := store.Queries.GetAllEdges(ctx)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	for _, row := range edges {
		if !row.TargetSymbol.Valid {
			continue
		}

		edge := Edge{
			Source: row.SourceSymbol,
			Target: row.TargetSymbol.String,
			Kind:   row.Kind,
		}

		g.Edges = append(g.Edges, edge)

		g.outgoing[edge.Source] =
			append(g.outgoing[edge.Source], edge)

		g.incoming[edge.Target] =
			append(g.incoming[edge.Target], edge)
	}
	return g, nil
}
