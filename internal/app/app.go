package app

import (
	"context"
	"os"

	"github.com/rs/zerolog"

	"github.com/maxviazov/excel-flow/internal/config"
	"github.com/maxviazov/excel-flow/internal/ingest"
	"github.com/maxviazov/excel-flow/internal/log"
	"github.com/maxviazov/excel-flow/internal/pipelines"
	"github.com/maxviazov/excel-flow/internal/writer"
)

type Pipeline interface {
	Run(ctx context.Context, cfg *config.Config, log zerolog.Logger) error
}

type realPipeline struct{}

func (r *realPipeline) Run(_ context.Context, cfg *config.Config, log zerolog.Logger) error {
	// Check if input file exists
	if _, err := os.Stat(cfg.Source.File); os.IsNotExist(err) {
		log.Error().Str("file", cfg.Source.File).Msg("Input file does not exist")
		return err
	}

	// Read Excel file
	log.Info().Str("file", cfg.Source.File).Str("sheet", cfg.Source.Sheet).Msg("Reading Excel file")
	data, err := ingest.ReadExcel(cfg.Source.File, cfg.Source.Sheet, cfg.Source.HeaderRow)
	if err != nil {
		return err
	}
	log.Info().Int("rows", len(data)).Msg("Excel file read successfully")

	// Process data through pipeline
	log.Info().Msg("Processing data through pipeline")
	groups, err := pipelines.ProcessSAPData(data)
	if err != nil {
		return err
	}
	log.Info().Int("groups", len(groups)).Msg("Data processed successfully")

	// Create output directory if it doesn't exist
	if err := os.MkdirAll("out", 0755); err != nil {
		return err
	}

	// Write staging file
	log.Info().Str("file", cfg.Output.StagingFile).Msg("Writing staging file")
	if err := writer.WriteStaging(cfg.Output.StagingFile, groups); err != nil {
		return err
	}
	log.Info().Msg("Staging file written successfully")

	// Write MOH file
	mohFile := "out/moh_staging.xlsx"
	log.Info().Str("file", mohFile).Msg("Writing MOH file")
	if err := writer.WriteMOH(mohFile, groups); err != nil {
		return err
	}
	log.Info().Msg("MOH file written successfully")

	return nil
}

func Run() error {
	// Очистка директории out
	if err := os.RemoveAll("out"); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll("out", 0755); err != nil {
		return err
	}

	cfg, err := config.LoadConfig("configs/pipeline.yaml")
	if err != nil {
		return err
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		return err
	}
	logger.Info().Msg("Logger initialized")
	pipeline := &realPipeline{}

	logger.Info().Msg("🚀 Starting Excel-Flow")

	ctx := context.Background()
	if err := pipeline.Run(ctx, cfg, logger); err != nil {
		logger.Error().Err(err).Msg("Pipeline execution failed")
		return err
	}

	logger.Info().Msg("✅ Excel-Flow finished successfully")
	return nil
}
