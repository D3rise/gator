-- name: CreateUser :one
INSERT INTO "user"
    (name, updated_at)
VALUES
    ($1, $2) RETURNING *;

-- name: GetUserByName :one
SELECT * FROM "user" WHERE name = $1;

-- name: CheckUserExistenceByName :one
SELECT EXISTS (SELECT 1 FROM "user" WHERE name = $1);