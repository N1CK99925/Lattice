package storage

import (
	"Lattice/internal/logger"
	"Lattice/internal/storage/generated"
	"database/sql"

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
