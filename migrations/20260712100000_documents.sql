CREATE TABLE IF NOT EXISTS documents (
    slug          TEXT PRIMARY KEY,
    title         TEXT NOT NULL,
    category      TEXT NOT NULL,
    word_count    INTEGER NOT NULL,
    reading_time  INTEGER NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
