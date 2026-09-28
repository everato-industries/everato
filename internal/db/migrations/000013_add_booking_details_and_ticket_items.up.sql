-- Add financial columns to bookings table
ALTER TABLE bookings ADD COLUMN total_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN final_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN coupon_code TEXT;
ALTER TABLE bookings ADD COLUMN notes TEXT;
ALTER TABLE bookings ADD COLUMN deleted_at TIMESTAMPTZ;

-- Create booking_tickets table for line items
CREATE TABLE booking_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    ticket_type_id UUID NOT NULL REFERENCES ticket_types(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_per_unit DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX idx_booking_tickets_booking ON booking_tickets(booking_id);
CREATE INDEX idx_booking_tickets_ticket_type ON booking_tickets(ticket_type_id);
CREATE INDEX idx_bookings_deleted ON bookings(deleted_at);
CREATE INDEX idx_bookings_event_status ON bookings(event_id, status);
