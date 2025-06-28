package utils

import (
	"os"

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
