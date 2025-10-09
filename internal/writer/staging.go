package writer

import (
	"fmt"
	"sort"
	"time"

	"github.com/maxviazov/excel-flow/internal/pipelines"

	"github.com/xuri/excelize/v2"
)

func WriteStaging(path string, groups map[pipelines.GroupKey]*pipelines.GroupVal) error {
	f := excelize.NewFile()
	sh := "data"
	f.NewSheet(sh)
	f.DeleteSheet("Sheet1")

	// заголовок
	headers := []string{
		"client_license_number", "order_id", "date",
		"total_weight", "total_packages", "rows_count",
		"client_name", "address",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sh, cell, h)
	}

	// отсортируем ключи для стабильности
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

	// строки
	row := 2
	for _, k := range keys {
		v := groups[k]
		values := []any{
			k.ClientLicense,
			k.OrderID, // один номер документа на строку
			time.Now().Format("2006-01-02"),
			fmt.Sprintf("%.3f", v.TotalWeight),
			fmt.Sprintf("%.3f", v.TotalPackages),
			v.Count,
			v.ClientName,
			v.Address,
		}
		for i, val := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			f.SetCellValue(sh, cell, val)
		}
		row++
	}

	return f.SaveAs(path)
}
