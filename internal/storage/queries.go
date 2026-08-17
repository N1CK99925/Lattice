package storage

import (
	"Lattice/internal/storage/generated"
	"context"
	"database/sql"
)

func (s *Store) GetDependencies(ctx context.Context, symbol string) ([]generated.GetDependenciesRow, error) {
	return s.Queries.GetDependencies(ctx, symbol)
}

func (s *Store) GetDependents(ctx context.Context, symbol string) ([]generated.GetDependentsRow, error) {
	return s.Queries.GetDependents(ctx, sql.NullString{
		String: symbol,
		Valid:  true,
	})
}

func (s *Store) GetCallers(ctx context.Context, symbol string) ([]generated.GetCallersRow, error) {
	return s.Queries.GetCallers(ctx, sql.NullString{
		String: symbol,
		Valid:  true,
	})
}
