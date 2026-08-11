package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrate applies pending embedded migrations in version order, in a
// transaction each. The file layout (NNNNNN_name.up.sql) and the
// schema_migrations version table match golang-migrate's conventions, so the
// Makefile's db-* targets keep working against the same files and a database
// previously migrated with the `migrate` CLI is picked up seamlessly.
func migrate(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT NOT NULL PRIMARY KEY,
		dirty BOOLEAN NOT NULL
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	var version int64
	var dirty bool
	err := db.QueryRow(`SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&version, &dirty)
	switch {
	case err == sql.ErrNoRows:
		// fresh database
	case err != nil:
		return fmt.Errorf("read schema_migrations: %w", err)
	case dirty:
		return fmt.Errorf("database is dirty at migration %d; resolve manually", version)
	}

	ms, err := listMigrations()
	if err != nil {
		return err
	}

	for _, m := range ms {
		if m.version <= version {
			continue
		}
		body, err := migrationsFS.ReadFile("migrations/" + m.file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", m.file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(body)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", m.file, err)
		}
		// single-row version tracking, same shape as golang-migrate
		if _, err := tx.Exec(`DELETE FROM schema_migrations`); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", m.file, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version, dirty) VALUES (?, 0)`, m.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", m.file, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", m.file, err)
		}
	}

	return nil
}

type migration struct {
	version int64
	file    string
}

func listMigrations() ([]migration, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("list migrations: %w", err)
	}

	var ms []migration
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		ver, ok := parseMigrationVersion(name)
		if !ok {
			return nil, fmt.Errorf("migration file %q does not start with a numeric version", name)
		}
		ms = append(ms, migration{version: ver, file: name})
	}

	sort.Slice(ms, func(i, j int) bool { return ms[i].version < ms[j].version })
	for i := 1; i < len(ms); i++ {
		if ms[i].version == ms[i-1].version {
			return nil, fmt.Errorf("duplicate migration version %d", ms[i].version)
		}
	}
	return ms, nil
}

// parseMigrationVersion extracts the leading version from golang-migrate style
// filenames: NNNNNN_name.up.sql
func parseMigrationVersion(name string) (int64, bool) {
	i := strings.IndexByte(name, '_')
	if i <= 0 {
		return 0, false
	}
	v, err := strconv.ParseInt(name[:i], 10, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
