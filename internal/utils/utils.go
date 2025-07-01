package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// GetEnv retrieves the value of the environment variable named by key.
// if the value is not set in the environment
// then it returns the default value it was passed to
func GetEnv(key, d_val string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return d_val
}

// This method accepts a string value of an UUID
// then tries to parse that into an actual UUID
//
// returns error with empty UUID or nil with parsed UUID
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

// This method accepts a string value of a time
// in RFC3339 / ISO 8061 format
// then tries to parse that into an actual time
//
// returns error with empty Timestampz or nil with parsed Timestampz
func StringToTime(s string) (pgtype.Timestamptz, error) {
	t := pgtype.Timestamptz{}
	err := t.Scan(s)

	// In case of error in parsing return an empty time struct
	if err != nil {
		return pgtype.Timestamptz{}, err
	}

	// Other wise return the time
	return t, nil
}

// Returns a pgtype.Text version of this string
func StringToText(s string) (pgtype.Text, error) {
	t := pgtype.Text{}

	err := t.Scan(s)

	// In case of error in parsing return an empty text struct
	if err != nil {
		return pgtype.Text{}, err
	}

	return t, nil
}

// This method accepts the title of an event and generates a slug based on that
//
// There are some rules for the slug:
// // 1. It should be lower case
// // 2. It should be unique
// // It returns the slug as a string or an error if it fails to generate
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
