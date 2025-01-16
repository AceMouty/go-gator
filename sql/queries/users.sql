-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * 
FROM users
WHERE users.name = $1;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1 
    FROM users 
    WHERE name = $1
) AS exists;

