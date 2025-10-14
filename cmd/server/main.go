package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	// Ensure writable directory for drivers DB
	os.MkdirAll("/tmp/data", 0755)

	// Copy drivers.db to writable location if not exists
	driverDBPath := "/tmp/data/drivers.db"
	if _, err := os.Stat(driverDBPath); os.IsNotExist(err) {
		if src, err := os.Open("configs/dictionaries/drivers.db"); err == nil {
			defer src.Close()
			if dst, err := os.Create(driverDBPath); err == nil {
				defer dst.Close()
				io.Copy(dst, src)
			}
		}
	}

	cityService = admin.NewCityService("configs/dictionaries/city.db")
	driverService = admin.NewDriverService(driverDBPath)

	// API endpoints only
	http.HandleFunc("/api/upload", handleUpload)
	http.HandleFunc("/api/process", handleProcess)
	http.HandleFunc("/api/download/", handleDownload)
	http.HandleFunc("/api/admin/cities", handleCities)
	http.HandleFunc("/api/admin/cities/alias", handleCityAlias)
	http.HandleFunc("/api/admin/cities/import", handleCitiesImport)
	http.HandleFunc("/api/admin/drivers", handleDrivers)
	http.HandleFunc("/api/admin/drivers/import", handleDriversImport)
	http.HandleFunc("/api/admin/drivers/template", handleDriversTemplate)
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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
