-- name: CreateExhibition :exec
--
-- Exhibitionを1件保存する。

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


-- name: GetExhibitionByID :one
--
-- IDを指定してExhibitionを1件取得する。
--
-- :one を指定するとsqlcが
-- 1件取得用のGoメソッドを自動生成する。
SELECT
    id,
    museum_id,
    title,
    description,
    created_at
FROM exhibitions
WHERE id = $1
LIMIT 1;