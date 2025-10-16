package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/maxviazov/excel-flow/internal/admin"
	"github.com/maxviazov/excel-flow/internal/app"
)

type ProcessRequest struct {
	InputFile  string `json:"inputFile"`
	OutputFile string `json:"outputFile"`
}

type ProcessResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	InputRows   int    `json:"inputRows"`
	OutputRows  int    `json:"outputRows"`
	OutputFile  string `json:"outputFile"`
	ProcessTime string `json:"processTime"`
}

var (
	cityService   *admin.CityService
	driverService *admin.DriverService
)

func main() {
	// Initialize admin services
	cityService = admin.NewCityService("configs/dictionaries/city.db")
	driverService = admin.NewDriverService("/tmp/data/drivers.db")
	
	// Ensure writable directory
	os.MkdirAll("/tmp/data", 0755)
	if _, err := os.Stat("/tmp/data/drivers.db"); os.IsNotExist(err) {
		if src, err := os.Open("configs/dictionaries/drivers.db"); err == nil {
			defer src.Close()
			if dst, err := os.Create("/tmp/data/drivers.db"); err == nil {
				defer dst.Close()
				io.Copy(dst, src)
			}
		}
	}
	
	// API endpoints
	http.HandleFunc("/api/upload", handleUpload)
	http.HandleFunc("/api/validate", handleValidate)
	http.HandleFunc("/api/process", handleProcess)
	http.HandleFunc("/api/export-csv/", handleExportCSV)
	http.HandleFunc("/api/download/", handleDownload)
	
	// Admin endpoints
	http.HandleFunc("/api/admin/cities", handleCities)
	http.HandleFunc("/api/admin/drivers", handleDrivers)
	
	http.HandleFunc("/health", handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"https://excel.viazov.dev":                                   true,
			"https://api.viazov.dev":                                     true,
			"http://excel-viazov-dev.s3-website-us-east-1.amazonaws.com": true,
			"https://d18sq2gf3s7zhe.cloudfront.net":                      true,
			"http://localhost:8080":                                      true,
			"http://localhost:3000":                                      true,
		}
		
		if origin != "" && allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Upload request from: %s, Method: %s", r.RemoteAddr, r.Method)
	
	if r.Method != "POST" {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filePath := filepath.Join("uploads", filename)

	log.Printf("Uploading file: %s -> %s", header.Filename, filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Failed to copy file: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}

	log.Printf("File uploaded successfully: %s", filePath)
	respondJSON(
		w, http.StatusOK,
		map[string]string{"filename": filename, "path": filePath, "fullPath": filepath.Join("./uploads", filename)},
	)
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	log.Printf("Processing file: %s", req.InputFile)

	start := time.Now()

	os.MkdirAll("./outputs", 0755)
	filename := fmt.Sprintf("דולינה_גרופ_%s.xlsx", time.Now().Format("2006-01-02"))
	outputPath := filepath.Join("./outputs", filename)

	inputRows, outputRows, err := app.ProcessFile(req.InputFile, outputPath)
	if err != nil {
		respondJSON(
			w, http.StatusInternalServerError, ProcessResponse{
				Success: false,
				Message: fmt.Sprintf("Processing failed: %v", err),
			},
		)
		return
	}

	respondJSON(
		w, http.StatusOK, ProcessResponse{
			Success:     true,
			Message:     "File processed successfully",
			InputRows:   inputRows,
			OutputRows:  outputRows,
			OutputFile:  filename,
			ProcessTime: time.Since(start).String(),
		},
	)
}

func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/api/export-csv/"):]
	excelPath := filepath.Join("./outputs", filename)

	if _, err := os.Stat(excelPath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Convert Excel to CSV (simplified - just return success for now)
	csvFilename := filename[:len(filename)-5] + ".csv"
	w.Header().Set("Content-Disposition", "attachment; filename="+csvFilename)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	
	// For now, return a simple message
	w.Write([]byte("CSV export feature - coming soon\n"))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/api/download/"):]
	filePath := filepath.Join("./outputs", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, filePath)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath string `json:"filePath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{"valid": false, "error": "Invalid request"})
		return
	}

	// Basic validation - check if file exists and is readable
	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"valid": false,
			"error": "File not found",
		})
		return
	}

	// Get file info
	fileInfo, _ := os.Stat(req.FilePath)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"valid":    true,
		"fileSize": fileInfo.Size(),
		"fileName": fileInfo.Name(),
		"warnings": []string{},
	})
}

// Admin handlers
func handleCities(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cities, err := cityService.ListCities()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, cities)
		
	case "POST":
		var req struct {
			Code    string `json:"code"`
			NameHeb string `json:"name_heb"`
			NameEng string `json:"name_eng"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			return
		}
		if err := cityService.AddCity(req.Code, req.NameHeb, req.NameEng); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"message": "City added"})
		
	case "DELETE":
		code := r.URL.Query().Get("code")
		if code == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing code"})
			return
		}
		if err := cityService.DeleteCity(code); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"message": "City deleted"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleDrivers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		drivers, err := driverService.ListDrivers()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, drivers)
		
	case "POST":
		var req struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Phone     string `json:"phone"`
			CarNumber string `json:"car_number"`
			Cities    string `json:"cities"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			return
		}
		
		if req.ID > 0 {
			// Update
			if err := driverService.UpdateDriver(req.ID, req.Name, req.Phone, req.CarNumber, req.Cities); err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			respondJSON(w, http.StatusOK, map[string]string{"message": "Driver updated"})
		} else {
			// Add
			if err := driverService.AddDriver(req.Name, req.Phone, req.CarNumber, req.Cities); err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			respondJSON(w, http.StatusOK, map[string]string{"message": "Driver added"})
		}
		
	case "DELETE":
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
			return
		}
		if err := driverService.DeleteDriver(id); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"message": "Driver deleted"})
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
