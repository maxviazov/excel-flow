package drivers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Driver struct {
	Name          string
	LicenseNumber string
	Phone         string
	CityCodes     []string
}

type Registry struct {
	drivers       []Driver
	byCityCode    map[string][]Driver
}

func LoadFromExcel(path string) (*Registry, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open drivers file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, fmt.Errorf("failed to read drivers file: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("drivers file is empty")
	}

	reg := &Registry{
		drivers:    make([]Driver, 0),
		byCityCode: make(map[string][]Driver),
	}

	for i, row := range rows[1:] {
		if len(row) < 4 || row[0] == "" {
			continue
		}

		driver := Driver{
			Name:          strings.TrimSpace(row[0]),
			LicenseNumber: strings.TrimSpace(row[1]),
			Phone:         strings.TrimSpace(row[2]),
		}

		// Парсим список кодов городов из строки вида "['F1381', 'F2376', ...]"
		citiesStr := strings.TrimSpace(row[3])
		citiesStr = strings.Trim(citiesStr, "[]")
		citiesStr = strings.ReplaceAll(citiesStr, "'", "")
		citiesStr = strings.ReplaceAll(citiesStr, " ", "")
		
		if citiesStr != "" {
			codes := strings.Split(citiesStr, ",")
			for _, code := range codes {
				code = strings.TrimSpace(code)
				if code != "" {
					driver.CityCodes = append(driver.CityCodes, strings.ToUpper(code))
				}
			}
		}

		if len(driver.CityCodes) == 0 {
			continue
		}

		reg.drivers = append(reg.drivers, driver)

		// Индексируем по кодам городов
		for _, code := range driver.CityCodes {
			reg.byCityCode[code] = append(reg.byCityCode[code], driver)
		}

		_ = i
	}

	return reg, nil
}

// GetRandomDriverForCity возвращает случайного водителя для указанного кода города
// Если для города нет водителей, возвращает случайного из всех доступных
func (r *Registry) GetRandomDriverForCity(cityCode string) *Driver {
	cityCode = strings.ToUpper(cityCode)
	drivers := r.byCityCode[cityCode]
	if len(drivers) == 0 {
		// Fallback: возвращаем случайного водителя из всех
		if len(r.drivers) == 0 {
			return nil
		}
		return &r.drivers[rand.Intn(len(r.drivers))]
	}
	return &drivers[rand.Intn(len(drivers))]
}

// ParseCityCodes парсит строку с кодами городов из JSON-подобного формата
func parseCityCodes(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	// Пробуем распарсить как JSON
	var codes []string
	if err := json.Unmarshal([]byte(s), &codes); err == nil {
		return codes, nil
	}

	// Если не JSON, парсим вручную
	s = strings.Trim(s, "[]")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, " ", "")
	
	if s == "" {
		return nil, nil
	}

	return strings.Split(s, ","), nil
}
