package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("testdata/drivers_summary.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "configs/dictionaries/drivers.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS drivers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			phone TEXT,
			car_number TEXT,
			city_codes TEXT,
			city_names TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("DELETE FROM drivers")

	count := 0
	for _, row := range rows[1:] {
		if len(row) < 4 || row[0] == "" {
			continue
		}

		name := strings.TrimSpace(row[0])
		carNumber := strings.TrimSpace(row[1])
		phone := strings.TrimSpace(row[2])
		citiesStr := strings.TrimSpace(row[3])

		citiesStr = strings.Trim(citiesStr, "[]")
		citiesStr = strings.ReplaceAll(citiesStr, "'", "")
		citiesStr = strings.ReplaceAll(citiesStr, " ", "")

		if citiesStr == "" {
			continue
		}

		_, err = db.Exec(`INSERT INTO drivers (name, phone, car_number, city_codes) VALUES (?, ?, ?, ?)`,
			name, phone, carNumber, citiesStr)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		count++
	}

	fmt.Printf("✅ Imported %d drivers\n", count)
}
