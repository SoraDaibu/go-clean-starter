-- name: CreateItem :one
INSERT INTO items (id, type_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetItem :one
SELECT * FROM items WHERE id = $1 LIMIT 1;

-- name: ListItems :many
SELECT * FROM items ORDER BY created_at DESC;

-- name: UpdateItem :one
UPDATE items
SET type_id = $2
WHERE id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;
