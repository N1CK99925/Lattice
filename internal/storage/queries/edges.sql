-- name: InsertEdge :exec
INSERT INTO edges (source_symbol, target_symbol, target_external, kind)
VALUES (?, ?, ?, ?);

-- name: DeleteEdgesByFile :exec
DELETE FROM edges
WHERE source_symbol IN (SELECT s.id FROM symbols s WHERE s.file_id = ?)
   OR target_symbol IN (SELECT s2.id FROM symbols s2 WHERE s2.file_id = ?);

-- name: GetDependencies :many
SELECT target_symbol, target_external, kind
FROM edges
WHERE source_symbol = ?;

-- name: GetDependents :many
SELECT source_symbol, target_external, kind
FROM edges
WHERE target_symbol = ?;

