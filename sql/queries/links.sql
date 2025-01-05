-- name: CreateShortLink :one
INSERT INTO shortly(id, short_link, long_link)
VALUES($1, $2, $3)
RETURNING *;



-- name: ListLinks :many
SELECT * FROM shortly
ORDER BY created_at DESC;