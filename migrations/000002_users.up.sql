CREATE TABLE IF NOT EXISTS users
(
    id            UUID        PRIMARY KEY,
    email         TEXT        UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    surname       TEXT        NOT NULL,
    phone         TEXT        NOT NULL,
    telegram      TEXT        NOT NULL,
    is_admin      BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);