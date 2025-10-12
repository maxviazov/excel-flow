package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/maxviazov/excel-flow/internal/admin"
)

var (
	cityService   *admin.CityService
	driverService *admin.DriverService
)

func main() {
	cityService = admin.NewCityService("configs/dictionaries/city.db")
	driverService = admin.NewDriverService("configs/dictionaries/drivers.db")

	http.HandleFunc("/api/cities", handleCities)
	http.HandleFunc("/api/cities/alias", handleCityAlias)
	http.HandleFunc("/api/drivers", handleDrivers)
	http.Handle("/", http.FileServer(http.Dir("web/admin")))

	log.Println("Admin panel running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleCities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		cities, err := cityService.ListCities()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cities)

	case "POST":
		var req struct {
			Code    string `json:"code"`
			NameHeb string `json:"name_heb"`
			NameEng string `json:"name_eng"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := cityService.AddCity(req.Code, req.NameHeb, req.NameEng); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case "DELETE":
		code := r.URL.Query().Get("code")
		if err := cityService.DeleteCity(code); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleCityAlias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		var req struct {
			AliasHeb string `json:"alias_heb"`
			CityCode string `json:"city_code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := cityService.AddAlias(req.AliasHeb, req.CityCode); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case "DELETE":
		alias := r.URL.Query().Get("alias")
		if err := cityService.DeleteAlias(alias); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleDrivers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		drivers, err := driverService.ListDrivers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(drivers)

	case "POST":
		var req struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := driverService.AddDriver(req.Name, req.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case "PUT":
		var req struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := driverService.UpdateDriver(req.ID, req.Name, req.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := driverService.DeleteDriver(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
