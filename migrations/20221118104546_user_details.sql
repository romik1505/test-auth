-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id uuid PRIMARY KEY default gen_random_uuid(),
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password bytea NOT NULL,
    phone VARCHAR(20) NOT NULL
);

CREATE TABLE sessions (
    id uuid PRIMARY KEY default gen_random_uuid(),
    user_id uuid NOT NULL,
    refresh_token bytea NOT NULL,
    expires_in INT NOT NULL,
    created_at timestamptz NOT NULL default now(),

    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE sessions;
DROP TABLE users;

-- +goose StatementEnd
