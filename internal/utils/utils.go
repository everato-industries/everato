package utils

import "os"

func GetEnv(key, d_val string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return d_val
}
