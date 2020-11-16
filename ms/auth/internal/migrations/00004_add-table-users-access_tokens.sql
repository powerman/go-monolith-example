-- +goose Up
CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TABLE users (
    id                  TEXT        NOT NULL CHECK (id ~ '^[a-z0-9-]{4,63}$'),
    pass_salt           BYTEA       NOT NULL,
    pass_hash           BYTEA       NOT NULL,
    email               TEXT        NOT NULL,
    display_name        TEXT        NOT NULL,
    role                role        NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_unique_lower_email_idx ON users (LOWER(email));

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE access_tokens (
    access_token        TEXT        NOT NULL,
    user_id             TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (access_token)
);

CREATE INDEX access_tokens_user_id_idx ON access_tokens (user_id);

-- +goose Down
DROP TABLE access_tokens;
DROP TABLE users;
DROP TYPE role;
