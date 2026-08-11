package storage

import (
	"Lattice/internal/logger"
	"Lattice/internal/storage/generated"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB      *sql.DB
	Queries *generated.Queries
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA temp_store = MEMORY;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, err
		}
	}

	// Apply pending embedded migrations in order before anything touches the
	// schema; a fresh database ends up fully migrated, an existing one only
	// gets the migrations it is missing.
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	db.SetMaxOpenConns(1)

	queries := generated.New(db)

	logger.Log.Info("DB is open")
	return &Store{
		DB:      db,
		Queries: queries,
	}, nil
}

func (s *Store) Close() error {
	logger.Log.Info("DB is closed")
	return s.DB.Close()
}
