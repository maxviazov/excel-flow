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
	// Открываем Excel файл с базой клиентов SAP
	f, err := excelize.OpenFile("testdata/База клиентов по адресам_ ТП_ районам.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Открываем базу данных
	db, err := sql.Open("sqlite3", "configs/dictionaries/city.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Читаем данные из Excel
	sheets := f.GetSheetList()
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Анализ городов из базы SAP...")

	// Находим колонку עיר
	cityColumnIndex := -1
	if len(rows) > 0 {
		for i, header := range rows[0] {
			if strings.Contains(header, "עיר") {
				cityColumnIndex = i
				break
			}
		}
	}

	if cityColumnIndex == -1 {
		fmt.Println("Колонка עיר не найдена")
		fmt.Println("Заголовки:", rows[0])
		return
	}

	fmt.Printf("Колонка עיר найдена в позиции %d\n", cityColumnIndex)

	// Собираем уникальные города
	cityMap := make(map[string]int)

	for i, row := range rows {
		if i == 0 || len(row) <= cityColumnIndex {
			continue
		}

		city := strings.TrimSpace(row[cityColumnIndex])
		if city != "" && len(city) > 1 {
			cityMap[city]++
		}
	}

	fmt.Printf("Найдено уникальных городов в SAP: %d\n\n", len(cityMap))

	// Проверяем каждый город в базе
	missingCities := []string{}
	aliases := []string{}

	for city, count := range cityMap {
		if count < 2 { // игнорируем редкие города
			continue
		}

		var dbCount int
		err := db.QueryRow(`SELECT COUNT(*) FROM v_city_lookup WHERE key_heb = ?`, city).Scan(&dbCount)
		if err != nil {
			log.Printf("Error checking city %s: %v", city, err)
			continue
		}

		if dbCount == 0 {
			// Ищем похожие города
			similar := findSimilar(db, city)
			if similar != "" {
				aliases = append(aliases, fmt.Sprintf("('%s', '%s')", city, similar))
			} else {
				missingCities = append(missingCities, fmt.Sprintf("-- %s (встречается %d раз)", city, count))
			}
		}
	}

	if len(aliases) > 0 {
		fmt.Println("Предлагаемые алиасы:")
		fmt.Println("INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES")
		for i, alias := range aliases {
			if i == len(aliases)-1 {
				fmt.Printf("%s;\n", alias)
			} else {
				fmt.Printf("%s,\n", alias)
			}
		}
		fmt.Println()
	}

	if len(missingCities) > 0 {
		fmt.Println("Города без похожих:")
		for _, city := range missingCities {
			fmt.Println(city)
		}
	} else if len(aliases) == 0 {
		fmt.Println("✅ Все города из SAP найдены в базе данных")
	}
}

func findSimilar(db *sql.DB, city string) string {
	patterns := []string{
		"%" + city + "%",
		city + "%",
		"%" + city,
	}

	for _, pattern := range patterns {
		var similar string
		err := db.QueryRow(`SELECT canon_heb FROM v_city_lookup WHERE canon_heb LIKE ? LIMIT 1`, pattern).Scan(&similar)
		if err == nil {
			return similar
		}
		// Логируем только неожиданные ошибки (не "no rows")
		if err.Error() != "sql: no rows in result set" {
			log.Printf("Error searching similar city for pattern %s: %v", pattern, err)
		}
	}

	return ""
}
