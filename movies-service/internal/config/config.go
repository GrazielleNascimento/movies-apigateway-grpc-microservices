package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GRPC     GRPCConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	ConnectionString string
	DatabaseName     string
	MaxPoolSize      int
}

type GRPCConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 10),
		},
		Database: DatabaseConfig{
			ConnectionString: getEnv("MONGODB_URI", "mongodb://mongodb:27017"),
			DatabaseName:     getEnv("DATABASE_NAME", "movies_db"),
			MaxPoolSize:      getEnvAsInt("MAX_POOL_SIZE", 10),
		},
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.ConnectionString == "" {
		log.Fatal("Database connection string is required")
	}
	if c.Database.DatabaseName == "" {
		log.Fatal("Database name is required")
	}
	if c.GRPC.Port == "" {
		log.Fatal("GRPC port is required")
	}
	return nil
}
