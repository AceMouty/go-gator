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

-- name: GetFeed :one
SELECT *
FROM feeds f
WHERE f.url = $1;

-- name: FeedExists :one
SELECT EXISTS (
    SELECT 1 
    FROM feeds f
    WHERE f.url = $1
) AS exists;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
  last_fetched_at = NOW(), updated_at = NOW()
WHERE
  id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, name, url, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

