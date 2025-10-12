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

	"github.com/maxviazov/excel-flow/internal/app"
)

type ProcessRequest struct {
	InputFile  string `json:"inputFile"`
	OutputFile string `json:"outputFile"`
}

type ProcessResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	InputRows    int    `json:"inputRows"`
	OutputRows   int    `json:"outputRows"`
	OutputFile   string `json:"outputFile"`
	ProcessTime  string `json:"processTime"`
}

func main() {
	http.HandleFunc("/api/upload", handleUpload)
	http.HandleFunc("/api/process", handleProcess)
	http.HandleFunc("/api/download/", handleDownload)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filePath := filepath.Join("uploads", filename)

	dst, err := os.Create(filePath)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"filename": filename, "path": filePath, "fullPath": filepath.Join("./uploads", filename)})
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

	start := time.Now()

	os.MkdirAll("./outputs", 0755)
	filename := fmt.Sprintf("דולינה_גרופ_%s.xlsx", time.Now().Format("2006-01-02"))
	outputPath := filepath.Join("./outputs", filename)

	inputRows, outputRows, err := app.ProcessFile(req.InputFile, outputPath)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Processing failed: %v", err),
		})
		return
	}

	respondJSON(w, http.StatusOK, ProcessResponse{
		Success:     true,
		Message:     "File processed successfully",
		InputRows:   inputRows,
		OutputRows:  outputRows,
		OutputFile:  filename,
		ProcessTime: time.Since(start).String(),
	})
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
