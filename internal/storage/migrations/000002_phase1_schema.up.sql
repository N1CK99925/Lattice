-- Phase 1 schema additions (ROADMAP point 3):
--   1. scope column on symbols: 'package' (declared at package level) vs
--      'local' (function-local declarations, captured by later parser passes).
--      Defaults to 'package'; the current extractor only emits package-level
--      declarations, so existing rows need no backfill.
--   2. Explicit edge kinds: calls, imports, defines, type_use. TypeRefs
--      captured by the parser (currently unused) wire in as type_use edges.
--      SQLite cannot add a CHECK constraint via ALTER TABLE, so edges is
--      rebuilt with the constraint in place.

ALTER TABLE symbols ADD COLUMN scope TEXT NOT NULL DEFAULT 'package';

CREATE TABLE edges_v2 (
    source_symbol   TEXT NOT NULL REFERENCES symbols(id),
    target_symbol   TEXT,
    target_external TEXT,
    kind            TEXT NOT NULL CHECK (kind IN ('calls', 'imports', 'defines', 'type_use'))
);

INSERT INTO edges_v2 (source_symbol, target_symbol, target_external, kind)
SELECT source_symbol, target_symbol, target_external, kind FROM edges;

DROP TABLE edges;
ALTER TABLE edges_v2 RENAME TO edges;

CREATE INDEX idx_edges_source ON edges(source_symbol);
CREATE INDEX idx_edges_target ON edges(target_symbol);
