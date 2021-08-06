package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv reads .env file and loads all variables as environment variables.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// ReadEnv reads an environment variable with a default value if it doesn't
// exist.
func ReadEnv(key string, _default string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return _default
	}
	return value
}

// MustReadEnv reads an environment variable and calls log.Panic if it doesn't
// exist.
func MustReadEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Panic("invalid key " + key)
	}

	return value
}
