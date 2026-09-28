-- Drop booking_tickets table
DROP TABLE IF EXISTS booking_tickets;

-- Remove columns from bookings table
ALTER TABLE bookings DROP COLUMN IF EXISTS total_amount;
ALTER TABLE bookings DROP COLUMN IF EXISTS discount_amount;
ALTER TABLE bookings DROP COLUMN IF EXISTS final_amount;
ALTER TABLE bookings DROP COLUMN IF EXISTS coupon_code;
ALTER TABLE bookings DROP COLUMN IF EXISTS notes;
ALTER TABLE bookings DROP COLUMN IF EXISTS deleted_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_booking_tickets_booking;
DROP INDEX IF EXISTS idx_booking_tickets_ticket_type;
DROP INDEX IF EXISTS idx_bookings_deleted;
DROP INDEX IF EXISTS idx_bookings_event_status;
