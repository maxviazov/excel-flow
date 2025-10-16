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
	Cities    string `json:"cities"`
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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS drivers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			phone TEXT,
			car_number TEXT,
			cities TEXT
		)
	`)
	return err
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

	rows, err := db.Query(`SELECT id, name, COALESCE(phone, ''), COALESCE(car_number, ''), COALESCE(cities, '') FROM drivers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var d Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.Phone, &d.CarNumber, &d.Cities); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (s *DriverService) AddDriver(name, phone, carNumber, cities string) error {
	if err := s.initDB(); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO drivers (name, phone, car_number, cities) VALUES (?, ?, ?, ?)`, name, phone, carNumber, cities)
	return err
}

func (s *DriverService) UpdateDriver(id int, name, phone, carNumber, cities string) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE drivers SET name = ?, phone = ?, car_number = ?, cities = ? WHERE id = ?`, name, phone, carNumber, cities, id)
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
