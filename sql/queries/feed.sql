-- name: CreateFeed :one
INSERT INTO "feed"
    (user_id, name, url)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: GetFeedListSortedByCreation :many
SELECT sqlc.embed(feed), sqlc.embed(u) FROM "feed"
    JOIN "user" u on "feed".user_id = u.id
    ORDER BY "feed"."created_at";

-- name: GetFeedById :one
SELECT sqlc.embed(feed), sqlc.embed(u) FROM "feed"
    JOIN "user" u ON feed.user_id = u."id"
    WHERE "feed"."id" = $1;

-- name: CheckFeedExistenceByName :one
SELECT EXISTS (SELECT 1 FROM "feed" WHERE "name" = $1);