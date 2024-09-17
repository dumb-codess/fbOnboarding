package main

import (
	"context"
	"fbOnboarding/db"
	s3base "fbOnboarding/s3"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func main() {
	// Initialize logger first
	logger = zerolog.New(os.Stderr).With().CallerWithSkipFrameCount(2).Caller().Timestamp().Logger()

	logger.Info().Msg("setting logger level........")

	ctx := context.Background()

	logger.Info().Msg("Loading the configuration........")
	if err := LoadConfig(ctx); err != nil {
		logger.Fatal().Msgf("Failed to Load config: %v", err.Error())
	}
	logger.Info().Msg("Configuration loaded!")

	logger.Info().Msg("setting up the db....")
	if err := db.Init(ctx, cfg.DatabaseConfig); err != nil {
		logger.Fatal().Msgf("Failed to Initialize db: %v", err.Error())
	}
	logger.Info().Msg("db setup finish!")

	logger.Info().Msg("setting up s3.....")
	if err := s3Setup(); err != nil {
		logger.Fatal().Msgf("failed to set up  s3: %v", err)
	}
	logger.Info().Msg("s3 setup finish!")

	controller, err := NewController()
	if err != nil {
		logger.Fatal().Msgf("Failed to Load config: %v", err.Error())
	}

	listenAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := http.Server{
		Addr:              listenAddr,
		Handler:           router(&cfg, controller),
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger.Info().Msgf("listening on :%v....", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal().Msg(err.Error())
	}

}

func s3Setup() error {
	s3, err := s3base.NewS3Client(cfg.AwsEnpoint)
	if err != nil {
		return err
	}

	if err := s3.CreateBucket(cfg.BucketName); err != nil {
		logger.Fatal().Msgf("failed to create bucket: %v", err)
		return err
	}

	return nil
}
