package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github/wry-0313/exchange/pkg/validator"

	"github.com/joho/godotenv"
)

const (
	keyDBHost     = "DB_HOST"
	keyDBPort     = "DB_PORT"
	keyDBName     = "DB_NAME"
	keyDBUser     = "DB_USER"
	keyDBPassword = "DB_PASSWORD"

	keyEnv             = "ENV"
	keyServerPort      = "SERVER_PORT"
	keyJWTSecret       = "JWT_SIGNING_KEY"
	keyJWTExpiration   = "JWT_EXPIRATION"
	keyInternalNetwork = "INTERNAL_NETWORK"

	keyKafkaBrokers = "KAFKA_BROKERS"

	ProdEnv = "PRODUCTION"
	DevEnv  = "DEVELOPMENT"
)

type Config struct {
	DB            DatabaseConfig
	ServerPort    string
	JwtSecret     string
	JwtExpiration int
	KafkaBrokers  []string
}

func Load(file string) (*Config, error) {
	env := os.Getenv(keyEnv)
	if env != ProdEnv {
		// Load .env file if in development
		err := godotenv.Load(file)
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting database config: %w", err)
	}

	serverPort := os.Getenv(keyServerPort)
	jwtSecret := os.Getenv(keyJWTSecret)
	jwtExpirationStr := os.Getenv(keyJWTExpiration)

	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration value: %w", err)
	}

	broker := os.Getenv(keyKafkaBrokers)
	KafkaBrokers := []string{broker}

	return &Config{
		DB:            databaseConfig,
		ServerPort:    serverPort,
		JwtSecret:     jwtSecret,
		JwtExpiration: jwtExpiration,
		KafkaBrokers: KafkaBrokers,
	}, nil
}

// DatabaseConfig encapsulates all the config values for the database.
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// Validate checks that all values are properly loaded into the database config.
func (dbConfig *DatabaseConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		return fmt.Errorf("missing database env var: %v", err)
	}
	return nil
}

func getDatabaseConfig() (DatabaseConfig, error) {
	databaseConfig := DatabaseConfig{
		Host:     os.Getenv(keyDBHost),
		Port:     os.Getenv(keyDBPort),
		Name:     os.Getenv(keyDBName),
		User:     os.Getenv(keyDBUser),
		Password: os.Getenv(keyDBPassword),
	}
	log.Printf("databaseConfig: %v\n", databaseConfig)

	// This allows running tests from outside the docker network assuming your local
	// development environment has ports exposed
	if os.Getenv(keyInternalNetwork) == "false" {
		databaseConfig.Host = "localhost"
	}

	// validate all db params are available
	if err := databaseConfig.Validate(); err != nil {
		return DatabaseConfig{}, err
	}

	return databaseConfig, nil
}