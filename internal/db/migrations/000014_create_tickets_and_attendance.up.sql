-- Migration 014: Replace old tickets table with new QR-based ticket system
-- Also creates the attendance table for check-in tracking

-- Step 1: Drop the old tickets table (it was an early schema, no data is lost)
-- The old tickets table had no meaningful data tied to the booking flow.
DROP TABLE IF EXISTS tickets CASCADE;

-- Step 2: Create the new tickets table with QR code support
-- One ticket row represents one ticket within a booking (tied to booking_tickets line items)
CREATE TABLE IF NOT EXISTS tickets (
    id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id       UUID        NOT NULL REFERENCES bookings(id)      ON DELETE CASCADE,
    ticket_type_id   UUID        NOT NULL REFERENCES ticket_types(id)  ON DELETE CASCADE,
    ticket_number    TEXT        NOT NULL UNIQUE,      -- human-readable, e.g. "EVT2026-00001"
    qr_code_data     TEXT        NOT NULL UNIQUE,      -- HMAC-signed JSON payload
    qr_code_image    TEXT,                             -- base64-encoded PNG (nullable)
    is_checked_in    BOOLEAN     NOT NULL DEFAULT FALSE,
    checked_in_at    TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_tickets_booking        ON tickets(booking_id);
CREATE INDEX idx_tickets_ticket_type    ON tickets(ticket_type_id);
CREATE INDEX idx_tickets_number         ON tickets(ticket_number);
CREATE INDEX idx_tickets_qr_data        ON tickets(qr_code_data);

-- Step 3: Drop old payments table (references old tickets.id) and recreate
-- Payments will be added properly in a future payment migration
DROP TABLE IF EXISTS payments CASCADE;

-- Step 4: Create attendance table for check-in audit log
CREATE TABLE IF NOT EXISTS attendance (
    id              UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id       UUID        NOT NULL REFERENCES tickets(id)      ON DELETE CASCADE,
    event_id        UUID        NOT NULL REFERENCES events(id)       ON DELETE CASCADE,
    checked_in_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    checked_in_by   UUID        NOT NULL REFERENCES super_users(id)  ON DELETE CASCADE,
    location        TEXT,
    device_info     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attendance_ticket  ON attendance(ticket_id);
CREATE INDEX idx_attendance_event   ON attendance(event_id);
CREATE INDEX idx_attendance_admin   ON attendance(checked_in_by);
