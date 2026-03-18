package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load .env first, then .env.local overrides
	_ = godotenv.Load(".env")

	if _, err := os.Stat(".env.local"); err == nil {
		_ = godotenv.Overload(".env.local")
	}
}
