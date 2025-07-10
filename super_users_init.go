package main

import (
	"context"
	"errors"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

func SuperUserInit(cfg *config.Config) error {
	// Initialize the super user with the provided configuration
	// This function should create a super user if it doesn't exist
	// and set up necessary permissions or roles
	logger := pkg.NewLogger()

	// Establish connection to the PostgreSQL database using connection string from environment
	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)

	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return err
	}

	// Initialize the repository with the database connection
	repo := repository.New(conn)

	tx, err := conn.Begin(context.Background())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction", "err", err.Error())
		return err
	}

	for _, su := range cfg.SuperUsers {
		logger.StdoutLogger.Info("Adding user with following details", "username", su.UserName, "email", su.Email)
		_, err = repo.WithTx(tx).CreateSuperUserIfNotExists(
			context.Background(),
			repository.CreateSuperUserIfNotExistsParams{
				Column1: su.UserName,
				Column2: su.Email,
				Column3: su.Password,
			},
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logger.StdoutLogger.Info("Super user already exists", "email", su.Email)
				continue // Not a fatal error, just skip
			}
			logger.StdoutLogger.Error("Error creating super user", "err", err.Error(), "email", su.Email)
			tx.Rollback(context.Background())
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		logger.StdoutLogger.Error("Error committing transaction", "err", err.Error())
		return err
	}

	return nil
}
