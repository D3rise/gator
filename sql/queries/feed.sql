-- name: CreateFeed :one
INSERT INTO "feed"
    (user_id, name, url)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: SetFeedFetchedAtToNowById :exec
UPDATE "feed" SET last_fetched_at = NOW() WHERE id = $1;

-- name: GetFeedListSortedByCreation :many
SELECT sqlc.embed(feed), sqlc.embed(u) FROM "feed"
    JOIN "user" u on "feed".user_id = u.id
    ORDER BY "feed"."created_at";

-- name: GetFeedListSortedByLastFetchedAt :many
SELECT * FROM "feed"
    ORDER BY "feed"."last_fetched_at" NULLS FIRST;

-- name: GetOldestFeedByUpdatedAt :one
SELECT * FROM feed ORDER BY updated_at LIMIT 1;

-- name: GetFeedById :one
SELECT sqlc.embed(feed), sqlc.embed(u) FROM "feed"
    JOIN "user" u ON feed.user_id = u."id"
    WHERE "feed"."id" = $1;

-- name: GetFeedByUrl :one
SELECT sqlc.embed(feed), sqlc.embed(u) FROM "feed"
    JOIN "user" u ON feed.user_id = u."id"
    WHERE "feed"."url" = $1;

-- name: CheckFeedExistenceByName :one
SELECT EXISTS (SELECT 1 FROM "feed" WHERE "name" = $1);