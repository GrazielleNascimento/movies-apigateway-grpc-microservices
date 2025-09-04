package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server      ServerConfig
	MovieService MovieServiceConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

type MovieServiceConfig struct {
	GRPCAddress string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 10),
		},
		MovieService: MovieServiceConfig{
			GRPCAddress: getEnv("MOVIE_SERVICE_GRPC_ADDRESS", "movies-service:50051"),
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
	if c.MovieService.GRPCAddress == "" {
		log.Fatal("Movie service GRPC address is required")
	}
	return nil
}