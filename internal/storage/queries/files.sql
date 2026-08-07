-- name: InsertFile :exec
INSERT INTO files (
    path,
    language,
    hash,
    parsed_at
)
VALUES (?, ?, ?, ?);

-- name: GetFileByPath :one
SELECT *
FROM files
WHERE path = ?;
