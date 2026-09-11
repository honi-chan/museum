-- name: CreateExhibition :exec
INSERT INTO exhibitions (
    id,
    museum_id,
    title,
    description,
    display_order,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    COALESCE(
        (
            SELECT MAX(display_order) + 1
            FROM exhibitions
            WHERE museum_id = $2
        ),
        0
    ),
    $5
);

-- name: GetExhibitionByID :one
SELECT
    id,
    museum_id,
    title,
    description,
    display_order,
    created_at
FROM exhibitions
WHERE id = $1
LIMIT 1;

-- name: ListExhibitionsByMuseumID :many
SELECT
    id,
    museum_id,
    title,
    description,
    display_order,
    created_at
FROM exhibitions
WHERE museum_id = $1
ORDER BY
    display_order ASC,
    id ASC;