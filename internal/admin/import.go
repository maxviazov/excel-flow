package admin

import (
	"database/sql"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ImportResult struct {
	Added   int `json:"added"`
	Updated int `json:"updated"`
	Skipped int `json:"skipped"`
}

func (s *CityService) ImportFromExcel(filePath string) (*ImportResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result := &ImportResult{}

	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		code := strings.TrimSpace(row[0])
		nameHeb := strings.TrimSpace(row[1])
		nameEng := ""
		if len(row) > 2 {
			nameEng = strings.TrimSpace(row[2])
		}

		if code == "" || nameHeb == "" {
			result.Skipped++
			continue
		}

		var exists int
		err := db.QueryRow(`SELECT COUNT(*) FROM city_codes WHERE city_code = ?`, code).Scan(&exists)
		if err != nil {
			return nil, err
		}

		if exists > 0 {
			_, err = db.Exec(`UPDATE city_codes SET city_heb = ? WHERE city_code = ?`,
				nameHeb, code)
			if err != nil {
				return nil, err
			}
			result.Updated++
		} else {
			_, err = db.Exec(`INSERT INTO city_codes (city_code, city_heb) VALUES (?, ?)`,
				code, nameHeb)
			if err != nil {
				return nil, err
			}
			result.Added++
		}
	}

	return result, nil
}

func (s *DriverService) ImportFromExcel(filePath string) (*ImportResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result := &ImportResult{}

	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}

		name := strings.TrimSpace(row[0])
		carNumber := strings.TrimSpace(row[1])
		phone := strings.TrimSpace(row[2])
		cities := ""
		if len(row) > 3 {
			cities = strings.TrimSpace(row[3])
		}

		if name == "" {
			result.Skipped++
			continue
		}

		var driverID int
		err := db.QueryRow(`SELECT id FROM drivers WHERE name = ?`, name).Scan(&driverID)

		if err == sql.ErrNoRows {
			_, err := db.Exec(`INSERT INTO drivers (name, phone, car_number, cities) VALUES (?, ?, ?, ?)`,
				name, phone, carNumber, cities)
			if err != nil {
				return nil, err
			}
			result.Added++
		} else if err != nil {
			return nil, err
		} else {
			_, err = db.Exec(`UPDATE drivers SET phone = ?, car_number = ?, cities = ? WHERE id = ?`,
				phone, carNumber, cities, driverID)
			if err != nil {
				return nil, err
			}
			result.Updated++
		}
	}

	return result, nil
}

func (s *DriverService) ExportTemplate(filePath string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := f.GetSheetName(0)
	f.SetCellValue(sheet, "A1", "שם נהג")
	f.SetCellValue(sheet, "B1", "מספר רכב")
	f.SetCellValue(sheet, "C1", "טלפון")
	f.SetCellValue(sheet, "D1", "ערים (מופרד בפסיק)")

	f.SetCellValue(sheet, "A2", "דוד כהן")
	f.SetCellValue(sheet, "B2", "12-345-67")
	f.SetCellValue(sheet, "C2", "050-1234567")
	f.SetCellValue(sheet, "D2", "ירושלים, תל אביב, חיפה")

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "D1", style)

	return f.SaveAs(filePath)
}
