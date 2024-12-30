-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UserByEmail :one
SELECT * 
FROM users
WHERE email = $1;

-- name: GetUserById :one 
SELECT *
FROM users 
WHERE id = $1; 

-- name: UpdateUser :exec
UPDATE users
SET email = $1 , hashed_password = $2
WHERE users.id = $3;

-- name: UpgradeToChirpyRed :exec 
UPDATE users 
SET is_chirpy_red = TRUE
WHERE users.id = $1;

