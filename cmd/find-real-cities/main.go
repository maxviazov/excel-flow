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
	f, err := excelize.OpenFile("testdata/משקל.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	db, err := sql.Open("sqlite3", "configs/dictionaries/city.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sheets := f.GetSheetList()
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Анализ городов из файла משקל.xlsx...")

	cityMap := make(map[string]int)

	for i, row := range rows {
		if i == 0 || len(row) < 8 {
			continue
		}

		// Колонка с адресом (כתובת)
		address := ""
		if len(row) > 7 {
			address = strings.TrimSpace(row[7])
		}

		if address == "" {
			continue
		}

		// Извлекаем город (до первой запятой)
		parts := strings.SplitN(address, ",", 2)
		city := strings.TrimSpace(parts[0])

		if len(city) > 2 && len(city) < 25 {
			cityMap[city]++
		}
	}

	fmt.Printf("Найдено городов: %d\n\n", len(cityMap))

	aliases := []string{}
	missing := []string{}

	for city, count := range cityMap {
		if count < 2 {
			continue
		}

		var dbCount int
		err := db.QueryRow(`SELECT COUNT(*) FROM v_city_lookup WHERE key_heb = ?`, city).Scan(&dbCount)
		if err != nil {
			log.Printf("Error checking city %s: %v", city, err)
			continue
		}

		if dbCount == 0 {
			similar := findSimilar(db, city)
			if similar != "" {
				aliases = append(aliases, fmt.Sprintf("('%s', '%s')", city, similar))
			} else {
				missing = append(missing, fmt.Sprintf("-- %s (встречается %d раз)", city, count))
			}
		}
	}

	if len(aliases) > 0 {
		fmt.Println("INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES")
		for i, alias := range aliases {
			if i == len(aliases)-1 {
				fmt.Printf("%s;\n", alias)
			} else {
				fmt.Printf("%s,\n", alias)
			}
		}
	}

	if len(missing) > 0 {
		fmt.Println("\nНе найдены похожие:")
		for _, m := range missing {
			fmt.Println(m)
		}
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
		// Логируем только неожиданные ошибки
		if err.Error() != "sql: no rows in result set" {
			log.Printf("Error searching similar city for pattern %s: %v", pattern, err)
		}
	}

	return ""
}
