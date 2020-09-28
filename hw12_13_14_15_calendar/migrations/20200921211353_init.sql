-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    email VARCHAR(256) NOT NULL UNIQUE,
    first_name VARCHAR(128),
    last_name VARCHAR(128),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL,
    user_id BIGINT NOT NULL,
    title VARCHAR(256) NOT NULL,
    content TEXT,
    date_from TIMESTAMPTZ,
    date_to TIMESTAMPTZ,
    PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS user_id_idx ON events(user_id);

-- +goose Down
DROP TABLE users;
DROP TABLE events;
