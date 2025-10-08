package main

import (
	"log"

	"github.com/maxviazov/excel-flow/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("Application failed:", err)
	}
}