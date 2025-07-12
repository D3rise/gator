-- name: CreateFeed :one
INSERT INTO "feed"
    (user_id, name, url)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: GetFeedListSortedByCreation :many
SELECT * FROM "feed" ORDER BY "feed"."created_at";

-- name: GetFeedById :one
SELECT * FROM "feed" WHERE id = $1;

-- name: CheckFeedExistenceByName :one
SELECT EXISTS (SELECT 1 FROM "feed" WHERE "name" = $1);