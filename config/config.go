package config

import (
	"clearWayTest/pkg/helpers"
	"errors"
	"fmt"
	"log"
	"os"
)

type (
	Config struct {
		App
		HTTP
		PG
	}

	App struct {
		Name    string `env:"APP_NAME"`
		Version string `env:"APP_VERSION"`
	}

	HTTP struct {
		Port     string `env:"HTTP_PORT"`
		CertName string `json:"HTTP_CERT_NAME"`
		CertKey  string `json:"HTTP_CERT_KEY"`
	}

	PG struct {
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `env:"DB_NAME"`
		SSLMode  string `env:"DB_SSL_MODE"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := helpers.ParseEnvFile()

	if err != nil {
		log.Fatalf("Can't parse env file: %s", err)
	}

	cfg.App = App{
		Name:    os.Getenv("APP_NAME"),
		Version: os.Getenv("APP_VERSION"),
	}

	cfg.HTTP = HTTP{
		Port:     fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
		CertKey:  os.Getenv("HTTP_CERT_KEY"),
		CertName: os.Getenv("HTTP_CERT_NAME"),
	}
	dbPort := os.Getenv("DB_PORT")

	if dbPort == "" {
		return nil, errors.New("db port is empty")
	}

	cfg.PG = PG{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	return cfg, nil
}
