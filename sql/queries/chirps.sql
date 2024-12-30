-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetOneChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1;

-- name: DeleteOneChirp :exec
DELETE FROM chirps 
WHERE id = $1 AND user_id = $2; 

