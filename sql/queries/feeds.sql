-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name FROM feeds f
JOIN users u
    on f.user_id = u.id;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE url = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
