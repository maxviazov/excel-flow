package app

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/maxviazov/excel-flow/internal/config"
	"github.com/maxviazov/excel-flow/internal/log"
)

type Pipeline interface {
	Run(ctx context.Context, cfg *config.Config, log zerolog.Logger) error
}

type mockPipeline struct{}

func (m *mockPipeline) Run(_ context.Context, _ *config.Config, log zerolog.Logger) error {
	log.Info().Msg("Mock pipeline executed")
	return nil
}

func Run() error {
	cfg, err := config.LoadConfig("configs/pipeline.yaml")
	if err != nil {
		return err
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		return err
	}
	logger.Info().Msg("Logger initialized")
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
