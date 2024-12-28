-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: DeleteRefreshTokens :exec
DELETE FROM refresh_tokens;

-- name: GetToken :one
SELECT * 
FROM refresh_tokens
WHERE refresh_tokens.token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
WHERE refresh_tokens.token = $1;
