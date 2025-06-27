package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// User data for seeding
var (
	firstNames = []string{
		"John", "Jane", "Michael", "Emma", "William", "Olivia", "James", "Sophia",
		"Robert", "Ava", "Thomas", "Isabella", "David", "Mia", "Richard", "Amelia",
		"Charles", "Harper", "Joseph", "Evelyn", "Daniel", "Abigail", "Matthew", "Emily",
		"Anthony", "Elizabeth", "Christopher", "Sofia", "Andrew", "Ella",
	}

	lastNames = []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis", "Garcia",
		"Rodriguez", "Wilson", "Martinez", "Anderson", "Taylor", "Thomas", "Hernandez", "Moore",
		"Martin", "Jackson", "Thompson", "White", "Lopez", "Lee", "Gonzalez", "Harris",
		"Clark", "Lewis", "Young", "Walker", "Hall", "Allen", "King", "Wright",
	}

	domains = []string{
		"gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "protonmail.com",
		"icloud.com", "aol.com", "mail.com", "zoho.com", "yandex.com",
	}
)

// hashPassword creates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// generateRandomUser creates a random user with realistic data
func generateRandomUser() (repository.CreateUserParams, error) {
	// Create a new random source for each user to enhance randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))

	firstName := firstNames[r.Intn(len(firstNames))]
	lastName := lastNames[r.Intn(len(lastNames))]

	// Generate email with firstName.lastName format
	domainName := domains[r.Intn(len(domains))]
	randomNum := r.Intn(999)
	email := fmt.Sprintf("%s.%s%03d@%s", firstName, lastName, randomNum, domainName)
	email = strings.ToLower(email) // Converting to lowercase

	// Generate a password (typically in production you'd have stronger password requirements)
	password := fmt.Sprintf("%s%s!%d", firstName, lastName, r.Intn(100))

	// Hash the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return repository.CreateUserParams{}, err
	}

	return repository.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
	}, nil
}

func main() {
	// Initialize the logger
	logger := pkg.NewLogger()
	defer logger.Close()
	logger.StdoutLogger.Info("Starting database seed")

	// Initialize database connection
	dbURL := utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?sslmode=disable")
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		logger.StdoutLogger.Error("Failed to connect to database", "error", err.Error())
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create the repository instance
	repo := repository.New(conn)

	// Create users
	logger.StdoutLogger.Info("Creating demo users")

	// Create 20 users with diverse realistic data
	totalUsers := 20
	for range totalUsers {
		// Small delay to ensure different timestamp for random seed
		time.Sleep(time.Millisecond * 10)
		// Generate random user data
		userParams, err := generateRandomUser()
		if err != nil {
			logger.StdoutLogger.Error("Error generating user data", "error", err.Error())
			continue
		}

		// Create user in the database
		user, err := repo.CreateUser(context.Background(), userParams)
		if err != nil {
			logger.StdoutLogger.Error("Failed to create user", "email", userParams.Email, "error", err.Error())
			continue
		}

		logger.StdoutLogger.Info("Created user",
			"id", user.ID,
			"email", user.Email,
			"name", user.FirstName+" "+user.LastName)
	}

	logger.StdoutLogger.Info("Seeding completed successfully", "total_users_created", totalUsers)
}
