// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: links.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createShortLink = `-- name: CreateShortLink :one
INSERT INTO shortly(id, short_link, long_link)
VALUES($1, $2, $3)
RETURNING id, short_link, long_link, created_at, updated_at
`

type CreateShortLinkParams struct {
	ID        uuid.UUID
	ShortLink string
	LongLink  string
}

func (q *Queries) CreateShortLink(ctx context.Context, arg CreateShortLinkParams) (Shortly, error) {
	row := q.db.QueryRowContext(ctx, createShortLink, arg.ID, arg.ShortLink, arg.LongLink)
	var i Shortly
	err := row.Scan(
		&i.ID,
		&i.ShortLink,
		&i.LongLink,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listLinks = `-- name: ListLinks :many
SELECT id, short_link, long_link, created_at, updated_at FROM shortly
ORDER BY created_at DESC
`

func (q *Queries) ListLinks(ctx context.Context) ([]Shortly, error) {
	rows, err := q.db.QueryContext(ctx, listLinks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shortly
	for rows.Next() {
		var i Shortly
		if err := rows.Scan(
			&i.ID,
			&i.ShortLink,
			&i.LongLink,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
