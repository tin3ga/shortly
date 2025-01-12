-- name: CreateShortLink :one
INSERT INTO shortly(id, user_id, short_link, long_link)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetLongLink :one
SELECT * FROM shortly
WHERE short_link = $1
LIMIT 1;

-- name: DeleteLink :exec
DELETE FROM shortly WHERE short_link = $1;

-- name: GetLinks :many
SELECT * FROM shortly
ORDER BY created_at DESC;


-- name: IncrementClickCount :exec
UPDATE shortly
SET click_count = click_count + 1, updated_at = NOW()
WHERE short_link = $1;


-- name: GetUserLinks :many
SELECT * FROM shortly
WHERE user_id = $1
ORDER BY created_at DESC;