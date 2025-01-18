-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows(id, user_id, feed_id, created_at, updated_at)
  VALUES($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT
  inserted_feed_follow.*
  ,f.name AS feed_name
  ,u.name as user_name
FROM inserted_feed_follow
JOIN feeds f ON inserted_feed_follow.feed_id = f.id
JOIN users u ON inserted_feed_follow.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT
  f.name feedname
  ,u.name username
FROM feed_follows AS ff
JOIN users AS u ON ff.user_id = u.id
JOIN feeds AS f ON ff.feed_id = f.id
WHERE u.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
USING feeds, users 
WHERE feed_follows.feed_id = feeds.id
  AND feed_follows.user_id = $1
  AND feeds.url = $2;
