package helpers

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv reads .env file and loads all variables as environment variables.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func ReadEnvDefault(key, _default string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return _default
	}
	return value
}

// ReadEnv reads an environment variable with a default value if it doesn't
// exist.
func ReadEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.New("non-existed environment variables")
	}
	return value, nil
}

// MustReadEnv reads an environment variable and calls log.Fatal if it doesn't
// exist.
func MustReadEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("non-existed key " + key)
	}

	return value
}
