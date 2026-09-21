package storage

import (
	"Lattice/internal/storage/generated"
	"context"
	"database/sql"
)

func (s *Store) GetCallers(ctx context.Context, symbol string) ([]generated.GetCallersRow, error) {
	return s.Queries.GetCallers(ctx, sql.NullString{
		String: symbol,
		Valid:  true,
	})
}

func (s *Store) Deadcode(ctx context.Context) ([]generated.DeadcodeRow, error) {
	return s.Queries.Deadcode(ctx)
}
