-- Ticket Operations

-- name: CreateTicket :one
-- Creates a new ticket record after a booking is confirmed.
INSERT INTO tickets (
    id,
    booking_id,
    ticket_type_id,
    ticket_number,
    qr_code_data,
    qr_code_image,
    is_checked_in,
    created_at
) VALUES (
    $1,  -- id (UUID)
    $2,  -- booking_id
    $3,  -- ticket_type_id
    $4,  -- ticket_number  (e.g. "EVT-2026-00001")
    $5,  -- qr_code_data   (HMAC-signed JSON string)
    $6,  -- qr_code_image  (base64 PNG, nullable)
    FALSE,
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetTicketByID :one
SELECT * FROM tickets
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTicketByQRData :one
-- Retrieves a ticket by its unique QR payload string (used during check-in validation).
SELECT * FROM tickets
WHERE qr_code_data = $1 AND deleted_at IS NULL;

-- name: GetTicketByNumber :one
SELECT * FROM tickets
WHERE ticket_number = $1 AND deleted_at IS NULL;

-- name: GetTicketsByBookingID :many
SELECT t.*, tt.name AS ticket_type_name
FROM tickets t
JOIN ticket_types tt ON t.ticket_type_id = tt.id
WHERE t.booking_id = $1 AND t.deleted_at IS NULL
ORDER BY t.created_at ASC;

-- name: MarkTicketCheckedIn :one
-- Marks a ticket as checked-in and records the timestamp.
UPDATE tickets
SET is_checked_in = TRUE,
    checked_in_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CountTicketsForBooking :one
SELECT COUNT(*) AS count
FROM tickets
WHERE booking_id = $1 AND deleted_at IS NULL;

-- name: GetLatestTicketNumberSeq :one
-- Returns the highest-numbered ticket in the system for sequential numbering.
SELECT ticket_number FROM tickets
WHERE ticket_number LIKE $1
ORDER BY created_at DESC
LIMIT 1;
