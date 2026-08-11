-- name: InsertSymbol :exec
INSERT INTO symbols (id, file_id, name, receiver, kind, start_line, end_line, exported, scope)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: DeleteSymbolsByFile :exec
DELETE FROM symbols WHERE file_id = ?;
