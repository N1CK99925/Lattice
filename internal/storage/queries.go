package storage

import (
	"Lattice/internal/storage/generated"
	"context"
)

func (s *Store) getDependencies(ctx context.Context, symbol string) ([]generated.GetDependenciesRow, error) {
	return s.Queries.GetDependencies(ctx, symbol)
}
