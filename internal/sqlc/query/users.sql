-- name: CreateUser :one
INSERT INTO users (id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
