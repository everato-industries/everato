-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ENUM TYPES
CREATE TYPE ROLE AS ENUM ('ADMIN', 'USER');
CREATE TYPE BOOKING_STATUS AS ENUM (
    'PENDING',
    'PENDING_PAYMENT',
    'TIMEOUT',
    'FILLED',
    'CONFIRMED',
    'CANCELLED'
);
CREATE TYPE PAYMENT_STATE AS ENUM ('PENDING', 'SUCCESS', 'FAILED');
CREATE TYPE SEAT_STATUS AS ENUM ('AVAILABLE', 'BOOKED');

-- USERS (already created above, shown here for completeness)
CREATE TABLE IF NOT EXISTS users (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(50) NOT NULL,
    last_name    VARCHAR(50) NOT NULL,
    username     VARCHAR(50) NOT NULL UNIQUE,
    email        VARCHAR(100) NOT NULL UNIQUE,
    password     VARCHAR(255) NOT NULL,
    verified     BOOLEAN NOT NULL DEFAULT FALSE,
    role         ROLE    NOT NULL DEFAULT 'USER',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- EVENTS
CREATE TABLE IF NOT EXISTS events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title           TEXT    NOT NULL,
    description     TEXT,
    banner          TEXT,
    admin_id        UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    location        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    available_seats INT    NOT NULL,
    total_seats     INT    NOT NULL
);

-- BOOKINGS
CREATE TABLE IF NOT EXISTS bookings (
    id         UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id   UUID           NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id    UUID           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status     BOOKING_STATUS NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- SEAT TYPES (per event)
CREATE TABLE IF NOT EXISTS seat_types (
    id          UUID    PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT    NOT NULL,
    price       INT     NOT NULL,
    description TEXT,
    event_id    UUID    NOT NULL REFERENCES events(id) ON DELETE CASCADE
);

-- SEATS
CREATE TABLE IF NOT EXISTS seats (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    seat_type_id UUID        NOT NULL REFERENCES seat_types(id) ON DELETE CASCADE,
    status       SEAT_STATUS NOT NULL DEFAULT 'AVAILABLE',
    booking_id   UUID        REFERENCES bookings(id) ON DELETE SET NULL
);


-- PAYMENTS
CREATE TABLE IF NOT EXISTS payments (
    id          UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    status      PAYMENT_STATE NOT NULL DEFAULT 'PENDING',
    event_id    UUID          NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id     UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    booking_id  UUID          REFERENCES bookings(id) ON DELETE SET NULL,
    amount      INT           NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- INDEXES / CONSTRAINTS

-- Each seat can only be linked to one booking at a time.
ALTER TABLE seats
  ADD CONSTRAINT unique_seat_booking UNIQUE (id, booking_id);

-- Optionally, enforce that when a booking is CONFIRMED, the number of booked seats
-- matches the event.available_seats. This would typically be done at application level
-- or via triggers.
