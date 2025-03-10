-- name: CreatePost :one
INSERT INTO posts(
  id
  ,title
  ,url
  ,description 
  ,feed_id
  ,published_at
  ,created_at
  ,updated_at
)
VALUES(
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
  p.id
  ,p.title
  ,f.name AS feed_name
  ,p.url
  ,p.description 
  ,p.feed_id
  ,p.published_at
  ,p.created_at
  ,p.updated_at
FROM posts AS p
JOIN feeds AS f ON p.feed_id = f.id
JOIN feed_follows AS ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;
