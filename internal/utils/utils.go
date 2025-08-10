// Package utils provides utility functions for the Everato application.
// It includes helpers for environment variables, data type conversions, and string manipulation.
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// GetEnv retrieves the value of the environment variable named by key.
// If the value is not set in the environment, it returns the provided default value.
//
// Parameters:
//   - key: The name of the environment variable to retrieve
//   - d_val: The default value to return if the environment variable is not set
//
// Returns:
//   - The value of the environment variable, or the default value if not set
func GetEnv(key, d_val string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return d_val
}

// StringToUUID converts a string representation of a UUID to a pgtype.UUID type.
// This is useful for converting UUID strings from requests to database-compatible types.
//
// Parameters:
//   - s: The string representation of a UUID (e.g., "550e8400-e29b-41d4-a716-446655440000")
//
// Returns:
//   - A pgtype.UUID object if successful
//   - An error if the string cannot be parsed as a valid UUID
func StringToUUID(s string) (pgtype.UUID, error) {
	uuid := pgtype.UUID{}
	err := uuid.Scan(s)

	// In case of error in parsing return an empty uuid struct
	if err != nil {
		return pgtype.UUID{}, err
	}

	// Other wise return the UUID
	return uuid, nil
}

// StringToTime converts a string representation of time in various formats
// to a pgtype.Timestamptz type for use with PostgreSQL.
//
// Supported formats include:
// - RFC3339 (e.g., "2023-04-01T15:30:00Z")
// - ISO8601 with timezone (e.g., "2023-04-01T15:30:00+05:30")
// - ISO8601 without timezone (e.g., "2023-04-01T15:30:00")
// - Date only (e.g., "2023-04-01")
// - Common formats with named timezones (e.g., "2023-04-01 15:30:00 GMT")
//
// Parameters:
//   - s: The string representation of time
//
// Returns:
//   - A pgtype.Timestamptz object if successful
//   - An error if the string cannot be parsed as a valid timestamp
func StringToTime(s string) (pgtype.Timestamptz, error) {
	t := pgtype.Timestamptz{}

	// Define a list of time formats to try
	formats := []string{
		time.RFC3339,                // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05",       // ISO8601 without timezone
		"2006-01-02 15:04:05",       // Common datetime format
		"2006-01-02 15:04:05 MST",   // With named timezone
		"2006-01-02 15:04:05 -0700", // With numeric timezone
		"2006-01-02",                // Date only
		time.RFC1123,                // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,               // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC822,                 // "02 Jan 06 15:04 MST"
		time.RFC822Z,                // "02 Jan 06 15:04 -0700"
	}

	var parsedTime time.Time
	var err error
	var lastError error

	// Try each format until we find one that works
	for _, format := range formats {
		parsedTime, err = time.Parse(format, s)
		if err == nil {
			// If parsing succeeds, break the loop
			break
		}
		lastError = err
	}

	// If none of the formats worked, return the last error
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("could not parse time '%s' in any supported format: %w", s, lastError)
	}

	// Set the Timestamptz value using the parsed time
	t.Time = parsedTime
	if err != nil {
		return pgtype.Timestamptz{}, fmt.Errorf("failed to convert parsed time to pgtype.Timestamptz: %w", err)
	}

	return t, nil
}

func RFCTimeToTimeStampZ(t time.Time) (pgtype.Timestamptz, error) {
	// Create a new Timestamptz object
	tz := pgtype.Timestamptz{}

	// Set the time value
	tz.Time = t

	// Other wise return the time
	return tz, nil
}

// StringToText converts a string to a pgtype.Text value for PostgreSQL compatibility.
// This is useful when passing strings to database queries that expect text types.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - A pgtype.Text object representing the string
//   - An error if conversion fails
func StringToText(s string) (pgtype.Text, error) {
	t := pgtype.Text{}

	err := t.Scan(s)

	// In case of error in parsing return an empty text struct
	if err != nil {
		return pgtype.Text{}, err
	}

	return t, nil
}

// GenerateSlug creates a URL-friendly slug from a given title string.
// The generated slug follows these rules:
// 1. It is converted to lowercase
// 2. Spaces are replaced with underscores
// 3. Only alphanumeric characters and underscores are retained
//
// Parameters:
//   - title: The original title to convert to a slug
//
// Returns:
//   - The generated slug string if successful
//   - An error if the resulting slug would be empty (e.g., title contains only special characters)
func GenerateSlug(title string) (string, error) {
	// Convert the title to lower case
	title_lower := strings.ToLower(title)
	title_wo_space := strings.TrimSpace(title_lower)
	title_wo_space = strings.ReplaceAll(title_wo_space, " ", "_")

	// Replace any special characters with an empty string
	slug := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		return -1 // Remove the character
	}, title_wo_space)

	// If the slug is empty, return an error
	if slug == "" {
		return "", errors.New("The title you provided was too short and also consists more special characters than allowed")
	}

	return slug, nil
}

// Sha256 computes the SHA256 hash of a given string and returns its hexadecimal representation.
//
// Parameters:
//   - String to be hashed
//
// Returns:
//   - Hashed version of the string it had
func Sha256(s string) string {
	// Create a new SHA256 hash
	h := sha256.New()

	// Write the string to the hash
	h.Write([]byte(s))

	// Return the hex representation of the hash
	return hex.EncodeToString(h.Sum(nil))
}

// BcryptHash hashes a string using bcrypt and returns the hashed value.
//
// Parameters:
//   - String to hash
//
// Returns:
//   - Generated hash or and empty string in case of error
//   - Nil or the error in case of error
func BcryptHash(s string) (string, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
