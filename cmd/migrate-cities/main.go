package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Открываем Excel файл
	f, err := excelize.OpenFile("testdata/רשימת יישובים עדכנית 16.4.24.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Получаем список листов
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		log.Fatal("No sheets found")
	}

	// Используем первый лист
	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	if len(rows) < 2 {
		log.Fatal("Not enough rows")
	}

	// Создаем SQL файл
	sqlFile, err := os.Create("migrations/012_add_missing_cities.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlFile.Close()

	// Заголовок SQL файла
	if _, err := sqlFile.WriteString("-- 012_add_missing_cities.sql\n"); err != nil {
		log.Fatal("Error writing header:", err)
	}
	if _, err := sqlFile.WriteString("-- Adding missing cities from official MOH list\n\n"); err != nil {
		log.Fatal("Error writing description:", err)
	}
	if _, err := sqlFile.WriteString("BEGIN;\n\n"); err != nil {
		log.Fatal("Error writing BEGIN:", err)
	}

	// Обрабатываем строки (пропускаем заголовок)
	for i, row := range rows {
		if i == 0 {
			// Выводим заголовки для понимания структуры
			fmt.Printf("Headers: %v\n", row)
			continue
		}

		if len(row) < 2 {
			continue
		}

		// Структура: название, код
		cityName := strings.TrimSpace(row[0])
		cityCode := strings.TrimSpace(row[1])

		if cityCode == "" || cityName == "" {
			continue
		}

		// Генерируем SQL INSERT с ON CONFLICT IGNORE
		sql := fmt.Sprintf(
			`INSERT OR IGNORE INTO city_codes (city_heb, city_code) VALUES ('%s', '%s');`+"\n",
			strings.ReplaceAll(cityName, "'", "''"), // экранируем апострофы
			cityCode,
		)

		_, err := sqlFile.WriteString(sql)
		if err != nil {
			log.Printf("Error writing SQL for city %s: %v", cityName, err)
			continue
		}
	}

	if _, err := sqlFile.WriteString("\nCOMMIT;\n"); err != nil {
		log.Fatal("Error writing COMMIT:", err)
	}
	fmt.Printf("Migration created: migrations/012_add_missing_cities.sql\n")
	fmt.Printf("Processed %d rows\n", len(rows)-1)
}
