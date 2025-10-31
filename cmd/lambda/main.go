package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

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

func init() {
	cityService = admin.NewCityService("configs/dictionaries/city.db")
	driverService = admin.NewDriverService("/tmp/data/drivers.db")

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
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "https://excel.viazov.dev",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
		"Content-Type":                 "application/json",
	}

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	switch request.Path {
	case "/process":
		return handleProcessLambda(request, headers)
	case "/admin/cities":
		return handleCitiesLambda(request, headers)
	case "/admin/drivers":
		return handleDriversLambda(request, headers)
	case "/health":
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       `{"status":"ok"}`,
		}, nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"path":"%s","resource":"%s","method":"%s"}`, request.Path, request.Resource, request.HTTPMethod),
		}, nil
	}
}

func handleProcessLambda(request events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var req ProcessRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       `{"error":"Invalid request"}`,
		}, nil
	}

	start := time.Now()
	filename := fmt.Sprintf("דולינה_גרופ_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))
	outputPath := filepath.Join("/tmp", filename)

	inputRows, outputRows, err := app.ProcessFile(req.InputFile, outputPath)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"success":false,"message":"Processing failed: %v"}`, err),
		}, nil
	}

	response := ProcessResponse{
		Success:     true,
		Message:     "File processed successfully",
		InputRows:   inputRows,
		OutputRows:  outputRows,
		OutputFile:  filename,
		ProcessTime: time.Since(start).String(),
	}

	body, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func handleCitiesLambda(request events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	cities, err := cityService.ListCities()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
		}, nil
	}
	body, _ := json.Marshal(cities)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func handleDriversLambda(request events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	drivers, err := driverService.ListDrivers()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error":"%s"}`, err.Error()),
		}, nil
	}
	body, _ := json.Marshal(drivers)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
