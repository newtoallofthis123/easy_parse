package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	ConnString   string
	Port         string
	GeminiAPIKey string
}

// Reads the .env file and returns an Env struct.
func ReadEnv() Env {
	if os.Getenv("GIN_MODE") == "release" {
		godotenv.Load(".env.prod")
	} else {
		godotenv.Load(".env")
	}

	return Env{
		ConnString:   getEnv("DATABASE_URL"),
		Port:         getEnv("PORT"),
		GeminiAPIKey: getEnv("GOOGLE_API_KEY"),
	}
}

// Returns the value of the given env var name.
func getEnv(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("Env var %s not found", name))
	}
	return val
}
