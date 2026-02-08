package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  int
	ContentDir  string
	TLSCertFile string
	TLSKeyFile  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error: .env file not found")
	}

	return &Config{
		ServerPort:  getEnvAsInt("SERVER_PORT"),
		ContentDir:  os.Getenv("CONTENT_DIR"),
		TLSCertFile: os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:  os.Getenv("TLS_KEY_FILE"),
	}
}

// TLSEnabled returns true if both TLS certificate and key files are specified
func (c *Config) TLSEnabled() bool {
	return c.TLSCertFile != "" && c.TLSKeyFile != ""
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string) int {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 8443 // default port
	}
	return value
}
