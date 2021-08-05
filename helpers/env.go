package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

var env map[string]string

// LoadEnv reads .env file and loads all variables as environment variables.
func LoadEnv() {
	e, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}

	env = e
}

// ReadEnv reads an environment variable with a default value if it doesn't
// exist.
func ReadEnv(key string, _default string) string {
	value, ok := env[key]
	if !ok {
		return _default
	}
	return value
}

// MustReadEnv reads an environment variable and calls log.Panic if it doesn't
// exist.
func MustReadEnv(key string) string {
	value, ok := env[key]
	if !ok {
		log.Panic("invalid key " + key)
	}

	return value
}
