-- name: CreateShortLink :one
INSERT INTO shortly(id, short_link, long_link)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetLongLink :one
SELECT * FROM shortly
WHERE short_link = $1
LIMIT 1;

-- name: DeleteLink :exec
DELETE FROM shortly WHERE short_link = $1;

-- name: ListLinks :many
SELECT * FROM shortly
ORDER BY created_at DESC;