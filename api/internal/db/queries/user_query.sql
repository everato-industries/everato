-- =========================================================
-- User related queries
-- =========================================================

-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email, password
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2,
    last_name = $3,
    email = $4,
    password = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
