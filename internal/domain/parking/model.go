package parking

import (
	"fmt"
	"sync"
)

const (
	Bicycle    = "Bicycle"
	Motorcycle = "Motorcycle"
	Automobile = "Automobile"
)

// ParkingSpotType represents the type of parking spot
type ParkingSpotType struct {
	VehicleType string
	IsActive    bool
}

// ParkingSpot represents a single parking spot
type ParkingSpot struct {
	Floor         int
	Row           int
	Column        int
	Type          ParkingSpotType
	IsOccupied    bool
	VehicleNumber string
}

// SpotID returns the ID of the parking spot in format "floor-row-column"
func (p *ParkingSpot) SpotID() string {
	return fmt.Sprintf("%d-%d-%d", p.Floor, p.Row, p.Column)
}

// ParkingLot represents the entire parking lot system
type ParkingLot struct {
	Floors  int
	Rows    int
	Columns int
	Spots   [][][]*ParkingSpot
	Gates   int

	// For thread safety
	mutex sync.RWMutex

	// For tracking vehicles and their history
	vehicleMap     map[string]string // vehicleNumber -> current spotID
	vehicleHistory map[string]string // vehicleNumber -> last spotID
}
