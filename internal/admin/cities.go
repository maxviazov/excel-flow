package admin

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type City struct {
	Code      string `json:"code"`
	NameHeb   string `json:"name_heb"`
	NameEng   string `json:"name_eng"`
	IsAlias   bool   `json:"is_alias"`
	CanonCode string `json:"canon_code,omitempty"`
}

type CityService struct {
	dbPath string
}

func NewCityService(dbPath string) *CityService {
	return &CityService{dbPath: dbPath}
}

func (s *CityService) ListCities() ([]City, error) {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT city_code, canon_heb, canon_eng, 0 as is_alias, '' as canon_code
		FROM city_codes
		UNION ALL
		SELECT alias_heb, alias_heb, '', 1, city_code
		FROM city_aliases
		ORDER BY canon_heb
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var c City
		if err := rows.Scan(&c.Code, &c.NameHeb, &c.NameEng, &c.IsAlias, &c.CanonCode); err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}
	return cities, nil
}

func (s *CityService) AddCity(code, nameHeb, nameEng string) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO city_codes (city_code, canon_heb, canon_eng) VALUES (?, ?, ?)`,
		code, nameHeb, nameEng)
	return err
}

func (s *CityService) AddAlias(aliasHeb, cityCode string) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists int
	err = db.QueryRow(`SELECT COUNT(*) FROM city_codes WHERE city_code = ?`, cityCode).Scan(&exists)
	if err != nil {
		return err
	}
	if exists == 0 {
		return fmt.Errorf("city code %s not found", cityCode)
	}

	_, err = db.Exec(`INSERT INTO city_aliases (alias_heb, city_code) VALUES (?, ?)`, aliasHeb, cityCode)
	return err
}

func (s *CityService) DeleteCity(code string) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM city_codes WHERE city_code = ?`, code)
	return err
}

func (s *CityService) DeleteAlias(aliasHeb string) error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM city_aliases WHERE alias_heb = ?`, aliasHeb)
	return err
}
