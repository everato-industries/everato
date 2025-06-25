-- Drop the UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";

-- Drop the tables in reverse order of creation
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS seat_types;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;

-- Drop the enums
DROP TYPE IF EXISTS role;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS payment_state;
DROP TYPE IF EXISTS seat_status;
