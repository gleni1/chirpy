// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING id, created_at, updated_at, email
`

type CreateUserParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}