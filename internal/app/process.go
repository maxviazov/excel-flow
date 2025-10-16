package app

import (
	"os"

	"github.com/maxviazov/excel-flow/internal/drivers"
	"github.com/maxviazov/excel-flow/internal/ingest"
	"github.com/maxviazov/excel-flow/internal/pipelines"
	"github.com/maxviazov/excel-flow/internal/writer"
)

func ProcessFile(inputPath, outputPath string) (int, int, error) {
	data, err := ingest.ReadExcel(inputPath, "Sheet1", 1)
	if err != nil {
		return 0, 0, err
	}

	groups, err := pipelines.ProcessSAPData(data)
	if err != nil {
		return 0, 0, err
	}

	var driverRegistry *drivers.Registry
	// Try loading from DB first
	driversDBPath := "/tmp/data/drivers.db"
	if _, err := os.Stat(driversDBPath); err == nil {
		driverRegistry, _ = drivers.LoadFromDB(driversDBPath)
	} else {
		// Fallback to Excel
		driversPath := "testdata/drivers_summary.xlsx"
		if _, err := os.Stat(driversPath); err == nil {
			driverRegistry, _ = drivers.LoadFromExcel(driversPath)
		}
	}

	if err := writer.WriteMOH(outputPath, groups, driverRegistry); err != nil {
		return 0, 0, err
	}

	return len(data), len(groups), nil
}
