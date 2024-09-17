package main

import (
	"context"
	"fbOnboarding/db"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		db.DatabaseConfig
		LoggerConfig
		Server
		JWTConfig
		S3
	}

	LoggerConfig struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	}

	Server struct {
		Port int    `envconfig:"DB_PORT"`
		Host string `envconfig:"DB_HOST"`
	}

	S3 struct {
		AwsEnpoint string `envconfig:"AWS_ENDPOINT"` //http://localhost:4566"
		BucketName string `envconfig:"BUCKET_NAME"`
	}

	JWTConfig struct {
		Secret            string        `envconfig:"SECRET_KEY"`
		ExpiresIn         string        `envconfig:"EXPIRES_IN"`
		ExpiresInDuration time.Duration `envconfig:"-"`
	}
)

var cfg Config

func LoadConfig(_ context.Context) error {
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	drtn, err := time.ParseDuration(cfg.ExpiresIn)
	if err != nil {
		logger.Error().Msg(err.Error())
		return err
	}

	cfg.ExpiresInDuration = drtn

	return nil
}
