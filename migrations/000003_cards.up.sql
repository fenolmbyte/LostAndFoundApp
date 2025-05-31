CREATE TABLE IF NOT EXISTS cards
(
    id              UUID                      PRIMARY KEY,
    title           TEXT                      NOT NULL CHECK (trim(title) <> ''),
    description     TEXT,
    owner_id        UUID                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    preview_url     TEXT                      NOT NULL,
    city            TEXT                      NOT NULL CHECK (trim(city) <> ''),
    street          TEXT                      NOT NULL,
    location        geography(Point, 4326),
    status          TEXT                      NOT NULL CHECK (status IN ('lost', 'found')),
    created_at      TIMESTAMP                 NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS card_images
(
    id              UUID     PRIMARY KEY,
    card_id      UUID     REFERENCES cards (id),
    url             TEXT     NOT NULL CHECK (trim(url) <> '')
);

CREATE INDEX IF NOT EXISTS idx_cards_location ON cards USING GIST (location);