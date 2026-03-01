# Backend TODO - MVP Focused

**Last Updated:** March 1, 2026  
**Based On:** REPORT.md MVP Gap Analysis  
**Reference:** Sahotsava Backend RBAC Implementation

---

## 🚨 **CRITICAL PATH - MVP BLOCKERS (P0)**

These features MUST be completed before MVP can ship to clients.

### **CP-1: User Registration System** ⏱️ 1-2 days

**Status:** 70% Complete - Route commented out in `auth_handler.go:91`  
**Priority:** P0 - CRITICAL  
**Files:** `internal/handlers/v1/api/auth_handler.go`, `www/src/pages/auth/register.tsx`

**Backend Tasks:**
- [ ] Uncomment registration route in `auth_handler.go` line 91
- [ ] Test user registration flow end-to-end
- [ ] Add email verification trigger after registration
- [ ] Test password hashing with existing bcrypt implementation
- [ ] Validate email uniqueness and format
- [ ] Return proper JWT tokens after registration
- [ ] Add user profile creation workflow

**Testing Checklist:**
- [ ] Register with valid email/password
- [ ] Attempt duplicate registration (should fail)
- [ ] Verify email format validation
- [ ] Confirm password is hashed in database
- [ ] Verify JWT token generation

**Reference:** Sahotsava uses pre-registered users only, but Everato needs public registration

---

### **CP-2: Ticket Booking System** ⏱️ 5-7 days

**Status:** 0% Complete - Database ready, no implementation  
**Priority:** P0 - CRITICAL BLOCKER  
**Effort:** Highest complexity - core business logic

**New Files to Create:**

1. **Handler:** `internal/handlers/v1/api/booking_handler.go`
   ```go
   type BookingHandler struct {
       Repo     *repository.Queries
       Conn     *pgx.Conn
       BasePath string
       Config   *config.Config
   }
   ```

2. **Service Directory:** `internal/services/booking/`
   - `booking_create.go` - Create booking with seat validation
   - `booking_validate.go` - Availability check, max tickets per user
   - `booking_dto.go` - Request/response DTOs
   - `booking_get.go` - Retrieve user bookings
   - `booking_cancel.go` - Cancel bookings with refund logic

3. **SQL Queries:** `internal/db/queries/booking_query.sql`
   ```sql
   -- name: CreateBooking :one
   INSERT INTO bookings (
       id, user_id, event_id, total_amount, 
       booking_status, created_at
   ) VALUES ($1, $2, $3, $4, $5, $6)
   RETURNING *;

   -- name: GetUserBookings :many
   SELECT b.*, e.title as event_title, e.event_date
   FROM bookings b
   JOIN events e ON b.event_id = e.id
   WHERE b.user_id = $1 AND b.deleted_at IS NULL
   ORDER BY b.created_at DESC;

   -- name: GetBookingsByEvent :many
   SELECT b.*, u.email as user_email
   FROM bookings b
   JOIN users u ON b.user_id = u.id
   WHERE b.event_id = $1 AND b.deleted_at IS NULL;

   -- name: UpdateBookingStatus :one
   UPDATE bookings 
   SET booking_status = $2, updated_at = NOW()
   WHERE id = $1
   RETURNING *;

   -- name: CreateBookingTicket :one
   INSERT INTO booking_tickets (
       id, booking_id, ticket_type_id, quantity, price
   ) VALUES ($1, $2, $3, $4, $5)
   RETURNING *;

   -- name: CheckTicketAvailability :one
   SELECT 
       tt.quantity_available - COALESCE(SUM(bt.quantity), 0) as remaining
   FROM ticket_types tt
   LEFT JOIN booking_tickets bt ON bt.ticket_type_id = tt.id
   LEFT JOIN bookings b ON bt.booking_id = b.id AND b.booking_status != 'CANCELLED'
   WHERE tt.id = $1 AND tt.deleted_at IS NULL
   GROUP BY tt.id, tt.quantity_available;
   ```

**API Endpoints:**
```
POST   /api/v1/bookings/create
GET    /api/v1/bookings/user/{userId}
GET    /api/v1/bookings/{bookingId}
PUT    /api/v1/bookings/{bookingId}/status
DELETE /api/v1/bookings/{bookingId}
GET    /api/v1/events/{eventId}/availability
```

**Request DTO Example:**
```go
type CreateBookingRequest struct {
    EventID string `json:"event_id" validate:"required,uuid"`
    Tickets []struct {
        TicketTypeID string `json:"ticket_type_id" validate:"required,uuid"`
        Quantity     int    `json:"quantity" validate:"required,min=1,max=10"`
    } `json:"tickets" validate:"required,min=1"`
    CouponCode string `json:"coupon_code,omitempty"`
}
```

**Business Logic Requirements:**
- [ ] Transaction-based booking (atomic operation)
- [ ] Real-time seat availability check
- [ ] Maximum tickets per user validation (configurable)
- [ ] Coupon code application and validation
- [ ] Price calculation with discounts
- [ ] Booking expiry (15-minute hold before payment)
- [ ] Concurrent booking conflict resolution (use database locks)

**Reference:** Similar to Sahotsava's participant registration flow but with payment

---

### **CP-3: Payment Integration** ⏱️ 4-6 days

**Status:** 0% Complete - Models exist, no implementation  
**Priority:** P0 - CRITICAL  
**Decision Required:** Choose payment gateway (Stripe recommended)

**New Files to Create:**

1. **Handler:** `internal/handlers/v1/api/payment_handler.go`
2. **Service Directory:** `internal/services/payment/`
   - `payment_create.go` - Initialize payment session
   - `payment_verify.go` - Verify payment callback
   - `payment_webhook.go` - Handle gateway webhooks
   - `payment_refund.go` - Handle refunds for cancellations

**Dependencies to Add:**
```bash
go get github.com/stripe/stripe-go/v76
```

**SQL Queries:** `internal/db/queries/payment_query.sql`
```sql
-- name: CreatePayment :one
INSERT INTO payments (
    id, booking_id, user_id, amount, 
    currency, payment_status, provider, 
    provider_payment_id, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET payment_status = $2, 
    provider_response = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetPaymentByBooking :one
SELECT * FROM payments
WHERE booking_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 1;
```

**API Endpoints:**
```
POST   /api/v1/payments/create-intent
POST   /api/v1/payments/verify
POST   /api/v1/payments/webhook
GET    /api/v1/payments/{paymentId}/status
```

**Environment Variables Needed:**
```env
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

**Payment Flow:**
1. User completes booking → Creates pending booking
2. Backend creates Stripe PaymentIntent
3. Frontend displays Stripe Elements
4. User completes payment
5. Stripe webhook confirms payment
6. Backend updates booking status to CONFIRMED
7. Trigger ticket generation + email

**Webhook Security:**
- [ ] Verify Stripe signature
- [ ] Idempotency key handling
- [ ] Handle all payment states: succeeded, failed, canceled, refunded

---

### **CP-4: QR Code Generation** ⏱️ 2-3 days

**Status:** 0% Complete - Database field exists  
**Priority:** P0 - CRITICAL  

**Dependencies:**
```bash
go get github.com/skip2/go-qrcode
```

**New Files to Create:**

1. **Service Directory:** `internal/services/ticket/`
   - `ticket_create.go` - Create tickets after successful payment
   - `ticket_generate_qr.go` - Generate unique QR codes
   - `ticket_validate.go` - Validate QR at check-in
   - `ticket_dto.go` - Ticket DTOs

**SQL Queries:** `internal/db/queries/ticket_query.sql`
```sql
-- name: CreateTicket :one
INSERT INTO tickets (
    id, booking_id, ticket_type_id, ticket_number,
    qr_code_data, qr_code_image, is_checked_in, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTicketByQRData :one
SELECT t.*, b.*, e.*, u.email as user_email
FROM tickets t
JOIN bookings b ON t.booking_id = b.id
JOIN events e ON b.event_id = e.id
JOIN users u ON b.user_id = u.id
WHERE t.qr_code_data = $1 AND t.deleted_at IS NULL;

-- name: MarkTicketCheckedIn :one
UPDATE tickets
SET is_checked_in = true, checked_in_at = NOW()
WHERE id = $1
RETURNING *;
```

**QR Code Data Structure:**
```go
type QRCodeData struct {
    TicketID   string `json:"ticket_id"`
    BookingID  string `json:"booking_id"`
    EventID    string `json:"event_id"`
    UserID     string `json:"user_id"`
    TicketType string `json:"ticket_type"`
    IssuedAt   string `json:"issued_at"`
    Signature  string `json:"signature"`  // HMAC signature for verification
}
```

**Implementation Tasks:**
- [ ] Generate unique ticket number (e.g., `EVT-2026-00001`)
- [ ] Create QR data with HMAC signature (use JWT secret)
- [ ] Generate QR code image as PNG
- [ ] Store QR code as base64 string in database
- [ ] Create validation function with signature verification
- [ ] Handle QR code regeneration on request

**Security:**
- Use HMAC-SHA256 signature to prevent QR code forgery
- Include timestamp to detect expired/reused codes
- Rate limit QR validation endpoint

---

### **CP-5: Email Ticket Delivery** ⏱️ 2-3 days

**Status:** 40% Complete - Mailer service exists  
**Priority:** P0 - CRITICAL  
**Files:** `pkg/templates/`, `internal/services/mailer/`

**Templates to Create:**

1. `templates/mail/ticket-confirmation.html` - Ticket email with QR code
2. `templates/mail/booking-confirmation.html` - Booking confirmation
3. `templates/mail/payment-receipt.html` - Payment receipt

**Email Template Structure:**
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your Ticket - {{.EventTitle}}</title>
</head>
<body>
    <h1>🎟️ Your Ticket for {{.EventTitle}}</h1>
    
    <div class="ticket-info">
        <p><strong>Booking Reference:</strong> {{.BookingID}}</p>
        <p><strong>Event Date:</strong> {{.EventDate}}</p>
        <p><strong>Venue:</strong> {{.Venue}}</p>
        <p><strong>Ticket Type:</strong> {{.TicketType}}</p>
    </div>

    <div class="qr-code">
        <img src="data:image/png;base64,{{.QRCodeBase64}}" alt="QR Code" />
        <p>Show this QR code at the venue entrance</p>
    </div>

    <div class="ticket-details">
        <p>Ticket Number: {{.TicketNumber}}</p>
        <p>Amount Paid: ${{.AmountPaid}}</p>
    </div>
</body>
</html>
```

**Implementation Tasks:**
- [ ] Create email service wrapper in `internal/services/mailer/`
- [ ] Implement ticket email with embedded QR code
- [ ] Trigger email after payment success webhook
- [ ] Add retry logic for failed email sends
- [ ] Store email send status in database
- [ ] Add "Resend Ticket" endpoint

**Email Triggers:**
- After successful payment → Send ticket email
- After booking confirmation → Send booking summary
- Before event (24hrs) → Send event reminder
- After check-in → Send thank you email

---

### **CP-6: QR Scanner for Admins** ⏱️ 3-4 days

**Status:** 0% Complete  
**Priority:** P0 - CRITICAL for event management  

**Migration Required:**
Create new migration: `make migrate-new` → `create_attendance_table`

```sql
CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id),
    event_id UUID NOT NULL REFERENCES events(id),
    checked_in_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    checked_in_by UUID NOT NULL REFERENCES super_users(id),
    location TEXT,
    device_info TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attendance_ticket ON attendance(ticket_id);
CREATE INDEX idx_attendance_event ON attendance(event_id);
```

**New Files to Create:**

1. **Handler:** `internal/handlers/v1/api/checkin_handler.go`
2. **Service Directory:** `internal/services/checkin/`
   - `validate_ticket.go` - Verify QR code authenticity
   - `mark_attendance.go` - Mark ticket as used
   - `get_attendance.go` - Get event attendance stats

**SQL Queries:** `internal/db/queries/checkin_query.sql`
```sql
-- name: ValidateTicket :one
SELECT t.*, b.booking_status, e.event_date
FROM tickets t
JOIN bookings b ON t.booking_id = b.id
JOIN events e ON b.event_id = e.id
WHERE t.qr_code_data = $1 AND t.deleted_at IS NULL;

-- name: CreateAttendance :one
INSERT INTO attendance (
    id, ticket_id, event_id, checked_in_by
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetEventAttendanceStats :one
SELECT 
    COUNT(DISTINCT a.ticket_id) as checked_in_count,
    COUNT(DISTINCT t.id) as total_tickets,
    COUNT(DISTINCT a.checked_in_at::date) as unique_days
FROM events e
LEFT JOIN bookings b ON b.event_id = e.id
LEFT JOIN tickets t ON t.booking_id = b.id
LEFT JOIN attendance a ON a.ticket_id = t.id
WHERE e.id = $1
GROUP BY e.id;
```

**API Endpoints:**
```
POST   /api/v1/checkin/validate
POST   /api/v1/checkin/mark
GET    /api/v1/checkin/event/{eventId}/stats
GET    /api/v1/checkin/event/{eventId}/attendees
```

**Validation Logic:**
- [ ] Verify QR signature is valid
- [ ] Check ticket exists and not deleted
- [ ] Verify booking is CONFIRMED status
- [ ] Check ticket not already checked in
- [ ] Verify event date is today or in past
- [ ] Admin authorization check
- [ ] Log all check-in attempts

**Response Example:**
```go
type ValidateTicketResponse struct {
    Valid        bool   `json:"valid"`
    Message      string `json:"message"`
    TicketNumber string `json:"ticket_number,omitempty"`
    UserName     string `json:"user_name,omitempty"`
    EventTitle   string `json:"event_title,omitempty"`
    TicketType   string `json:"ticket_type,omitempty"`
    CheckedInAt  string `json:"checked_in_at,omitempty"`
}
```

---

### **CP-7: Deployment Documentation** ⏱️ 1-2 days

**Status:** 0% Complete  
**Priority:** P0 - Required for client self-hosting  

**Files to Create:**

1. `DEPLOYMENT.md` - Comprehensive deployment guide
2. `docker-compose.production.yaml` - Production Docker setup
3. `scripts/setup_admin.sh` - Admin account creation script
4. `scripts/backup_db.sh` - Database backup script

**DEPLOYMENT.md Structure:**

```markdown
# Everato Deployment Guide

## System Requirements
- Go 1.24+
- PostgreSQL 15+
- Node.js 20+ (for frontend build)
- 2GB RAM minimum
- 20GB storage

## Quick Deploy (Docker)
1. Clone repository
2. Copy .env.example to .env
3. Configure database credentials
4. Run: docker-compose -f docker-compose.production.yaml up -d
5. Create admin: ./scripts/setup_admin.sh
6. Access: http://your-domain.com

## Manual Deployment
### Database Setup
### Backend Deployment
### Frontend Build & Deployment
### Nginx Configuration
### SSL/TLS Setup

## Environment Variables Reference
## Health Check Endpoints
## Backup & Restore
## Troubleshooting
```

**Admin Setup Script:**
```bash
#!/bin/bash
# scripts/setup_admin.sh

echo "Creating super admin account..."

read -p "Enter admin email: " email
read -sp "Enter password: " password

# SQL to create admin
psql $DATABASE_URL -c "
INSERT INTO super_users (email, password, created_at)
VALUES ('$email', crypt('$password', gen_salt('bf')), NOW());
"

echo "Admin created successfully!"
```

---

## 🔥 **HIGH PRIORITY - MVP ENHANCEMENT (P1)**

These features significantly improve MVP but not strict blockers.

### **P1-1: Dashboard Statistics APIs** ⏱️ 2-3 days

**Files to Create:**
- `internal/handlers/v1/api/dashboard_handler.go`
- `internal/services/dashboard/stats.go`
- `internal/db/queries/dashboard_query.sql`

**Endpoints:**
```
GET /api/v1/dashboard/stats
GET /api/v1/dashboard/recent-events
GET /api/v1/dashboard/activity
```

**SQL Queries:**
```sql
-- name: GetDashboardStats :one
SELECT 
    COUNT(DISTINCT e.id) as total_events,
    COUNT(DISTINCT u.id) as total_users,
    COALESCE(SUM(p.amount), 0) as total_revenue,
    COUNT(DISTINCT CASE WHEN e.event_date > NOW() THEN e.id END) as upcoming_events,
    COUNT(DISTINCT t.id) as total_tickets_sold,
    COUNT(DISTINCT CASE WHEN e.event_status = 'PUBLISHED' THEN e.id END) as active_events
FROM events e
CROSS JOIN users u
LEFT JOIN bookings b ON b.event_id = e.id
LEFT JOIN payments p ON p.booking_id = b.id AND p.payment_status = 'DONE'
LEFT JOIN tickets t ON t.booking_id = b.id
WHERE e.deleted_at IS NULL AND u.deleted_at IS NULL;
```

---

### **P1-2: User Dashboard** ⏱️ 2-3 days

**Endpoints:**
```
GET /api/v1/users/{userId}/dashboard
GET /api/v1/users/{userId}/tickets
GET /api/v1/users/{userId}/bookings
```

---

### **P1-3: Coupon System Integration** ⏱️ 1-2 days

**Status:** 100% backend complete, 0% frontend  
**Task:** Connect booking flow to existing coupon APIs

---

## 📊 **NICE TO HAVE - POST-MVP (P2)**

### **P2-1: Event Analytics**
- Views tracking
- Booking trends
- Revenue reports

### **P2-2: Advanced Search & Filtering**
- Full-text search
- Category filtering
- Location-based search
- Price range filters

### **P2-3: Social Features**
- Event reviews and ratings
- Share events on social media
- Event recommendations

---

## 📅 **IMPLEMENTATION ROADMAP**

### **Week 1-2: Core Booking Flow (Days 1-14)**
- ✅ Day 1-2: Enable user registration (CP-1)
- ✅ Day 3-9: Build ticket booking system (CP-2)
- ✅ Day 10-14: Integrate payment gateway (CP-3)

### **Week 3-4: Ticketing & Communication (Days 15-28)**
- ✅ Day 15-17: QR code generation (CP-4)
- ✅ Day 18-20: Email ticket delivery (CP-5)
- ✅ Day 21-24: QR scanner for check-in (CP-6)
- ✅ Day 25-26: Deployment documentation (CP-7)
- ✅ Day 27-28: Integration testing

### **Week 5-6: Enhancement & Polish (Days 29-42)**
- ✅ Day 29-31: Dashboard APIs (P1-1)
- ✅ Day 32-34: User dashboard (P1-2)
- ✅ Day 35-36: Coupon integration (P1-3)
- ✅ Day 37-40: Bug fixes and optimization
- ✅ Day 41-42: Production testing

---

## 🔧 **TECHNICAL NOTES**

### **Database Transaction Pattern**
```go
// Use for booking creation
tx, err := h.Conn.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Perform operations
// ...

err = tx.Commit(ctx)
```

### **RBAC Pattern (from Sahotsava)**
```go
// Middleware pattern for role-based access
func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := r.Context().Value("user")
            // Validate role
            next.ServeHTTP(w, r)
        })
    }
}
```

### **Error Handling Pattern**
```go
if err != nil {
    if err == pgx.ErrNoRows {
        wr.Status(404).Json(utils.M{"error": "Not found"})
        return
    }
    logger.Error("Operation failed", "error", err)
    wr.Status(500).Json(utils.M{"error": "Internal error"})
    return
}
```

---

## ✅ **MVP COMPLETION CHECKLIST**

**Before Production Release:**

### Core Features
- [ ] User registration enabled and tested
- [ ] Ticket booking flow working end-to-end
- [ ] Payment integration tested (test mode)
- [ ] QR code generation working
- [ ] Email delivery confirmed
- [ ] QR scanner functional for admins

### Security
- [ ] SQL injection prevention verified
- [ ] JWT token security validated
- [ ] Password hashing confirmed
- [ ] HTTPS enforced
- [ ] CORS configured correctly
- [ ] Rate limiting enabled

### Documentation
- [ ] DEPLOYMENT.md complete
- [ ] API documentation generated
- [ ] Admin setup guide written
- [ ] Environment variables documented
- [ ] Troubleshooting guide created

### Testing
- [ ] Unit tests for booking logic
- [ ] Integration tests for payment flow
- [ ] Load testing for concurrent bookings
- [ ] Cross-browser frontend testing
- [ ] Mobile responsiveness verified

---

**Total Estimated Effort:** 20-30 development days (4-6 weeks for 1 developer)  
**Success Probability:** 85%+ (Solid foundation, clear requirements)
