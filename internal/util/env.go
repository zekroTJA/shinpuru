package util

import (
	"os"
	"strings"
)

const envPrefix = "SP_"

// GetEnv returns a value from environment variables
// or the given default value, if not existent.
func GetEnv(key, def string) string {
	res := os.Getenv(envPrefix + strings.ToUpper(key))
	if res == "" {
		return def
	}
	return res
}
