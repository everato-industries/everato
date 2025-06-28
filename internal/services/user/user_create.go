package user

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/services/mailer"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

func CreateUser(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the function exits

	user_dto := &CreateUserDTO{}
	err := wr.ParseBody(user_dto)
	if err != nil {
		logger.StdoutLogger.Error(
			"Error parsing body",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"error": err.Error(),
			},
		)
		return
	}

	// Validate whether the sent data is valid or not
	if err := user_dto.Validate(); err != nil {
		logger.StdoutLogger.Error(
			"Error parsing body",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"error":   err.Error(),
				"message": "Invalid user data provided",
			},
		)
		return
	}

	// Create the actual user in the database
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error(
			"Error starting a transaction",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		logger.FileLogger.Error(
			"Error starting a transaction",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"error":   "Failed to begin transaction",
				"message": err.Error(),
			},
		)
		return
	}

	// Find if the user with the same email already exists or not
	_, err = repo.WithTx(tx).GetUserByEmail(
		wr.R.Context(),
		user_dto.Email,
	)
	if err == nil {
		// If no error, it means the user exists
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusConflict).Json(
			utils.M{
				"error":   "User with this email already exists",
				"message": "Please try with a different email address",
			},
		)
		return
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// Check if the error is anything other than "user not found"
		// If it's a different error, it indicates a database issue
		logger.StdoutLogger.Error("Error Finding user", "err", err.Error()) // Log the error
		tx.Rollback(wr.R.Context())                                         // Rollback the transaction
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Failed to check if user already exists",
				"error":   err.Error(),
			},
		)
		return
	}

	// Hash the password inside the DTO object
	err = user_dto.HashPassword()
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
	}

	// Create the user in the database
	user, err := repo.WithTx(tx).CreateUser(
		wr.R.Context(),
		user_dto.ToCreteUserParams(),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error creating user", "err", err.Error()) // Log the error
		tx.Rollback(wr.R.Context())                                          // Rollback the transaction
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"error":   "Failed to create user",
				"message": err.Error(),
			},
		)
		return
	}

	tx.Commit(wr.R.Context()) // Commit the transaction

	// =================================================================
	// GOROUTINE starts
	//
	// Send the verification mail on a seperate thread
	go func() {
		logger := pkg.NewLogger()
		defer logger.Close()

		// Initialize the payload
		type EmailPayload struct {
			VerificationLink string // Verification link
		}

		// Parse the mail template here
		tpl_path := filepath.Join("templates", "mail", "verify_email.html")
		tpl, err := template.ParseFiles(tpl_path)
		if err != nil {
			logger.StdoutLogger.Error("Error parsing the mail template", "err", err.Error())
			logger.FileLogger.Error("Error parsing the mail template", "err", err.Error())
			return
		}

		// Build the verification URL
		verification_url := fmt.Sprintf("%s/auth/verify-email", utils.GetEnv("API_URL", "http://localhost:8080/api/v1"))

		var body bytes.Buffer // Body data for the email
		err = tpl.Execute(&body, EmailPayload{
			VerificationLink: verification_url,
		})
		if err != nil {
			logger.StdoutLogger.Error("Error parsing the mail template", "err", err.Error())
			logger.FileLogger.Error("Error parsing the mail template", "err", err.Error())
			return
		}

		smtp_port, err := strconv.Atoi(utils.GetEnv("SMTP_PORT", "587"))
		if err != nil {
			logger.StdoutLogger.Error("Error parsing SMTP server PORT", "err", err.Error())
			logger.FileLogger.Error("Error parsing SMTP server PORT", "err", err.Error())
			return
		}

		// Instantiate the mailer service
		mail_service := mailer.NewMailService(&mailer.MailerParameters{
			To:      user.Email,
			Subject: "Verify your account!!",
			Body:    &body,
			Options: &mailer.MailerOptions{
				Host:        utils.GetEnv("SMTP_HOST", "smtp.gmail.com"),
				Port:        uint16(smtp_port),
				SenderEmail: utils.GetEnv("SMTP_EMAIL", "dev.bosepiush@gmail.com"),
				AppPassword: utils.GetEnv("SMTP_PASSWORD", "SUPER_SECRET_PASSWORD"),
			},
		})

		logger.StdoutLogger.Info("Sending verifcation email", "user_email", user.Email)

		res, err := mail_service.SendEmail(wr) // Send the email
		if res != mailer.MailerSuccess && err != nil {
			logger.StdoutLogger.Error("Error sending verification email", "err", err.Error())
			logger.FileLogger.Error("Error sending verification email", "err", err.Error())
			return
		}
	}() // GOROUTINE finishes
	// =================================================================

	// Return the actual user data
	wr.Status(http.StatusCreated).Json(
		utils.M{
			"success": true,
			"message": "User registration endpoint reached successfully",
			"note":    "Kindly check your email and verify your account withing 24 hours",
			"data":    user,
		},
	)
}
