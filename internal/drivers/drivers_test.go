package drivers

import (
	"testing"
)

func TestLoadFromExcel(t *testing.T) {
	reg, err := LoadFromExcel("../../testdata/drivers_summary.xlsx")
	if err != nil {
		t.Fatalf("Failed to load drivers: %v", err)
	}

	if len(reg.drivers) == 0 {
		t.Fatal("No drivers loaded")
	}

	t.Logf("Loaded %d drivers", len(reg.drivers))

	// Проверяем, что водители проиндексированы по кодам городов
	if len(reg.byCityCode) == 0 {
		t.Fatal("No city codes indexed")
	}

	t.Logf("Indexed %d city codes", len(reg.byCityCode))
}

func TestGetRandomDriverForCity(t *testing.T) {
	reg, err := LoadFromExcel("../../testdata/drivers_summary.xlsx")
	if err != nil {
		t.Fatalf("Failed to load drivers: %v", err)
	}

	// Тестируем получение водителя для известного кода города
	driver := reg.GetRandomDriverForCity("I400")
	if driver == nil {
		t.Fatal("Expected driver for city code I400, got nil")
	}

	t.Logf("Driver for I400: %s, %s, %s", driver.Name, driver.LicenseNumber, driver.Phone)

	// Тестируем несуществующий код города - должен вернуть случайного водителя
	driver = reg.GetRandomDriverForCity("XXXXX")
	if driver == nil {
		t.Fatal("Expected random driver for non-existent city code, got nil")
	}
	t.Logf("Fallback driver for XXXXX: %s", driver.Name)
}
