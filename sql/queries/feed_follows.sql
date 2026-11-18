-- name: CreateFeedFollow :one
WITH insert_result AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING *
)
SELECT ff.id, ff.created_at, ff.updated_at, u.name, f.name FROM insert_result ff
JOIN feeds f
    ON f.url = ff.feed_id
JOIN users u
    ON u.id = ff.user_id;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, u.name, f.name FROM feed_follows ff
JOIN users u
    ON ff.user_id = u.id
JOIN feeds f
    ON f.url = ff.feed_id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows ff where ff.user_id = $1 and ff.feed_id = $2
RETURNING *;
