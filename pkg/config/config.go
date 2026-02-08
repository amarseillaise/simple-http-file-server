package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort  int
	ContentDir  string
	TLSCertFile string
	TLSKeyFile  string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		ContentDir:  getEnv("CONTENT_DIR", "./content"),
		TLSCertFile: getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:  getEnv("TLS_KEY_FILE", ""),
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

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
