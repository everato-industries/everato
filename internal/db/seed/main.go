package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
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

	eventTitles = []string{
		"Tech Conference 2026", "Music Festival", "Art Exhibition", "Business Summit",
		"Food & Wine Festival", "Sports Tournament", "Science Fair", "Book Fair",
		"Fashion Show", "Charity Gala", "Gaming Convention", "Film Festival",
		"Startup Pitch Night", "Career Fair", "Health & Wellness Expo", "Comedy Night",
		"Photography Workshop", "Dance Competition", "Robotics Challenge", "Cultural Festival",
	}

	eventDescriptions = []string{
		"Join us for an amazing experience with industry leaders and enthusiasts.",
		"An unforgettable event bringing together the best in the field.",
		"Discover, learn, and connect with like-minded individuals.",
		"Experience the future of innovation and creativity.",
		"A celebration of talent, passion, and community.",
		"Don't miss this exclusive opportunity to be part of something special.",
		"Network with professionals and expand your horizons.",
		"Transform your perspective with expert insights and hands-on sessions.",
	}

	eventTypes    = []string{"conference", "workshop", "seminar", "festival", "expo", "meetup", "competition"}
	categories    = []string{"technology", "music", "art", "business", "food", "sports", "education", "entertainment", "health"}
	cities        = []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	states        = []string{"NY", "CA", "IL", "TX", "AZ", "PA", "TX", "CA", "TX", "CA"}
	countries     = []string{"USA", "United States", "US"}
	organizations = []string{"Tech Corp", "Innovation Labs", "Global Events Inc", "Creative Studios", "Future Foundation"}
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

// generateRandomEvent creates a random event with realistic data
func generateRandomEvent(adminID pgtype.UUID) (repository.CreateEventParams, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))

	title := eventTitles[r.Intn(len(eventTitles))]
	description := eventDescriptions[r.Intn(len(eventDescriptions))]

	// Generate slug from title
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	slug = fmt.Sprintf("%s-%d", slug, r.Intn(1000))

	// Generate realistic timestamps
	now := time.Now()
	daysUntilEvent := r.Intn(60) + 1 // Event 1-60 days in the future
	startTime := now.AddDate(0, 0, daysUntilEvent)
	endTime := startTime.Add(time.Hour * time.Duration(r.Intn(8)+2)) // 2-10 hours duration

	bookingStartTime := now.Add(time.Hour * time.Duration(r.Intn(24)))
	bookingEndTime := startTime.Add(-time.Hour * 24) // Booking ends 1 day before event

	totalSeats := int32(r.Intn(450) + 50) // 50-500 seats
	availableSeats := totalSeats

	city := cities[r.Intn(len(cities))]
	state := states[r.Intn(len(states))]
	country := countries[r.Intn(len(countries))]

	// Generate random coordinates (simplified)
	latitude := fmt.Sprintf("%f", 25.0+float64(r.Intn(25)))
	longitude := fmt.Sprintf("-%f", 70.0+float64(r.Intn(50)))

	eventType := eventTypes[r.Intn(len(eventTypes))]
	category := categories[r.Intn(len(categories))]
	organization := organizations[r.Intn(len(organizations))]

	// Generate contact info
	contactEmail := fmt.Sprintf("contact@%s.com", strings.ToLower(strings.ReplaceAll(organization, " ", "")))
	contactPhone := fmt.Sprintf("+1-555-%04d", r.Intn(10000))

	var latitudeNumeric pgtype.Numeric
	if err := latitudeNumeric.Scan(latitude); err != nil {
		return repository.CreateEventParams{}, err
	}

	var longitudeNumeric pgtype.Numeric
	if err := longitudeNumeric.Scan(longitude); err != nil {
		return repository.CreateEventParams{}, err
	}

	return repository.CreateEventParams{
		Title:       title,
		Description: description,
		Slug:        slug,
		Banner:      fmt.Sprintf("https://images.unsplash.com/photo-%d", r.Intn(10000)),
		Icon:        fmt.Sprintf("https://icons.example.com/%s.png", eventType),
		AdminID:     adminID,
		StartTime: pgtype.Timestamptz{
			Time:  startTime,
			Valid: true,
		},
		EndTime: pgtype.Timestamptz{
			Time:  endTime,
			Valid: true,
		},
		Location: pgtype.Text{
			String: fmt.Sprintf("%s Convention Center", city),
			Valid:  true,
		},
		TotalSeats:     totalSeats,
		AvailableSeats: availableSeats,
		Status:         repository.EventStatusCREATED,
		OrganizerName: pgtype.Text{
			String: organization,
			Valid:  true,
		},
		OrganizerEmail: pgtype.Text{
			String: contactEmail,
			Valid:  true,
		},
		OrganizerPhone: pgtype.Text{
			String: contactPhone,
			Valid:  true,
		},
		Organization: pgtype.Text{
			String: organization,
			Valid:  true,
		},
		ContactEmail: pgtype.Text{
			String: contactEmail,
			Valid:  true,
		},
		ContactPhone: pgtype.Text{
			String: contactPhone,
			Valid:  true,
		},
		RefundPolicy: pgtype.Text{
			String: "Full refund available up to 7 days before the event.",
			Valid:  true,
		},
		TermsAndConditions: pgtype.Text{
			String: "By attending this event, you agree to follow all venue rules and regulations.",
			Valid:  true,
		},
		EventType: pgtype.Text{
			String: eventType,
			Valid:  true,
		},
		Category: pgtype.Text{
			String: category,
			Valid:  true,
		},
		MaxTicketsPerUser: pgtype.Int4{
			Int32: int32(r.Intn(5) + 1), // 1-6 tickets per user
			Valid: true,
		},
		BookingStartTime: pgtype.Timestamptz{
			Time:  bookingStartTime,
			Valid: true,
		},
		BookingEndTime: pgtype.Timestamptz{
			Time:  bookingEndTime,
			Valid: true,
		},
		Tags: []string{category, eventType, city},
		WebsiteUrl: pgtype.Text{
			String: fmt.Sprintf("https://www.%s.com", strings.ToLower(strings.ReplaceAll(title, " ", ""))),
			Valid:  true,
		},
		FacebookUrl: pgtype.Text{
			String: fmt.Sprintf("https://facebook.com/%s", slug),
			Valid:  true,
		},
		TwitterUrl: pgtype.Text{
			String: fmt.Sprintf("https://twitter.com/%s", slug),
			Valid:  true,
		},
		InstagramUrl: pgtype.Text{
			String: fmt.Sprintf("https://instagram.com/%s", slug),
			Valid:  true,
		},
		LinkedinUrl: pgtype.Text{
			String: fmt.Sprintf("https://linkedin.com/company/%s", slug),
			Valid:  true,
		},
		VenueName: pgtype.Text{
			String: fmt.Sprintf("%s Convention Center", city),
			Valid:  true,
		},
		AddressLine1: pgtype.Text{
			String: fmt.Sprintf("%d Main Street", r.Intn(9999)+1),
			Valid:  true,
		},
		AddressLine2: pgtype.Text{
			String: fmt.Sprintf("Suite %d", r.Intn(500)+1),
			Valid:  true,
		},
		City: pgtype.Text{
			String: city,
			Valid:  true,
		},
		State: pgtype.Text{
			String: state,
			Valid:  true,
		},
		PostalCode: pgtype.Text{
			String: fmt.Sprintf("%05d", r.Intn(100000)),
			Valid:  true,
		},
		Country: pgtype.Text{
			String: country,
			Valid:  true,
		},
		Latitude:  latitudeNumeric,
		Longitude: longitudeNumeric,
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
	var createdUsers []repository.User
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

		createdUsers = append(createdUsers, user)

		logger.StdoutLogger.Info("Created user",
			"id", user.ID,
			"id_valid", user.ID.Valid,
			"email", user.Email,
			"name", user.FirstName+" "+user.LastName)
	}

	logger.StdoutLogger.Info("User seeding completed", "total_users_created", len(createdUsers))

	// Fetch all admins (super_users) to use as event admins
	admins, err := repo.GetAllAdmins(context.Background())
	if err != nil {
		logger.StdoutLogger.Error("Failed to fetch admins", "error", err.Error())
		log.Fatalf("Failed to fetch admins: %v", err)
	}

	if len(admins) == 0 {
		logger.StdoutLogger.Error("No admins found in database. Please create super users first before seeding events.")
		return
	}

	logger.StdoutLogger.Info("Found admins for event creation", "count", len(admins))

	// Create events
	logger.StdoutLogger.Info("Creating demo events")

	// Create 30 events with diverse realistic data
	totalEvents := 30
	var createdEventsCount int
	for range totalEvents {
		// Small delay to ensure different timestamp for random seed
		time.Sleep(time.Millisecond * 10)

		// Randomly select an admin as event organizer
		admin := admins[rand.Intn(len(admins))]

		// Verify the admin ID is valid
		if !admin.ID.Valid {
			logger.StdoutLogger.Error("Selected admin has invalid ID", "email", admin.Email)
			continue
		}

		// Generate random event data
		eventParams, err := generateRandomEvent(admin.ID)
		if err != nil {
			logger.StdoutLogger.Error("Error generating event data", "error", err.Error())
			continue
		}

		// Create event in the database
		event, err := repo.CreateEvent(context.Background(), eventParams)
		if err != nil {
			logger.StdoutLogger.Error("Failed to create event", "title", eventParams.Title, "error", err.Error())
			continue
		}

		createdEventsCount++

		logger.StdoutLogger.Info("Created event",
			"id", event.ID,
			"title", event.Title,
			"admin", admin.Email,
			"start_time", event.StartTime.Time.Format(time.RFC3339))
	}

	logger.StdoutLogger.Info("Event seeding completed", "total_events_created", createdEventsCount)
	logger.StdoutLogger.Info("Seeding completed successfully",
		"total_users_created", len(createdUsers),
		"total_events_created", createdEventsCount)
}
