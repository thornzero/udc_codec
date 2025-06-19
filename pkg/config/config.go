package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath  string
	DataDir string
	Port    string
}

var (
	cfg  *Config
	once sync.Once
)

// Load once, singleton pattern
func Load() *Config {
	once.Do(func() {
		// Load .env file if present
		_ = godotenv.Load()

		cfg = &Config{
			DBPath:  getEnv("DB_PATH", "tags.db"),
			DataDir: getEnv("DATA_DIR", "data"),
			Port:    getEnv("SERVER_PORT", "8080"),
		}
	})
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (c *Config) Path(subpath string) string {
	return fmt.Sprintf("%s/%s", c.DataDir, subpath)
}
