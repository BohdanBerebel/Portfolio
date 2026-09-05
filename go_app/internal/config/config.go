package config

import (
	"errors"
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	JWTSecret  string
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)

	if value == "" {
		return "", errors.New(
			"missing required environment variable: " + name,
		)
	}

	return value, nil
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbHost, err := requiredEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	dbPort, err := requiredEnv("DB_PORT")
	if err != nil {
		return Config{}, err
	}

	dbName, err := requiredEnv("DB_NAME")
	if err != nil {
		return Config{}, err
	}

	dbUser, err := requiredEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}

	dbPassword, err := requiredEnv("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:       port,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		JWTSecret:  jwtSecret,
	}, nil
}
