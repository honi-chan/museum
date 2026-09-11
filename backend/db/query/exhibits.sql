-- name: LockExhibitionForExhibit :one
SELECT id FROM exhibitions WHERE id = $1 FOR UPDATE;

-- name: CreateExhibit :one
INSERT INTO exhibits (id, exhibition_id, title, caption, image_url, display_order, created_at)
VALUES ($1, $2, $3, $4, $5,
    COALESCE((SELECT MAX(display_order) + 1 FROM exhibits WHERE exhibition_id = $2), 0),
    $6)
RETURNING *;

-- name: ListExhibitsByExhibitionID :many
SELECT * FROM exhibits WHERE exhibition_id = $1 ORDER BY display_order ASC, id ASC;
