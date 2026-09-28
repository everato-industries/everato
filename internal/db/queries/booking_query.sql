-- Booking Operations

-- name: CreateBooking :one
INSERT INTO bookings (
    id,
    user_id,
    event_id,
    total_amount,
    status,
    coupon_code,
    discount_amount,
    final_amount,
    created_at,
    updated_at
) VALUES (
    $1,  -- id
    $2,  -- user_id
    $3,  -- event_id
    $4,  -- total_amount
    $5,  -- status
    $6,  -- coupon_code
    $7,  -- discount_amount
    $8,  -- final_amount
    CURRENT_TIMESTAMP,  -- created_at
    CURRENT_TIMESTAMP   -- updated_at
) RETURNING *;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserBookings :many
SELECT b.*, e.title as event_title, e.start_time as event_start_time, e.end_time as event_end_time
FROM bookings b
JOIN events e ON b.event_id = e.id
WHERE b.user_id = $1 AND b.deleted_at IS NULL
ORDER BY b.created_at DESC;

-- name: GetBookingsByEvent :many
SELECT b.*, u.email as user_email, u.first_name, u.last_name
FROM bookings b
JOIN users u ON b.user_id = u.id
WHERE b.event_id = $1 AND b.deleted_at IS NULL
ORDER BY b.created_at DESC;

-- name: UpdateBookingStatus :one
UPDATE bookings
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CancelBooking :one
UPDATE bookings
SET status = 'CANCELLED', deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: CountBookingsByEvent :one
SELECT COUNT(*) as count FROM bookings
WHERE event_id = $1 AND status != 'CANCELLED' AND deleted_at IS NULL;

-- name: GetUserBookingCountForEvent :one
SELECT COUNT(*) as count FROM bookings
WHERE user_id = $1 AND event_id = $2 AND status != 'CANCELLED' AND deleted_at IS NULL;

-- Booking Ticket Operations

-- name: CreateBookingTicket :one
INSERT INTO booking_tickets (
    id,
    booking_id,
    ticket_type_id,
    quantity,
    price_per_unit,
    subtotal,
    created_at
) VALUES (
    $1,  -- id
    $2,  -- booking_id
    $3,  -- ticket_type_id
    $4,  -- quantity
    $5,  -- price_per_unit
    $6,  -- subtotal
    CURRENT_TIMESTAMP  -- created_at
) RETURNING *;

-- name: GetBookingTicketsByBookingID :many
SELECT bt.*, tt.name as ticket_type_name
FROM booking_tickets bt
JOIN ticket_types tt ON bt.ticket_type_id = tt.id
WHERE bt.booking_id = $1;

-- name: GetTicketAvailability :one
SELECT
    tt.id,
    tt.name,
    tt.price,
    tt.available_tickets as total_available,
    tt.available_tickets - COALESCE(SUM(bt.quantity), 0) as remaining
FROM ticket_types tt
LEFT JOIN booking_tickets bt ON bt.ticket_type_id = tt.id
LEFT JOIN bookings b ON bt.booking_id = b.id AND b.status != 'CANCELLED' AND b.deleted_at IS NULL
WHERE tt.id = $1
GROUP BY tt.id, tt.name, tt.price, tt.available_tickets;

-- name: GetEventTicketAvailability :many
SELECT
    tt.id,
    tt.name,
    tt.price,
    tt.available_tickets as total_available,
    tt.available_tickets - COALESCE(SUM(bt.quantity), 0) as remaining
FROM ticket_types tt
LEFT JOIN booking_tickets bt ON bt.ticket_type_id = tt.id
LEFT JOIN bookings b ON bt.booking_id = b.id AND b.status != 'CANCELLED' AND b.deleted_at IS NULL
WHERE tt.event_id = $1
GROUP BY tt.id, tt.name, tt.price, tt.available_tickets
ORDER BY tt.price ASC;

-- name: DecrementTicketAvailability :one
UPDATE ticket_types
SET available_tickets = available_tickets - $2
WHERE id = $1 AND available_tickets >= $2
RETURNING *;

-- name: IncrementTicketAvailability :one
UPDATE ticket_types
SET available_tickets = available_tickets + $2
WHERE id = $1
RETURNING *;
