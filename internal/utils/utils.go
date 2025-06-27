package utils

import "os"

// GetEnv retrieves the value of the environment variable named by key.
// if the value is not set in the environment
// then it returns the default value it was passed to
func GetEnv(key, d_val string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return d_val
}
