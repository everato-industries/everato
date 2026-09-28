-- Check-in / Attendance Operations

-- name: CreateAttendance :one
-- Records a successful check-in event.
INSERT INTO attendance (
    id,
    ticket_id,
    event_id,
    checked_in_by,
    location,
    device_info,
    checked_in_at
) VALUES (
    $1,  -- id
    $2,  -- ticket_id
    $3,  -- event_id
    $4,  -- checked_in_by (super_user id)
    $5,  -- location (nullable)
    $6,  -- device_info (nullable)
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetAttendanceByTicket :one
SELECT * FROM attendance
WHERE ticket_id = $1
ORDER BY checked_in_at DESC
LIMIT 1;

-- name: GetEventAttendanceStats :one
-- Returns check-in statistics for a given event.
SELECT
    COUNT(DISTINCT t.id)              AS total_tickets,
    COUNT(DISTINCT a.ticket_id)       AS checked_in_count,
    COUNT(DISTINCT t.id) - COUNT(DISTINCT a.ticket_id) AS not_checked_in_count
FROM bookings b
JOIN tickets t   ON t.booking_id = b.id AND t.deleted_at IS NULL
LEFT JOIN attendance a ON a.ticket_id = t.id
WHERE b.event_id = $1
  AND b.deleted_at IS NULL;

-- name: GetEventAttendees :many
-- Lists all attendees (checked-in or not) for an event.
SELECT
    t.id          AS ticket_id,
    t.ticket_number,
    t.is_checked_in,
    t.checked_in_at,
    u.email       AS user_email,
    u.first_name,
    u.last_name,
    tt.name       AS ticket_type_name
FROM bookings b
JOIN users u         ON u.id = b.user_id
JOIN tickets t       ON t.booking_id = b.id AND t.deleted_at IS NULL
JOIN ticket_types tt ON tt.id = t.ticket_type_id
WHERE b.event_id = $1
  AND b.deleted_at IS NULL
ORDER BY b.created_at ASC;
