-- name: GetFileByPath :one
SELECT *
FROM files
WHERE path = ?;

-- name: UpsertFile :one
INSERT INTO files (path, language, hash, parsed_at)
VALUES (?, ?, ?, ?)
ON CONFLICT(path) DO UPDATE SET
    hash = excluded.hash,
    parsed_at = excluded.parsed_at
RETURNING id;
