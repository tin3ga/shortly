// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Shortly struct {
	ID        uuid.UUID
	ShortLink string
	LongLink  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
