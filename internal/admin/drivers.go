package admin

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Driver struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	CarNumber string `json:"car_number"`
	CityCodes string `json:"city_codes"`
	CityNames string `json:"city_names"`
}

type DriverService struct {
	dbPath string
}

func NewDriverService(dbPath string) *DriverService {
	return &DriverService{dbPath: dbPath}
}

func (s *DriverService) initDB() error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create table if not exists
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
		return err
	}

	// Migrate existing table - add missing columns
	db.Exec(`ALTER TABLE drivers ADD COLUMN car_number TEXT`)
	db.Exec(`ALTER TABLE drivers ADD COLUMN city_codes TEXT`)
	db.Exec(`ALTER TABLE drivers ADD COLUMN city_names TEXT`)
	// Migrate old 'cities' column to 'city_codes' if exists
	db.Exec(`UPDATE drivers SET city_codes = cities WHERE cities IS NOT NULL AND city_codes IS NULL`)
	
	return nil
}

func (s *DriverService) ListDrivers() ([]Driver, error) {
	if err := s.initDB(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, COALESCE(phone, ''), COALESCE(car_number, ''), COALESCE(city_codes, ''), COALESCE(city_names, '') FROM drivers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drivers := make([]Driver, 0)
	for rows.Next() {
		var d Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.Phone, &d.CarNumber, &d.CityCodes, &d.CityNames); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (s *DriverService) AddDriver(name, phone, carNumber, cityCodes, cityNames string) error {
	if err := s.initDB(); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO drivers (name, phone, car_number, city_codes, city_names) VALUES (?, ?, ?, ?, ?)`, name, phone, carNumber, cityCodes, cityNames)
	return err
}

func (s *DriverService) UpdateDriver(id int, name, phone, carNumber, cityCodes, cityNames string) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE drivers SET name = ?, phone = ?, car_number = ?, city_codes = ?, city_names = ? WHERE id = ?`, name, phone, carNumber, cityCodes, cityNames, id)
	return err
}

func (s *DriverService) DeleteDriver(id int) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM drivers WHERE id = ?`, id)
	return err
}
