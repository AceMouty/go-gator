-- name: CreateFeed :one
INSERT INTO feeds(id, name, url, user_id, created_at, updated_at)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT 
  f.name
  ,f.url
  ,u.name username
FROM feeds f
JOIN users u ON f.user_id = u.id;
