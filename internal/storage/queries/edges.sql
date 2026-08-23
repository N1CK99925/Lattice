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

-- name: GetCallers :many
SELECT source_symbol, target_external, kind
FROM edges
WHERE target_symbol = ?
AND kind = 'calls';
-- name: Deadcode :many
SELECT s.id, s.name
FROM symbols s
WHERE s.scope = 'package'
  AND s.exported = 1
  AND NOT EXISTS (
      SELECT 1
      FROM edges e
      WHERE e.target_symbol = s.id
        AND e.kind IN ('calls', 'type_use')
  );
