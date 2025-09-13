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
    organizer_name,
    organizer_email,
    organizer_phone,
    organization,
    contact_email,
    contact_phone,
    refund_policy,
    terms_and_conditions,
    event_type,
    category,
    max_tickets_per_user,
    booking_start_time,
    booking_end_time,
    tags,
    website_url,
    facebook_url,
    twitter_url,
    instagram_url,
    linkedin_url,
    venue_name,
    address_line1,
    address_line2,
    city,
    state,
    postal_code,
    country,
    latitude,
    longitude,
    created_at,
    updated_at
) VALUES (
    $1,  -- title
    $2,  -- description
    $3,  -- slug
    $4,  -- banner
    $5,  -- icon
    $6,  -- admin_id
    $7,  -- start_time
    $8,  -- end_time
    $9,  -- location
    $10, -- total_seats
    $11, -- available_seats
    $12, -- status
    $13, -- organizer_name
    $14, -- organizer_email
    $15, -- organizer_phone
    $16, -- organization
    $17, -- contact_email
    $18, -- contact_phone
    $19, -- refund_policy
    $20, -- terms_and_conditions
    $21, -- event_type
    $22, -- category
    $23, -- max_tickets_per_user
    $24, -- booking_start_time
    $25, -- booking_end_time
    $26, -- tags
    $27, -- website_url
    $28, -- facebook_url
    $29, -- twitter_url
    $30, -- instagram_url
    $31, -- linkedin_url
    $32, -- venue_name
    $33, -- address_line1
    $34, -- address_line2
    $35, -- city
    $36, -- state
    $37, -- postal_code
    $38, -- country
    $39, -- latitude
    $40, -- longitude
    CURRENT_TIMESTAMP, -- created_at
    CURRENT_TIMESTAMP  -- updated_at
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

-- Ticket Type Operations

-- name: CreateTicketType :one
INSERT INTO ticket_types (
    name,
    event_id,
    price,
    available_tickets
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetTicketTypesByEventID :many
SELECT * FROM ticket_types
WHERE event_id = $1
ORDER BY price ASC;

-- name: UpdateTicketTypeAvailability :one
UPDATE ticket_types
SET available_tickets = $2
WHERE id = $1
RETURNING *;

-- Coupon Operations

-- name: CreateCoupon :one
INSERT INTO coupons (
    event_id,
    code,
    discount_percentage,
    valid_from,
    valid_until,
    usage_limit,
    usage_count
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    0
) RETURNING *;

-- name: GetCouponsByEventID :many
SELECT * FROM coupons
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: GetValidCouponByCode :one
SELECT * FROM coupons
WHERE code = $1
    AND valid_from <= CURRENT_TIMESTAMP
    AND valid_until >= CURRENT_TIMESTAMP
    AND usage_count < usage_limit
LIMIT 1;
