-- name: CreateEvent :one
INSERT INTO events (
    title,
    description,
    banner,
    icon,
    admin_id,
    start_time,
    end_time,
    location,
    total_seats,
    available_seats,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1;

-- name: SearchByName :many
SELECT * FROM events
    WHERE title ILIKE '%' || $1 || '%'
    OR description ILIKE '%' || $1 || '%'
    ORDER BY start_time DESC
LIMIT $2 OFFSET $3;

-- name: ListEvents :many
SELECT * FROM events
    ORDER BY start_time DESC
LIMIT $1 OFFSET $2;

-- name: UpdateEvent :one
UPDATE events
SET
    title = $2,
    description = $3,
    banner = $4,
    icon = $5,
    admin_id = $6,
    start_time = $7,
    end_time = $8,
    location = $9,
    total_seats = $10,
    available_seats = $11,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;

-- name: ListEventsByAdmin :many
SELECT * FROM events
    WHERE admin_id = $1
ORDER BY start_time DESC;
