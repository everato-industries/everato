package ticket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// GenerateTicketsForBooking creates one Ticket DB record per booked ticket
// (i.e. for a quantity of 3 it creates 3 rows), each with a unique QR code.
//
// Call this after a booking is confirmed. It runs inside a provided transaction
// so the caller can roll back on failure.
//
// Parameters:
//   - ctx:  Request context
//   - qtx:  Repository bound to a transaction
//   - req:  Ticket generation parameters
//
// Returns the list of created TicketResponse objects, or an error.
func GenerateTicketsForBooking(
	ctx context.Context,
	qtx *repository.Queries,
	req GenerateTicketsRequest,
	cfg *config.Config,
) ([]TicketResponse, error) {
	logger := pkg.NewLogger()
	defer logger.Close()

	secret := req.HMACSecret
	if secret == "" {
		secret = utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
	}

	results := make([]TicketResponse, 0, req.Quantity)

	for i := 0; i < req.Quantity; i++ {
		ticketID := uuid.New()
		ticketIDStr := ticketID.String()

		// Sequential-style ticket number: SLUG-YYYYMMDD-UUID[:6]
		ticketNum := fmt.Sprintf("%s-%s-%s",
			req.EventSlug,
			time.Now().UTC().Format("20060102"),
			ticketIDStr[:6],
		)

		// Build HMAC-signed QR payload
		qrData, err := BuildSignedQRData(
			ticketIDStr,
			req.BookingID,
			req.EventID,
			req.UserID,
			ticketNum,
			secret,
		)
		if err != nil {
			logger.Error("Failed to build QR data", "error", err)
			return nil, fmt.Errorf("build QR data for ticket %d: %w", i+1, err)
		}

		// Generate base64 PNG QR image (256px)
		qrImage, err := GenerateQRImage(qrData, 256)
		if err != nil {
			// Non-fatal: proceed without image, log warning
			logger.StdoutLogger.Warn("Failed to generate QR image, ticket saved without image", "error", err)
			qrImage = ""
		}

		// Build pgtype values
		pgTicketID := pgtype.UUID{Bytes: ticketID, Valid: true}

		bookingUUID := pgtype.UUID{}
		if err := bookingUUID.Scan(req.BookingID); err != nil {
			return nil, fmt.Errorf("parse booking_id UUID: %w", err)
		}
		typeUUID := pgtype.UUID{}
		if err := typeUUID.Scan(req.TicketTypeID); err != nil {
			return nil, fmt.Errorf("parse ticket_type_id UUID: %w", err)
		}

		qrImagePg := pgtype.Text{}
		if qrImage != "" {
			qrImagePg = pgtype.Text{String: qrImage, Valid: true}
		}

		// Persist ticket
		ticket, err := qtx.CreateTicket(ctx, repository.CreateTicketParams{
			ID:            pgTicketID,
			BookingID:     bookingUUID,
			TicketTypeID:  typeUUID,
			TicketNumber:  ticketNum,
			QrCodeData:    qrData,
			QrCodeImage:   qrImagePg,
		})
		if err != nil {
			logger.Error("Failed to insert ticket", "error", err, "index", i)
			return nil, fmt.Errorf("insert ticket %d: %w", i+1, err)
		}

		results = append(results, TicketResponse{
			ID:           ticket.ID.String(),
			BookingID:    ticket.BookingID.String(),
			TicketNumber: ticket.TicketNumber,
			QRCodeData:   ticket.QrCodeData,
			QRCodeImage:  pgTextToString(ticket.QrCodeImage),
			IsCheckedIn:  ticket.IsCheckedIn,
			CreatedAt:    pgTimeToString(ticket.CreatedAt),
		})
	}

	return results, nil
}

// GetTicketsByBooking fetches all tickets for a booking ID and returns them
// as TicketResponse slices for API consumption.
func GetTicketsByBooking(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	bookingIDStr := utils.GetParam(wr.R, "bookingID")
	if bookingIDStr == "" {
		bookingIDStr = utils.GetParam(wr.R, "bookingId")
	}
	if bookingIDStr == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "bookingID is required"})
		return
	}

	bookingUUID, err := utils.StringToUUID(bookingIDStr)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "invalid booking ID"})
		return
	}

	rows, err := repo.GetTicketsByBookingID(context.Background(), bookingUUID)
	if err != nil {
		logger.Error("Failed to fetch tickets", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to fetch tickets"})
		return
	}

	tickets := make([]utils.M, 0, len(rows))
	for _, row := range rows {
		tickets = append(tickets, utils.M{
			"id":             row.ID.String(),
			"booking_id":     row.BookingID.String(),
			"ticket_number":  row.TicketNumber,
			"is_checked_in":  row.IsCheckedIn,
			"checked_in_at":  pgTimeToString(row.CheckedInAt),
			"qr_code_data":   row.QrCodeData,
			"qr_code_image":  pgTextToString(row.QrCodeImage),
			"ticket_type_name": row.TicketTypeName,
			"created_at":     pgTimeToString(row.CreatedAt),
		})
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"data": utils.M{
			"tickets": tickets,
			"count":   len(tickets),
		},
	})
}

// GetTicketByID fetches a single ticket by its UUID.
func GetTicketByID(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	ticketIDStr := utils.GetParam(wr.R, "ticketID")
	if ticketIDStr == "" {
		ticketIDStr = utils.GetParam(wr.R, "ticketId")
	}
	if ticketIDStr == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "ticketID is required"})
		return
	}

	ticketUUID, err := utils.StringToUUID(ticketIDStr)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "invalid ticket ID"})
		return
	}

	t, err := repo.GetTicketByID(context.Background(), ticketUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusNotFound).Json(utils.M{"error": "Ticket not found"})
			return
		}
		logger.Error("Failed to fetch ticket", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to fetch ticket"})
		return
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"data": utils.M{
			"id":             t.ID.String(),
			"booking_id":     t.BookingID.String(),
			"ticket_type_id": t.TicketTypeID.String(),
			"ticket_number":  t.TicketNumber,
			"qr_code_data":   t.QrCodeData,
			"qr_code_image":  pgTextToString(t.QrCodeImage),
			"is_checked_in":  t.IsCheckedIn,
			"checked_in_at":  pgTimeToString(t.CheckedInAt),
			"created_at":     pgTimeToString(t.CreatedAt),
		},
	})
}

