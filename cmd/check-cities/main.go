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
	// Открываем Excel файл
	f, err := excelize.OpenFile("testdata/רשימת יישובים עדכנית 16.4.24.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Читаем данные из Excel
	sheets := f.GetSheetList()
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		log.Fatal(err)
	}

	// Открываем базу данных
	db, err := sql.Open("sqlite3", "configs/dictionaries/city.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Проверка соответствия Excel файла и базы данных...")

	missingInDB := []string{}
	excelCities := make(map[string]string) // code -> name

	// Собираем данные из Excel (пропускаем заголовок)
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		cityName := strings.TrimSpace(row[0])
		cityCode := strings.TrimSpace(row[1])

		if cityName != "" && cityCode != "" {
			excelCities[cityCode] = cityName
		}
	}

	// Проверяем каждый город из Excel в базе
	for code, name := range excelCities {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM city_codes WHERE city_code = ?`, code).Scan(&count)
		if err != nil {
			log.Printf("Error checking code %s: %v", code, err)
			continue
		}

		if count == 0 {
			missingInDB = append(missingInDB, fmt.Sprintf("'%s', '%s'", name, code))
		}
	}

	// Выводим результаты
	fmt.Printf("Городов в Excel: %d\n", len(excelCities))

	if len(missingInDB) > 0 {
		fmt.Printf("Города из Excel, отсутствующие в базе (%d):\n", len(missingInDB))
		fmt.Println("INSERT OR IGNORE INTO city_codes (city_heb, city_code) VALUES")
		for i, city := range missingInDB {
			if i == len(missingInDB)-1 {
				fmt.Printf("(%s);\n", city)
			} else {
				fmt.Printf("(%s),\n", city)
			}
		}
	} else {
		fmt.Println("✅ Все города из Excel есть в базе данных")
	}
}
