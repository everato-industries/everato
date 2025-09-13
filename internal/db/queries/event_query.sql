-- name: CreateEvent :one
INSERT INTO events (
    title,
    description,
    slug,
    banner,
    icon,
    admin_id,
    start_time,
    end_time,
    location,
    total_seats,
    available_seats,
    status,
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
    $11,
    $12,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1;

-- name: GetEventBySlug :one
SELECT * FROM events
WHERE slug = $1;

-- name: SearchSlug :one
SELECT slug FROM events
WHERE slug = $1;

-- name: SearchByName :many
SELECT * FROM events
    WHERE title ILIKE '%' || $1 || '%'
    OR description ILIKE '%' || $1 || '%'
    OR slug ILIKE '%' || $1 || '%'
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
    slug = $4,
    banner = $5,
    icon = $6,
    admin_id = $7,
    start_time = $8,
    end_time = $9,
    location = $10,
    total_seats = $11,
    available_seats = $12,
    status = $13,
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

-- name: CountTotalEvents :one
SELECT COUNT(*) as total_events FROM events;

-- name: CountEventsByStatus :one
SELECT COUNT(*) as count FROM events
WHERE status = $1;

-- name: CountUpcomingEvents :one
SELECT COUNT(*) as upcoming_events FROM events
WHERE start_time > CURRENT_TIMESTAMP;

-- name: GetDashboardStats :one
SELECT
    COUNT(*) as total_events,
    COUNT(CASE WHEN status = 'CREATED' THEN 1 END) as created_events,
    COUNT(CASE WHEN status = 'STARTED' THEN 1 END) as active_events,
    COUNT(CASE WHEN status = 'COMPLETED' THEN 1 END) as completed_events,
    COUNT(CASE WHEN status = 'CANCELLED' THEN 1 END) as cancelled_events,
    COUNT(CASE WHEN start_time > CURRENT_TIMESTAMP THEN 1 END) as upcoming_events
FROM events;

-- name: GetRecentEvents :many
SELECT * FROM events
WHERE start_time >= CURRENT_TIMESTAMP - INTERVAL '1 day'
ORDER BY
    CASE
        WHEN start_time > CURRENT_TIMESTAMP THEN ABS(EXTRACT(EPOCH FROM (start_time - CURRENT_TIMESTAMP)))
        ELSE ABS(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - start_time)))
    END ASC
LIMIT $1;
