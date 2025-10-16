package pipelines

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/maxviazov/excel-flow/internal/mapping"
)

// ProcessSAPData processes raw Excel data and returns grouped results
func ProcessSAPData(data []map[string]string) (map[GroupKey]*GroupVal, error) {
	groups := make(map[GroupKey]*GroupVal)

	// Apply SAP to internal mapping
	mappingRules := mapping.MapSAPtoInternal()
	mappedData := mapping.ApplySelect(data, mappingRules)

	var skippedReasons []string
	for i, row := range mappedData {
		// Extract and validate required fields
		clientLicense := strings.TrimSpace(row["client_license_number"])
		if clientLicense == "" {
			skippedReasons = append(skippedReasons, fmt.Sprintf("row %d: empty client_license", i+2))
			continue
		}

		orderID := strings.TrimSpace(row["order_id"])
		if orderID == "" {
			skippedReasons = append(skippedReasons, fmt.Sprintf("row %d: empty order_id", i+2))
			continue
		}

		// Дата не используется для группировки

		// Parse weight - convert from tons to kg
		weightStr := normalizeNumber(row["total_weight_raw"])
		weightTons, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			weightTons = 0
		}
		weight := weightTons / 1000 // конверсия грамм в кг

		// Parse and normalize packages
		packagesStr := normalizeNumber(row["total_packaging_raw"])
		packages, err := strconv.ParseFloat(packagesStr, 64)
		if err != nil {
			packages = 0
		}

		// Get client info - use ONLY Hebrew name
		clientName := strings.TrimSpace(row["client_name_he"])
		fullAddress := strings.TrimSpace(row["client_address"])

		// Извлекаем город из адреса (до первой запятой)
		cityName, streetAddress := extractCityFromAddress(fullAddress)
		cityCode, cityNameHeb := lookupCityInfo(cityName)
		// В колонку address записываем только улицу, город указывается кодом
		address := streetAddress

		// Skip rows where weight <= 0
		if weight <= 0 {
			skippedReasons = append(skippedReasons, fmt.Sprintf("row %d: weight <= 0 (raw: %s, parsed: %.3f)", i+2, row["total_weight_raw"], weight))
			continue
		}

		// Create group key - по клиенту И номеру документа
		key := GroupKey{
			ClientLicense: clientLicense,
			OrderID:       orderID,     // группируем по номеру документа
			Date:          time.Time{}, // не группируем по дате
		}

		// Aggregate values
		if existing, ok := groups[key]; ok {
			existing.TotalWeight += weight
			existing.TotalPackages += packages
			existing.Count++
			// Update client info if empty
			if existing.ClientName == "" && clientName != "" {
				existing.ClientName = clientName
			}
			if existing.Address == "" && address != "" {
				existing.Address = address
			}
			if existing.CityName == "" && cityName != "" {
				existing.CityName = cityName
			}
			if existing.CityNameHeb == "" && cityNameHeb != "" {
				existing.CityNameHeb = cityNameHeb
			}
			if existing.CityCode == "" && cityCode != "" {
				existing.CityCode = cityCode
			}
			// Добавляем номер документа в список
			if orderID != "" {
				// Проверяем что номер еще не добавлен
				found := false
				for _, existingOrderID := range existing.OrderIDs {
					if existingOrderID == orderID {
						found = true
						break
					}
				}
				if !found {
					existing.OrderIDs = append(existing.OrderIDs, orderID)
				}
			}
		} else {
			orderIDs := []string{}
			if orderID != "" {
				orderIDs = append(orderIDs, orderID)
			}
			groups[key] = &GroupVal{
				TotalWeight:   weight,
				TotalPackages: packages,
				Count:         1,
				ClientName:    clientName,
				Address:       address,
				OrderIDs:      orderIDs,
				CityName:      cityName,
				CityNameHeb:   cityNameHeb,
				CityCode:      cityCode,
			}
		}
	}

	if len(groups) == 0 {
		errMsg := fmt.Sprintf("no valid data found in %d rows", len(data))
		if len(skippedReasons) > 0 {
			errMsg += ". Reasons: " + strings.Join(skippedReasons[:min(5, len(skippedReasons))], "; ")
		}
		return nil, fmt.Errorf("%s", errMsg)
	}

	return groups, nil
}

// parseDate tries multiple date formats
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"02/01/06",
		"02/01/2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// normalizeNumber replaces comma with dot and trims spaces
func normalizeNumber(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	return s
}

// extractCityFromAddress извлекает город до первой запятой и оставляет адрес
func extractCityFromAddress(fullAddress string) (city, address string) {
	parts := strings.SplitN(fullAddress, ",", 2)
	if len(parts) >= 2 {
		city = strings.TrimSpace(parts[0])
		address = strings.TrimSpace(parts[1])
	} else {
		// Если нет запятой, весь адрес считаем городом
		city = strings.TrimSpace(fullAddress)
		address = ""
	}
	return city, address
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// lookupCityInfo ищет код города и название на иврите в SQLite базе
func lookupCityInfo(cityName string) (code, hebName string) {
	cityName = strings.TrimSpace(cityName)
	if cityName == "" {
		return "9999", "" // код по умолчанию
	}

	// Открываем базу данных
	db, err := sql.Open("sqlite3", "configs/dictionaries/city.db")
	if err != nil {
		return "9999", cityName // ошибка подключения, возвращаем исходное название
	}
	defer db.Close()

	// Ищем код города и название на иврите через view v_city_lookup
	// Приоритет: сначала не-алиасы (is_alias=0), потом по алфавиту кода
	var cityCode, cityHeb string
	query := `SELECT city_code, canon_heb FROM v_city_lookup WHERE key_heb = ? ORDER BY is_alias ASC, city_code ASC LIMIT 1`
	err = db.QueryRow(query, cityName).Scan(&cityCode, &cityHeb)
	if err != nil {
		return "9999", cityName // город не найден, возвращаем исходное название
	}

	return cityCode, cityHeb
}
