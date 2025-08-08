package admin

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/services/mailer"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func SendVerificationEmail(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract the admin ID from URL parameters
	vars := mux.Vars(wr.R)
	targetAdminID := vars["id"]
	if targetAdminID == "" {
		logger.StdoutLogger.Error("Admin ID not provided in URL parameters")
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Missing admin ID",
				"error":   "Admin ID must be provided in the URL",
			},
		)
		return
	}

	// Get the current admin ID from the request context
	currentAdminID, ok := wr.R.Context().Value("admin_id").(string)
	if !ok {
		logger.StdoutLogger.Error("Failed to get admin ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required",
				"error":   "Admin ID not found in request context",
			},
		)
		return
	}

	// Start a transaction for consistency and atomicity
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Convert the admin IDs to UUIDs
	currentAdminUUID, err := utils.StringToUUID(currentAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid current admin UUID", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	targetAdminUUID, err := utils.StringToUUID(targetAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid target admin UUID", "adminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Get current admin details to check permissions
	currentAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), currentAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve current admin data", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed",
				"error":   "Admin account not found",
			},
		)
		return
	}

	// Get target admin details
	targetAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), targetAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve target admin data", "adminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Target admin not found",
				"error":   "The admin account you're trying to send verification email to doesn't exist",
			},
		)
		return
	}

	// Check permissions for sending verification email
	// 1. SUPER_ADMINs can send verification emails to any account
	// 2. Admins with MANAGE_USERS permission can send verification emails to non-SUPER_ADMIN accounts
	// 3. Admins can send verification emails to themselves
	isSelfAction := currentAdminID == targetAdminID
	isSuperAdmin := currentAdmin.Role == "SUPER_ADMIN"
	hasManageUsersPermission := false

	// Check if the current admin has MANAGE_USERS permission
	for _, perm := range currentAdmin.Permissions {
		if perm == "MANAGE_USERS" {
			hasManageUsersPermission = true
			break
		}
	}

	// Check if the current admin is allowed to send verification email to the target admin
	hasSendPermission := isSelfAction || isSuperAdmin ||
		(hasManageUsersPermission && targetAdmin.Role != "SUPER_ADMIN")

	if !hasSendPermission {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Admin lacks permission to send verification email",
			"currentAdminID", currentAdminID,
			"targetAdminID", targetAdminID)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "You don't have permission to send verification email to this admin account",
			},
		)
		return
	}

	// Commit the transaction now that we've verified permissions
	// We don't need to keep the transaction open for email sending
	if err = tx.Commit(wr.R.Context()); err != nil {
		logger.StdoutLogger.Error("Failed to commit transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to complete database operation",
			},
		)
		return
	}

	// Start a goroutine to send the email without blocking the response
	go func() {
		// Create a new logger for this goroutine
		emailLogger := pkg.NewLogger()
		defer emailLogger.Close()

		// Initialize the payload for the email template
		type EmailPayload struct {
			AdminName        string
			AdminEmail       string
			VerificationLink string
		}

		// Parse the mail template
		tpl, err := pkg.GetTemplate("templates/mail/admin_verify_email.html")
		if err != nil {
			// If the template doesn't exist, fall back to a generic template
			tpl, err = pkg.GetTemplate("templates/mail/verify_email.html")
			if err != nil {
				emailLogger.StdoutLogger.Error("Error loading email template", "err", err.Error())
				emailLogger.FileLogger.Error("Error loading email template", "err", err.Error())
				return
			}
		}

		// Build the verification URL
		verificationURL := fmt.Sprintf(
			"%s/admin/verify/%s",
			utils.GetEnv("API_URL", "http://localhost:8080/api/v1"),
			targetAdmin.ID.String(),
		)

		// Prepare email body
		var body bytes.Buffer
		err = tpl.Execute(&body, EmailPayload{
			AdminName:        targetAdmin.Name,
			AdminEmail:       targetAdmin.Email,
			VerificationLink: verificationURL,
		})
		if err != nil {
			emailLogger.StdoutLogger.Error("Error executing email template", "err", err.Error())
			emailLogger.FileLogger.Error("Error executing email template", "err", err.Error())
			return
		}

		// Get SMTP port from environment or use default
		smtpPort, err := strconv.Atoi(utils.GetEnv("SMTP_PORT", "587"))
		if err != nil {
			emailLogger.StdoutLogger.Error("Error parsing SMTP server PORT", "err", err.Error())
			emailLogger.FileLogger.Error("Error parsing SMTP server PORT", "err", err.Error())
			return
		}

		// Instantiate the mailer service
		mailService := mailer.NewMailService(&mailer.MailerParameters{
			To:      targetAdmin.Email,
			Subject: "Verify your admin account - Everato Platform",
			Body:    &body,
			Options: &mailer.MailerOptions{
				Host:        utils.GetEnv("SMTP_HOST", "smtp.gmail.com"),
				Port:        uint16(smtpPort),
				SenderEmail: utils.GetEnv("SMTP_EMAIL", "dev.bosepiush@gmail.com"),
				AppPassword: utils.GetEnv("SMTP_PASSWORD", "SUPER_SECRET_PASSWORD"),
			},
		})

		emailLogger.StdoutLogger.Info("Sending verification email", "admin_email", targetAdmin.Email)

		// Send the email
		res, err := mailService.SendEmail(nil) // We're using nil because we're in a goroutine
		if res != mailer.MailerSuccess || err != nil {
			emailLogger.StdoutLogger.Error("Error sending verification email", "err", err.Error())
			emailLogger.FileLogger.Error("Error sending verification email", "err", err.Error())
			return
		}

		emailLogger.StdoutLogger.Info("Verification email sent successfully", "admin_email", targetAdmin.Email)
	}()

	// Return a success response immediately without waiting for the email to be sent
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Verification email queued for sending",
			"admin": utils.M{
				"id":       targetAdmin.ID.String(),
				"email":    targetAdmin.Email,
				"username": targetAdmin.Username,
				"name":     targetAdmin.Name,
			},
		},
	)
}
