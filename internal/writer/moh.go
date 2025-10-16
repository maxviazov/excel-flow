package writer

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/maxviazov/excel-flow/internal/drivers"
	"github.com/maxviazov/excel-flow/internal/pipelines"
	"github.com/maxviazov/excel-flow/internal/textutil"
)

func WriteMOH(path string, groups map[pipelines.GroupKey]*pipelines.GroupVal, driverRegistry *drivers.Registry) error {
	// Open template file
	f, err := excelize.OpenFile("testdata/sample.xlsx")
	if err != nil {
		return fmt.Errorf("failed to open template: %w", err)
	}
	sh := "Sheet1"

	// Sort keys for stability
	keys := make([]pipelines.GroupKey, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(
		keys, func(i, j int) bool {
			if keys[i].ClientLicense != keys[j].ClientLicense {
				return keys[i].ClientLicense < keys[j].ClientLicense
			}
			if keys[i].OrderID != keys[j].OrderID {
				return keys[i].OrderID < keys[j].OrderID
			}
			return keys[i].Date.Before(keys[j].Date)
		},
	)

	// Create date style (date only, no time)
	dateStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 14, // dd/mm/yyyy format
	})

	// Write data rows
	row := 2
	for _, k := range keys {
		v := groups[k]
		
		// Copy styles from template row 2
		for col := 1; col <= 30; col++ {
			templateCell, _ := excelize.CoordinatesToCellName(col, 2)
			targetCell, _ := excelize.CoordinatesToCellName(col, row)
			styleID, _ := f.GetCellStyle(sh, templateCell)
			f.SetCellStyle(sh, targetCell, targetCell, styleID)
		}
		
		// Set values
		f.SetCellStr(sh, fmt.Sprintf("A%d", row), "דולינה גרופ בע\"מ")
		f.SetCellValue(sh, fmt.Sprintf("B%d", row), 511777856)
		f.SetCellStr(sh, fmt.Sprintf("C%d", row), "P1908")
		// Set date with date-only format
		f.SetCellValue(sh, fmt.Sprintf("D%d", row), time.Now())
		f.SetCellStyle(sh, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), dateStyle)
		
		// 5-7: vehicle and driver info
		if driverRegistry != nil {
			driver := driverRegistry.GetRandomDriverForCity(v.CityCode)
			if driver != nil {
				f.SetCellStr(sh, fmt.Sprintf("E%d", row), driver.LicenseNumber)
				f.SetCellStr(sh, fmt.Sprintf("F%d", row), driver.Name)
				f.SetCellStr(sh, fmt.Sprintf("G%d", row), driver.Phone)
			}
		}
		
		if v.ClientName != "" {
			f.SetCellStr(sh, fmt.Sprintf("H%d", row), textutil.SanitizeForMOH(v.ClientName))
		}
		f.SetCellStr(sh, fmt.Sprintf("I%d", row), "קמעונאי")
		if v.CityCode != "" {
			f.SetCellStr(sh, fmt.Sprintf("J%d", row), v.CityCode)
		}
		if v.Address != "" {
			f.SetCellStr(sh, fmt.Sprintf("K%d", row), textutil.SanitizeForMOH(v.Address))
		}
		f.SetCellValue(sh, fmt.Sprintf("L%d", row), parseNumber(k.ClientLicense))
		f.SetCellValue(sh, fmt.Sprintf("M%d", row), 0)
		f.SetCellValue(sh, fmt.Sprintf("N%d", row), parseNumber(k.OrderID))
		// O-V are empty (meat/fish categories)
		f.SetCellValue(sh, fmt.Sprintf("W%d", row), v.TotalWeight) // Col 23
		// X, Y are empty
		f.SetCellValue(sh, fmt.Sprintf("Z%d", row), v.TotalWeight) // Col 26
		f.SetCellValue(sh, fmt.Sprintf("AA%d", row), v.TotalWeight) // Col 27
		f.SetCellValue(sh, fmt.Sprintf("AB%d", row), 1) // Col 28
		// AC, AD are empty
		row++
	}

	return f.SaveAs(path)
}

func boolPtr(b bool) *bool {
	return &b
}

func parseNumber(s string) interface{} {
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num
	}
	return s
}
