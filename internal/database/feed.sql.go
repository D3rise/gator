// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feed.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const checkFeedExistenceByName = `-- name: CheckFeedExistenceByName :one
SELECT EXISTS (SELECT 1 FROM "feed" WHERE "name" = $1)
`

func (q *Queries) CheckFeedExistenceByName(ctx context.Context, name string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkFeedExistenceByName, name)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createFeed = `-- name: CreateFeed :one
INSERT INTO "feed"
    (user_id, name, url)
VALUES
    ($1, $2, $3)
RETURNING id, user_id, name, url, created_at, updated_at, last_fetched_at
`

type CreateFeedParams struct {
	UserID uuid.UUID
	Name   string
	Url    string
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed, arg.UserID, arg.Name, arg.Url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeedById = `-- name: GetFeedById :one
SELECT feed.id, feed.user_id, feed.name, feed.url, feed.created_at, feed.updated_at, feed.last_fetched_at, u.id, u.name, u.created_at, u.updated_at FROM "feed"
    JOIN "user" u ON feed.user_id = u."id"
    WHERE "feed"."id" = $1
`

type GetFeedByIdRow struct {
	Feed Feed
	User User
}

func (q *Queries) GetFeedById(ctx context.Context, id uuid.UUID) (GetFeedByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getFeedById, id)
	var i GetFeedByIdRow
	err := row.Scan(
		&i.Feed.ID,
		&i.Feed.UserID,
		&i.Feed.Name,
		&i.Feed.Url,
		&i.Feed.CreatedAt,
		&i.Feed.UpdatedAt,
		&i.Feed.LastFetchedAt,
		&i.User.ID,
		&i.User.Name,
		&i.User.CreatedAt,
		&i.User.UpdatedAt,
	)
	return i, err
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT feed.id, feed.user_id, feed.name, feed.url, feed.created_at, feed.updated_at, feed.last_fetched_at, u.id, u.name, u.created_at, u.updated_at FROM "feed"
    JOIN "user" u ON feed.user_id = u."id"
    WHERE "feed"."url" = $1
`

type GetFeedByUrlRow struct {
	Feed Feed
	User User
}

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (GetFeedByUrlRow, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i GetFeedByUrlRow
	err := row.Scan(
		&i.Feed.ID,
		&i.Feed.UserID,
		&i.Feed.Name,
		&i.Feed.Url,
		&i.Feed.CreatedAt,
		&i.Feed.UpdatedAt,
		&i.Feed.LastFetchedAt,
		&i.User.ID,
		&i.User.Name,
		&i.User.CreatedAt,
		&i.User.UpdatedAt,
	)
	return i, err
}

const getFeedListSortedByCreation = `-- name: GetFeedListSortedByCreation :many
SELECT feed.id, feed.user_id, feed.name, feed.url, feed.created_at, feed.updated_at, feed.last_fetched_at, u.id, u.name, u.created_at, u.updated_at FROM "feed"
    JOIN "user" u on "feed".user_id = u.id
    ORDER BY "feed"."created_at"
`

type GetFeedListSortedByCreationRow struct {
	Feed Feed
	User User
}

func (q *Queries) GetFeedListSortedByCreation(ctx context.Context) ([]GetFeedListSortedByCreationRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedListSortedByCreation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedListSortedByCreationRow
	for rows.Next() {
		var i GetFeedListSortedByCreationRow
		if err := rows.Scan(
			&i.Feed.ID,
			&i.Feed.UserID,
			&i.Feed.Name,
			&i.Feed.Url,
			&i.Feed.CreatedAt,
			&i.Feed.UpdatedAt,
			&i.Feed.LastFetchedAt,
			&i.User.ID,
			&i.User.Name,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
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

const getFeedListSortedByLastFetchedAt = `-- name: GetFeedListSortedByLastFetchedAt :many
SELECT id, user_id, name, url, created_at, updated_at, last_fetched_at FROM "feed"
    ORDER BY "feed"."last_fetched_at" NULLS FIRST
`

func (q *Queries) GetFeedListSortedByLastFetchedAt(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeedListSortedByLastFetchedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Url,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastFetchedAt,
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

const getOldestFeedByUpdatedAt = `-- name: GetOldestFeedByUpdatedAt :one
SELECT id, user_id, name, url, created_at, updated_at, last_fetched_at FROM feed ORDER BY updated_at LIMIT 1
`

func (q *Queries) GetOldestFeedByUpdatedAt(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getOldestFeedByUpdatedAt)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const setFeedFetchedAtToNowById = `-- name: SetFeedFetchedAtToNowById :exec
UPDATE "feed" SET last_fetched_at = NOW(), updated_at = NOW() WHERE id = $1
`

func (q *Queries) SetFeedFetchedAtToNowById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, setFeedFetchedAtToNowById, id)
	return err
}
