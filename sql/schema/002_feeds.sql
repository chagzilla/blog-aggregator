-- +goose Up
CREATE TABLE feeds (
    name VARCHAR NOT NULL,
    url VARCHAR UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;
