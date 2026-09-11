CREATE TABLE exhibits (
    id TEXT PRIMARY KEY,
    exhibition_id TEXT NOT NULL REFERENCES exhibitions(id) ON DELETE CASCADE,
    title TEXT NOT NULL DEFAULT '',
    caption TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL,
    display_order INTEGER NOT NULL CHECK (display_order >= 0),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_exhibits_exhibition_order ON exhibits (exhibition_id, display_order, id);
