-- name: CreateUser :one
INSERT INTO
    users (email, password_hash)
VALUES ($1, $2) RETURNING id,
    email,
    created_at,
    updated_at;

-- name: FindUser :one
SELECT id, email, created_at, updated_at
FROM users
WHERE
    email = $1
LIMIT 1;
