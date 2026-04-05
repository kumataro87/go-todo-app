CREATE TABLE IF NOT EXISTS todos (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    completed   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 完了フラグでのフィルタリング用インデックス
CREATE INDEX idx_todos_completed ON todos (completed);

-- デフォルトソート（新着順）用インデックス
CREATE INDEX idx_todos_created_at ON todos (created_at DESC);
