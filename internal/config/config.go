package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	IsDebug       bool `envconfig:"IS_DEBUG" default:"false"`
	IsDevelopment bool `envconfig:"IS_DEVELOP" default:"false"`
	Listen        struct {
		BindIP string `envconfig:"BIND_IP" default:"0.0.0.0"`
		Port   string `envconfig:"PORT" default:"8080"`
	}
	Server struct {
		ReadTimeout  int    `envconfig:"READ_TIMEOUT" default:"10"`
		WriteTimeout int    `envconfig:"WRITE_TIMEOUT" default:"10"`
		IdleTimeout  int    `envconfig:"IDLE_TIMEOUT" default:"100"`
		ServerType   string `envconfig:"SERVER_TYPE" default:"http"`
	}
	AppConfig struct {
		LogLevel string `envconfig:"LOGLEVEL" default:"debug"`
	}
	Database struct {
		Host     string `envconfig:"DB_HOST" default:"db"`
		DBName   string `envconfig:"DB_NAME" default:"postgres"`
		Port     string `envconfig:"DB_PORT" default:"5432"`
		User     string `envconfig:"DB_USER" default:"postgres"`
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
		SslMode  string `envconfig:"SSL_MODE" default:"disable"`
	}
	Bucket struct {
		IPLimit             int `envconfig:"IP_LIMIT" default:"1000"`
		LoginLimit          int `envconfig:"LOGIN_LIMIT" default:"10"`
		PasswordLimit       int `envconfig:"PASSWORD_LIMIT" default:"100"`
		ResetBucketInterval int `envconfig:"RESET_BUCKET_INTERVAL" default:"60"`
	}
}

func LoadAll() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("envconfig.Process: %w", err)
	}

	return &cfg, nil
}
