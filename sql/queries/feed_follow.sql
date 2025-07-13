-- name: CreateNewFeedFollow :one
INSERT INTO feed_follow (user_id, feed_id) VALUES ($1, $2) RETURNING *;

-- name: CheckFeedFollowExistence :one
SELECT EXISTS (SELECT 1 FROM feed_follow WHERE user_id = $1 AND feed_id = $2);

-- name: GetFeedFollowListByUserId :many
SELECT sqlc.embed(feed_follow), sqlc.embed(feed)
    FROM feed_follow
    JOIN public.feed on feed_follow.feed_id = feed.id
    WHERE feed_follow.user_id = $1
    ORDER BY feed_follow.created_at;

-- name: GetFeedFollowListByFeedId :many
SELECT sqlc.embed(feed_follow), sqlc.embed(feed)
    FROM feed_follow
    JOIN public.feed on feed_follow.feed_id = feed.id
    WHERE feed_id = $1 ORDER BY feed_follow.created_at;

-- name: GetFeedFollowListPaginated :many
SELECT * FROM feed_follow;