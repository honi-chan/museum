-- exhibitions はMUSEUMの展示室を保存するテーブル。
--
-- DomainモデルをそのままDB構造に合わせるのではなく、
-- 永続化に必要なデータだけを保持する。

CREATE TABLE exhibitions (
    -- Exhibition ID。
    --
    -- UUID自体はApplication側で生成するため、
    -- PostgreSQL側ではTEXTとして保持する。
    id TEXT PRIMARY KEY,

    -- Exhibitionが所属するMuseum。
    --
    -- Museumテーブルはまだ作成していないため、
    -- 現段階ではForeign Keyを設定しない。
    museum_id TEXT NOT NULL,

    -- 展示室タイトル。
    title TEXT NOT NULL,

    -- 展示室説明。
    --
    -- NULLと空文字の2状態を作らないように、
    -- NOT NULL + default empty string とする。
    description TEXT NOT NULL DEFAULT '',

    -- Exhibition作成日時。
    created_at TIMESTAMPTZ NOT NULL
);

-- 今後Museum単位でExhibition一覧を取得するため、
-- museum_idにindexを作成する。
CREATE INDEX idx_exhibitions_museum_id
    ON exhibitions (museum_id);