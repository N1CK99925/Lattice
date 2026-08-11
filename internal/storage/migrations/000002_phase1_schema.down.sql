-- Reverse 000002: drop the scope column and restore the unconstrained edges
-- table. SQLite cannot drop columns via ALTER TABLE, so both tables are
-- rebuilt in their 000001 shape.

CREATE TABLE symbols_v1 (
    id TEXT PRIMARY KEY,
    file_id INTEGER NOT NULL REFERENCES files(id),
    name TEXT NOT NULL,
    receiver TEXT,
    kind TEXT NOT NULL,
    start_line INTEGER NOT NULL,
    end_line INTEGER NOT NULL,
    exported BOOLEAN NOT NULL
);

INSERT INTO symbols_v1 (id, file_id, name, receiver, kind, start_line, end_line, exported)
SELECT id, file_id, name, receiver, kind, start_line, end_line, exported FROM symbols;

DROP TABLE symbols;
ALTER TABLE symbols_v1 RENAME TO symbols;

CREATE TABLE edges_v1 (
    source_symbol  TEXT NOT NULL REFERENCES symbols(id),
    target_symbol  TEXT,
    target_external TEXT,
    kind           TEXT NOT NULL
);

INSERT INTO edges_v1 (source_symbol, target_symbol, target_external, kind)
SELECT source_symbol, target_symbol, target_external, kind FROM edges;

DROP TABLE edges;
ALTER TABLE edges_v1 RENAME TO edges;

CREATE INDEX idx_edges_source ON edges(source_symbol);
CREATE INDEX idx_edges_target ON edges(target_symbol);
