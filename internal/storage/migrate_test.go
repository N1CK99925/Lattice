package storage

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func schemaVersion(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	var v int64
	if err := db.QueryRow(`SELECT version FROM schema_migrations LIMIT 1`).Scan(&v); err != nil {
		t.Fatalf("read schema_migrations: %v", err)
	}
	return v
}

func hasColumn(t *testing.T, db *sql.DB, table, column string) bool {
	t.Helper()
	rows, err := db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatal(err)
		}
		if name == column {
			return true
		}
	}
	return false
}

func TestMigrateFresh(t *testing.T) {
	db := openTestDB(t)
	if err := migrate(db); err != nil {
		t.Fatal(err)
	}

	if v := schemaVersion(t, db); v != 2 {
		t.Fatalf("expected schema version 2, got %d", v)
	}
	for _, col := range []string{"id", "path", "language", "hash", "parsed_at"} {
		if !hasColumn(t, db, "files", col) {
			t.Fatalf("files table missing column %q", col)
		}
	}
	for _, col := range []string{"id", "file_id", "name", "kind", "start_line", "end_line", "exported", "scope"} {
		if !hasColumn(t, db, "symbols", col) {
			t.Fatalf("symbols table missing column %q", col)
		}
	}

	// scope defaults to 'package' for existing-style rows
	if _, err := db.Exec(`INSERT INTO files (path, language, hash, parsed_at) VALUES ('a.go', 'go', 'h', datetime('now'))`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO symbols (id, file_id, name, kind, start_line, end_line, exported) VALUES ('f::a.go::F', 1, 'F', 'function', 1, 1, 1)`); err != nil {
		t.Fatal(err)
	}
	var scope string
	if err := db.QueryRow(`SELECT scope FROM symbols WHERE id = 'f::a.go::F'`).Scan(&scope); err != nil {
		t.Fatal(err)
	}
	if scope != "package" {
		t.Fatalf("expected scope 'package', got %q", scope)
	}

	// edge kinds are constrained to the four roadmap kinds
	if _, err := db.Exec(`INSERT INTO edges (source_symbol, target_symbol, kind) VALUES ('f::a.go::F', 'x', 'type_use')`); err != nil {
		t.Fatalf("type_use edge rejected: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO edges (source_symbol, target_symbol, kind) VALUES ('f::a.go::F', 'x', 'bogus')`); err == nil {
		t.Fatal("expected CHECK constraint to reject unknown edge kind")
	}
}

func TestMigrateFromV1(t *testing.T) {
	// simulate a database created by golang-migrate at version 1
	db := openTestDB(t)
	body, err := migrationsFS.ReadFile("migrations/000001_init.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(body)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE schema_migrations (version BIGINT NOT NULL PRIMARY KEY, dirty BOOLEAN NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO schema_migrations (version, dirty) VALUES (1, 0)`); err != nil {
		t.Fatal(err)
	}

	if err := migrate(db); err != nil {
		t.Fatal(err)
	}
	if v := schemaVersion(t, db); v != 2 {
		t.Fatalf("expected schema version 2, got %d", v)
	}
	if !hasColumn(t, db, "symbols", "scope") {
		t.Fatal("scope column not added on upgrade")
	}
	// 000001 must not have re-run: its tables still exist with data intact
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM files`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected 0 files after upgrade, got %d", n)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	db := openTestDB(t)
	for i := 0; i < 2; i++ {
		if err := migrate(db); err != nil {
			t.Fatal(err)
		}
	}
	if v := schemaVersion(t, db); v != 2 {
		t.Fatalf("expected schema version 2, got %d", v)
	}
}

func TestMigrateDirty(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`CREATE TABLE schema_migrations (version BIGINT NOT NULL PRIMARY KEY, dirty BOOLEAN NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO schema_migrations (version, dirty) VALUES (1, 1)`); err != nil {
		t.Fatal(err)
	}
	err := migrate(db)
	if err == nil || !strings.Contains(err.Error(), "dirty") {
		t.Fatalf("expected dirty-state error, got %v", err)
	}
}
