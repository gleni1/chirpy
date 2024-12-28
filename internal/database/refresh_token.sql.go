// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_token.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken,
		arg.Token,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.ExpiresAt,
		arg.RevokedAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const deleteRefreshTokens = `-- name: DeleteRefreshTokens :exec
DELETE FROM refresh_tokens
`

func (q *Queries) DeleteRefreshTokens(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteRefreshTokens)
	return err
}

const getToken = `-- name: GetToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at 
FROM refresh_tokens
WHERE refresh_tokens.token = $1
`

func (q *Queries) GetToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const revokeToken = `-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
WHERE refresh_tokens.token = $1
`

func (q *Queries) RevokeToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeToken, token)
	return err
}
