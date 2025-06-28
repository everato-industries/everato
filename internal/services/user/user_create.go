package user

import (
	"errors"
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
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
		logger.StdoutLogger.Error("Error Finding user") // Log the error
		tx.Rollback(wr.R.Context())                     // Rollback the transaction
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

	// Return the actual user data
	wr.Status(http.StatusCreated).Json(
		utils.M{
			"success": true,
			"message": "User registration endpoint reached successfully",
			"data":    user,
		},
	)
}
