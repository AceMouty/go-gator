// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
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
RETURNING id, title, url, description, feed_id, published_at, created_at, updated_at
`

type CreatePostParams struct {
	ID          uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
	FeedID      uuid.UUID
	PublishedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.FeedID,
		arg.PublishedAt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.FeedID,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :many
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
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type GetPostsForUserRow struct {
	ID          uuid.UUID
	Title       string
	FeedName    string
	Url         string
	Description sql.NullString
	FeedID      uuid.UUID
	PublishedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]GetPostsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsForUserRow
	for rows.Next() {
		var i GetPostsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.FeedName,
			&i.Url,
			&i.Description,
			&i.FeedID,
			&i.PublishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
