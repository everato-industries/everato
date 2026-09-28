// Package api provides HTTP handlers for the Everato platform REST API.
package api

import (
	"context"
	"io"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/internal/services/payment"
	"github.com/dtg-lucifer/everato/internal/services/ticket"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// PaymentHandler manages payment operations with multi-provider support.
// Supported gateways: Razorpay, Stripe, Cashfree.
//
// Route prefix:
//   - `/api/v1/payments`
//
// Routes:
//   - GET /payments/config - Get active payment gateway public client config (public)
//   - POST /payments/create-order - Create order / payment intent for a booking (protected)
//   - POST /payments/verify - Verify payment and confirm booking (protected)
//   - POST /payments/webhook - Webhook endpoint for async gateway updates (public)
type PaymentHandler struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	BasePath string
	Cfg      *config.Config
	Gateway  payment.PaymentGateway
}

var _ handlers.Handler = (*PaymentHandler)(nil)

// NewPaymentHandler initializes the payment handler and selected gateway.
func NewPaymentHandler(cfg *config.Config) *PaymentHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("PaymentHandler: Error connecting to the postgres db", "err", err.Error())
	}

	var repo *repository.Queries
	if conn != nil {
		repo = repository.New(conn)
	}

	// Try to initialize payment gateway based on config
	gateway, err := payment.NewPaymentGateway(cfg)
	if err != nil {
		logger.StdoutLogger.Warn(
			"Payment gateway initialization warning (free bookings will still work)",
			"error", err.Error(),
		)
	}

	return &PaymentHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/payments",
		Cfg:      cfg,
		Gateway:  gateway,
	}
}

// RegisterRoutes registers all payment-related endpoints.
func (h *PaymentHandler) RegisterRoutes(router *mux.Router) {
	payments := router.PathPrefix(h.BasePath).Subrouter()

	// Public endpoints
	payments.HandleFunc("/config", h.GetConfig).Methods(http.MethodGet)
	payments.HandleFunc("/webhook", h.HandleWebhook).Methods(http.MethodPost)

	// Protected user endpoints
	guard := middlewares.NewAuthMiddleware(h.Repo, h.Conn, false)
	protected := payments.NewRoute().Subrouter()
	protected.Use(guard.Guard)

	protected.HandleFunc("/create-order", h.CreateOrder).Methods(http.MethodPost)
	protected.HandleFunc("/verify", h.VerifyPayment).Methods(http.MethodPost)
}

// GetConfig returns the active payment provider name and public keys for the frontend.
func (h *PaymentHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	if h.Gateway == nil {
		wr.Status(http.StatusOK).Json(utils.M{
			"provider": "none",
			"enabled":  false,
			"message":  "No payment gateway credentials configured. Free bookings will be auto-confirmed.",
		})
		return
	}

	clientConfig := h.Gateway.GetClientConfig()
	clientConfig["enabled"] = "true"

	wr.Status(http.StatusOK).Json(utils.M{
		"data": clientConfig,
	})
}

// CreateOrderRequest is the payload for /payments/create-order
type createOrderBody struct {
	BookingID string `json:"booking_id"`
}

// CreateOrder prepares an order/intent with the selected payment provider.
func (h *PaymentHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	logger := pkg.NewLogger()
	defer logger.Close()

	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}

	body := &createOrderBody{}
	if err := wr.ParseBody(body); err != nil || body.BookingID == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "booking_id is required"})
		return
	}

	bookingUUID, err := utils.StringToUUID(body.BookingID)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "invalid booking_id UUID"})
		return
	}

	ctx := r.Context()
	booking, err := h.Repo.GetBookingByID(ctx, bookingUUID)
	if err != nil {
		wr.Status(http.StatusNotFound).Json(utils.M{"error": "Booking not found"})
		return
	}

	// Calculate final float amount
	finalAmount := 0.0
	if booking.FinalAmount.Valid {
		fVal, _ := booking.FinalAmount.Float64Value()
		if fVal.Valid {
			finalAmount = fVal.Float64
		} else {
			// Fallback calculate from big.Int
			intVal := new(big.Float).SetInt(booking.FinalAmount.Int)
			exp := booking.FinalAmount.Exp
			multiplier := new(big.Float).SetFloat64(1.0)
			for i := int32(0); i < -exp; i++ {
				multiplier.Quo(multiplier, big.NewFloat(10))
			}
			res, _ := new(big.Float).Mul(intVal, multiplier).Float64()
			finalAmount = res
		}
	}

	// If booking is free, auto-confirm immediately
	if finalAmount <= 0 {
		wr.Status(http.StatusOK).Json(utils.M{
			"message":  "Booking is free, no payment required.",
			"provider": "free",
			"amount":   0,
		})
		return
	}

	if h.Gateway == nil {
		wr.Status(http.StatusServiceUnavailable).Json(utils.M{
			"error": "Payment gateway is not configured on this server. Please set payment credentials in .env.",
		})
		return
	}

	// Fetch event for metadata
	event, _ := h.Repo.GetEventByID(ctx, booking.EventID)

	orderReq := payment.CreateOrderRequest{
		BookingID:     booking.ID.String(),
		Amount:        finalAmount,
		Currency:      h.Cfg.Payment.Currency,
		CustomerEmail: "",
		CustomerName:  event.Title,
		Receipt:       "bkg_" + booking.ID.String()[:8],
	}

	orderResp, err := h.Gateway.CreateOrder(ctx, orderReq)
	if err != nil {
		logger.Error("Failed to create payment order", "provider", h.Gateway.GetProvider(), "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{
			"error":   "Failed to initialize payment with provider",
			"details": err.Error(),
		})
		return
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"data": orderResp,
	})
}

// VerifyPayment verifies the payment signature/status and issues QR tickets.
func (h *PaymentHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	logger := pkg.NewLogger()
	defer logger.Close()

	if h.Repo == nil || h.Conn == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}

	req := payment.VerifyPaymentRequest{}
	if err := wr.ParseBody(&req); err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "Invalid verification body"})
		return
	}

	if h.Gateway == nil {
		wr.Status(http.StatusServiceUnavailable).Json(utils.M{"error": "Payment gateway is not configured"})
		return
	}

	ctx := r.Context()
	verifyResp, err := h.Gateway.VerifyPayment(ctx, req)
	if err != nil || !verifyResp.Success {
		logger.Error("Payment verification failed", "error", err)
		wr.Status(http.StatusBadRequest).Json(utils.M{
			"error":   "Payment verification failed",
			"details": err.Error(),
		})
		return
	}

	// Payment confirmed! Now update booking status to CONFIRMED and ensure tickets exist
	bookingUUID, err := utils.StringToUUID(req.BookingID)
	if err == nil {
		tx, txErr := h.Conn.Begin(ctx)
		if txErr == nil {
			defer tx.Rollback(ctx) //nolint:errcheck
			qtx := h.Repo.WithTx(tx)

			// Update booking status
			_, _ = qtx.UpdateBookingStatus(ctx, repository.UpdateBookingStatusParams{
				ID:     bookingUUID,
				Status: "CONFIRMED",
			})

			// Check if tickets already generated
			ticketCount, _ := qtx.CountTicketsForBooking(ctx, bookingUUID)
			if ticketCount == 0 {
				booking, _ := qtx.GetBookingByID(ctx, bookingUUID)
				event, _ := qtx.GetEventByID(ctx, booking.EventID)
				ticketItems, _ := qtx.GetBookingTicketsByBookingID(ctx, bookingUUID)

				for _, item := range ticketItems {
					_, _ = ticket.GenerateTicketsForBooking(ctx, qtx, ticket.GenerateTicketsRequest{
						BookingID:    booking.ID.String(),
						EventID:      event.ID.String(),
						EventSlug:    event.Slug,
						TicketTypeID: item.TicketTypeID.String(),
						Quantity:     int(item.Quantity),
						UserID:       booking.UserID.String(),
						HMACSecret:   utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY"),
					}, h.Cfg)
				}
			}

			_ = tx.Commit(ctx)
		}
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"message": "Payment verified successfully and tickets confirmed!",
		"data":    verifyResp,
	})
}

// HandleWebhook receives provider webhook events.
func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	logger := pkg.NewLogger()
	defer logger.Close()

	if h.Gateway == nil {
		wr.Status(http.StatusOK).Json(utils.M{"status": "ignored"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "Failed to read body"})
		return
	}

	signature := r.Header.Get("X-Razorpay-Signature")
	if signature == "" {
		signature = r.Header.Get("Stripe-Signature")
	}

	event, err := h.Gateway.HandleWebhook(r.Context(), body, signature)
	if err != nil {
		logger.Error("Webhook verification error", "error", err)
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": err.Error()})
		return
	}

	logger.Info("Payment webhook processed", "type", event.EventType, "status", event.Status)
	wr.Status(http.StatusOK).Json(utils.M{"status": "ok"})
}
