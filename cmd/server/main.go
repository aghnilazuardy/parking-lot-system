package main

import (
	"log"
	"parking-lot-system/internal/api/handler"
	"parking-lot-system/internal/config"
	"parking-lot-system/internal/domain/parking"
	"parking-lot-system/internal/repository"
)

func main() {
	// Load configuration
	cfg := config.NewAppConfig()

	parkingRepo := repository.NewParkingRepository()

	parkingService := parking.NewParkingService(parkingRepo)

	// Create a new parking lot with 3 floors, 5 rows, 10 columns, and 2 gates
	err := parkingService.InitializeParkingLot(3, 5, 10, 2)
	if err != nil {
		log.Fatalf("Error creating parking lot: %v\n", err)
	}

	// Configure some spots
	configureSpots := []struct {
		floor    int
		row      int
		column   int
		spotType string
	}{
		{0, 0, 0, "B-1"}, // Bicycle spot
		{0, 0, 1, "B-1"}, // Bicycle spot
		{0, 1, 0, "M-1"}, // Motorcycle spot
		{0, 1, 1, "M-1"}, // Motorcycle spot
		{0, 2, 0, "A-1"}, // Automobile spot
		{0, 2, 1, "A-1"}, // Automobile spot
		{0, 2, 2, "X-0"}, // Inactive spot
		{1, 0, 0, "B-1"}, // Bicycle spot
		{1, 0, 1, "M-1"}, // Motorcycle spot
		{1, 1, 0, "A-1"}, // Automobile spot
	}

	for _, cfg := range configureSpots {
		err := parkingService.ConfigureSpot(cfg.floor, cfg.row, cfg.column, cfg.spotType)
		if err != nil {
			log.Printf("Error configuring spot at (%d,%d,%d): %v\n",
				cfg.floor, cfg.row, cfg.column, err)
		}
	}

	// Create a new handler with the parking service
	parkingHandler := handler.NewParkingHandler(parkingService)

	// Start the HTTP server on port 8080
	log.Fatal(parkingHandler.StartServer(cfg.ServerPort))
}
