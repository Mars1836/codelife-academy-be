CREATE TABLE IF NOT EXISTS user_document_progress (
    user_id             TEXT NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    document_slug       TEXT NOT NULL REFERENCES documents(slug) ON DELETE CASCADE,
    status              TEXT NOT NULL DEFAULT 'unread'
                        CHECK (status IN ('unread', 'studying', 'completed')),
    scroll_position     INTEGER NOT NULL DEFAULT 0 CHECK (scroll_position >= 0),
    note                TEXT NOT NULL DEFAULT '',
    checked_flashcards  JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, document_slug)
);

CREATE INDEX IF NOT EXISTS idx_user_document_progress_updated
    ON user_document_progress (user_id, updated_at DESC);
