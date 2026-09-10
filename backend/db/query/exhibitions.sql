-- name: CreateExhibition :exec
--
-- Exhibitionを1件保存する。
--
-- Go側のQueryコードはsqlcが自動生成するため、
-- pgx.Exec()をRepository内で直接書かない。

INSERT INTO exhibitions (
    id,
    museum_id,
    title,
    description,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);