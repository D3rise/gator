-- name: UpsertPostOnUrl :exec
INSERT INTO "post"
    (title, url, description, published_at, feed_id)
    VALUES
    ($1, $2, $3, $4, $5)
    ON CONFLICT (url) DO NOTHING;

-- name: GetPostsByUserFeedFollows :many
SELECT post.* FROM "post"
    JOIN public.feed f on f.id = post.feed_id
    JOIN public.feed_follow ff on f.id = ff.feed_id
    JOIN public."user" u on u.id = ff.user_id
    WHERE u.id = $1;