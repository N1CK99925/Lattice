CREATE TABLE files (
    id INTEGER PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL,
    hash TEXT NOT NULL,
    parsed_at DATETIME NOT NULL
);

CREATE TABLE symbols (
    id TEXT PRIMARY KEY,
    file_id INTEGER NOT NULL REFERENCES files(id),
    name TEXT NOT NULL,
    receiver TEXT,
    kind TEXT NOT NULL,
    start_line INTEGER NOT NULL,
    end_line INTEGER NOT NULL,
    exported BOOLEAN NOT NULL
);

CREATE TABLE edges (
    source_symbol TEXT NOT NULL REFERENCES symbols(id),
    target_symbol TEXT,
    target_external TEXT,
    kind TEXT NOT NULL
);

CREATE INDEX idx_symbols_name ON symbols(name);
CREATE INDEX idx_edges_source ON edges(source_symbol);
CREATE INDEX idx_edges_target ON edges(target_symbol);
