// Package checkin provides services for QR-based event attendance management.
// Admins (super_users) scan attendee QR codes to mark them as checked-in.
package checkin

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/services/ticket"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// ValidateAndCheckinDTO is the request body for scanning a QR code and checking in.
type ValidateAndCheckinDTO struct {
	QRData     string `json:"qr_data"`      // the raw QR string scanned by camera
	Location   string `json:"location"`     // optional — gate/hall label
	DeviceInfo string `json:"device_info"`  // optional — scanner device name
}

// ValidateAndCheckin processes a QR scan, verifies the HMAC signature,
// marks the ticket as checked-in, and records the attendance.
//
// Route: POST /api/v1/checkin/scan
//
// Steps:
//  1. Parse request body
//  2. Verify HMAC signature on QR payload
//  3. Look up ticket in DB
//  4. Reject if already checked in
//  5. Mark ticket checked in, write attendance record
//  6. Return success with attendee info
func ValidateAndCheckin(
	wr *utils.HttpWriter,
	repo *repository.Queries,
	conn *pgx.Conn,
	hmacSecret string,
) {
	logger := pkg.NewLogger()
	defer logger.Close()

	dto := &ValidateAndCheckinDTO{}
	if err := wr.ParseBody(dto); err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	if dto.QRData == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "qr_data is required"})
		return
	}

	// 1. Verify HMAC signature
	payload, err := ticket.VerifyQRData(dto.QRData, hmacSecret)
	if err != nil {
		logger.StdoutLogger.Warn("Invalid QR scan attempt", "error", err, "qr_data", dto.QRData[:20])
		wr.Status(http.StatusUnprocessableEntity).Json(utils.M{
			"error":   "Invalid or tampered QR code",
			"message": err.Error(),
		})
		return
	}

	ctx := context.Background()

	// 2. Look up ticket by the QR data string
	t, err := repo.GetTicketByQRData(ctx, dto.QRData)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusNotFound).Json(utils.M{"error": "Ticket not found"})
		} else {
			logger.Error("DB error fetching ticket", "error", err)
			wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Database error"})
		}
		return
	}

	// 3. Reject double check-in
	if t.IsCheckedIn {
		checkedAt := ""
		if t.CheckedInAt.Valid {
			checkedAt = t.CheckedInAt.Time.UTC().Format(time.RFC3339)
		}
		wr.Status(http.StatusConflict).Json(utils.M{
			"error":          "Ticket already checked in",
			"ticket_number":  t.TicketNumber,
			"checked_in_at":  checkedAt,
		})
		return
	}

	// 4. Start transaction: mark ticket + insert attendance
	tx, err := conn.Begin(ctx)
	if err != nil {
		logger.Error("Failed to start transaction", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Internal server error"})
		return
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	qtx := repo.WithTx(tx)

	// 5. Mark ticket checked in
	updatedTicket, err := qtx.MarkTicketCheckedIn(ctx, t.ID)
	if err != nil {
		logger.Error("Failed to mark ticket checked in", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to check in ticket"})
		return
	}

	// 6. Parse IDs for attendance record
	eventUUID := pgtype.UUID{}
	if scanErr := eventUUID.Scan(payload.EventID); scanErr != nil {
		logger.Error("Failed to parse event UUID", "error", scanErr)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Internal server error"})
		return
	}

	// Get the admin user ID from the request context (set by admin/auth middleware)
	var adminIDVal any = wr.R.Context().Value("admin_id")
	if adminIDVal == nil {
		adminIDVal = wr.R.Context().Value("uid")
	}
	adminUUID := pgtype.UUID{}
	if adminIDVal != nil {
		if idStr, ok := adminIDVal.(string); ok {
			_ = adminUUID.Scan(idStr)
		}
	}
	if !adminUUID.Valid {
		admins, err := repo.GetAllAdmins(ctx)
		if err == nil && len(admins) > 0 {
			adminUUID = admins[0].ID
		}
	}

	locationPg := pgtype.Text{}
	if dto.Location != "" {
		locationPg = pgtype.Text{String: dto.Location, Valid: true}
	}
	devicePg := pgtype.Text{}
	if dto.DeviceInfo != "" {
		devicePg = pgtype.Text{String: dto.DeviceInfo, Valid: true}
	}

	attendanceID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	_, err = qtx.CreateAttendance(ctx, repository.CreateAttendanceParams{
		ID:          attendanceID,
		TicketID:    t.ID,
		EventID:     eventUUID,
		CheckedInBy: adminUUID,
		Location:    locationPg,
		DeviceInfo:  devicePg,
	})
	if err != nil {
		logger.Error("Failed to write attendance record", "error", err)
		// Non-fatal: ticket is already marked, still commit
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		logger.Error("Failed to commit check-in transaction", "error", commitErr)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to finalize check-in"})
		return
	}

	logger.Info("Ticket checked in successfully",
		"ticket_id", t.ID.String(),
		"ticket_number", t.TicketNumber,
	)

	wr.Status(http.StatusOK).Json(utils.M{
		"message": "Check-in successful",
		"data": utils.M{
			"ticket_id":      updatedTicket.ID.String(),
			"ticket_number":  updatedTicket.TicketNumber,
			"checked_in_at":  updatedTicket.CheckedInAt.Time.UTC().Format(time.RFC3339),
		},
	})
}

// GetEventAttendanceStats returns check-in statistics for an event.
//
// Route: GET /api/v1/checkin/events/{eventID}/stats
func GetEventAttendanceStats(wr *utils.HttpWriter, repo *repository.Queries) {
	logger := pkg.NewLogger()
	defer logger.Close()

	eventIDStr := utils.GetParam(wr.R, "eventID")
	if eventIDStr == "" {
		eventIDStr = utils.GetParam(wr.R, "eventId")
	}
	if eventIDStr == "" {
		eventIDStr = mux.Vars(wr.R)["eventID"]
	}
	if eventIDStr == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "eventID is required"})
		return
	}

	eventUUID, err := utils.StringToUUID(eventIDStr)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "invalid event ID"})
		return
	}

	stats, err := repo.GetEventAttendanceStats(context.Background(), eventUUID)
	if err != nil {
		logger.Error("Failed to fetch attendance stats", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to fetch stats"})
		return
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"data": utils.M{
			"total_tickets":        stats.TotalTickets,
			"checked_in_count":     stats.CheckedInCount,
			"not_checked_in_count": stats.NotCheckedInCount,
		},
	})
}

// GetEventAttendees returns the full attendee list for an event.
//
// Route: GET /api/v1/checkin/events/{eventID}/attendees
func GetEventAttendees(wr *utils.HttpWriter, repo *repository.Queries) {
	logger := pkg.NewLogger()
	defer logger.Close()

	eventIDStr := utils.GetParam(wr.R, "eventID")
	if eventIDStr == "" {
		eventIDStr = utils.GetParam(wr.R, "eventId")
	}
	if eventIDStr == "" {
		eventIDStr = mux.Vars(wr.R)["eventID"]
	}
	if eventIDStr == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "eventID is required"})
		return
	}

	eventUUID, err := utils.StringToUUID(eventIDStr)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "invalid event ID"})
		return
	}

	rows, err := repo.GetEventAttendees(context.Background(), eventUUID)
	if err != nil {
		logger.Error("Failed to fetch attendees", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to fetch attendees"})
		return
	}

	attendees := make([]utils.M, 0, len(rows))
	for _, row := range rows {
		checkedInAt := ""
		if row.CheckedInAt.Valid {
			checkedInAt = row.CheckedInAt.Time.UTC().Format(time.RFC3339)
		}
		attendees = append(attendees, utils.M{
			"ticket_id":        row.TicketID.String(),
			"ticket_number":    row.TicketNumber,
			"is_checked_in":    row.IsCheckedIn,
			"checked_in_at":    checkedInAt,
			"user_email":       row.UserEmail,
			"first_name":       row.FirstName,
			"last_name":        row.LastName,
			"ticket_type_name": row.TicketTypeName,
		})
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"data": utils.M{
			"attendees": attendees,
			"count":     len(attendees),
		},
	})
}
