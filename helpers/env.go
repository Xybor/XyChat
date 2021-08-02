package helpers

import (
	"errors"
	"log"

	"github.com/joho/godotenv"
)

var env map[string]string

func LoadEnv() {
	e, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}

	env = e
}

func ReadEnv(key string, _default string) string {
	value, ok := env[key]
	if !ok {
		return _default
	}
	return value
}

func MustReadEnv(key string) (string, error) {
	value, ok := env[key]
	if !ok {
		return "", errors.New("KeyError: " + key)
	}

	return value, nil
}
