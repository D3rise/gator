-- +goose Up
CREATE TABLE "post" (
    id UUID UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES "feed" (id)
);

-- +goose Down
DROP TABLE "post";