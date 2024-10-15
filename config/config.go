package config

import (
	"os"
)

var (
	StoreID    = getEnv("SSLCOMMERZ_STORE_ID", "testbox")
	StorePass  = getEnv("SSLCOMMERZ_STORE_PASS", "qwerty")
	IS_SANDBOX = getEnv("SSLCOMMERZ_IS_SANDBOX", "true")
)

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
