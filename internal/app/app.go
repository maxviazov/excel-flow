package app

import (
	"context"
	"os"

	"github.com/rs/zerolog"

	"github.com/maxviazov/excel-flow/internal/config"
)

type Pipeline interface {
	Run(ctx context.Context, cfg *config.Config, log zerolog.Logger) error
}

type mockPipeline struct{}

func (m *mockPipeline) Run(ctx context.Context, cfg *config.Config, log zerolog.Logger) error {
	log.Info().Msg("Mock pipeline executed")
	return nil
}

func Run() error {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		return err
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	pipeline := &mockPipeline{}

	logger.Info().Msg("🚀 Starting Excel-Flow")

	ctx := context.Background()
	if err := pipeline.Run(ctx, cfg, logger); err != nil {
		logger.Error().Err(err).Msg("Pipeline execution failed")
		return err
	}

	logger.Info().Msg("✅ Excel-Flow finished successfully")
	return nil
}
