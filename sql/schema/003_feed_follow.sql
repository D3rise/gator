-- +goose Up
CREATE TABLE feed_follow (
    id UUID UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES "feed" (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follow;