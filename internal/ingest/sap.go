package ingest

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ReadExcel(filePath, sheet string, headerRow int) ([]map[string]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the excel file: %w", err)
	}
	defer f.Close()

	if sheet == "" {
		sheet = f.GetSheetName(0)
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to read an excel file: %w", err)
	}
	if len(rows) < headerRow {
		return nil, fmt.Errorf("no data rows: headerRow=%d", headerRow)
	}
	headers := rows[headerRow-1] // headerRow is 1-based

	// Создаем маппинг заголовков к индексам (берем первое вхождение)
	headerToIndex := make(map[string]int)
	for i, header := range headers {
		if header != "" {
			if _, exists := headerToIndex[header]; !exists {
				headerToIndex[header] = i // берем первый индекс
			}
		}
	}

	var out []map[string]string
	for _, r := range rows[headerRow:] {
		m := map[string]string{}
		for header, index := range headerToIndex {
			if index < len(r) {
				m[header] = r[index]
			}
		}
		if len(m) > 0 {
			out = append(out, m)
		}
	}
	return out, nil
}
