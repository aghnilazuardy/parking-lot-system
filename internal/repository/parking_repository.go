package repository

import (
	"errors"
	"fmt"
	pkgerrors "parking-lot-system/pkg/errors"
	"sync"
)

// represents a single parking spot in the repository
type ParkingSpot struct {
	Floor         int
	Row           int
	Column        int
	VehicleType   string
	IsActive      bool
	IsOccupied    bool
	VehicleNumber string
}

type ParkingRepository interface {
	InitializeParkingLot(floors, rows, columns, gates int) error
	ConfigureSpot(floor, row, column int, vehicleType string, isActive bool) error
	IsValidLocation(floor, row, column int) bool
	IsSpotOccupied(floor, row, column int) (bool, error)
	FindAvailableSpot(vehicleType string) (string, error)
	ParkVehicle(spotID string, vehicleNumber string) error
	UnparkVehicle(floor, row, column int, vehicleNumber string) error
	IsVehicleParked(vehicleNumber string) (bool, string, error)
	GetAvailableSpots(vehicleType string) ([]string, error)
	SearchVehicle(vehicleNumber string) (string, bool, error)
	ParseSpotID(spotID string) (int, int, int, error)
}

type InMemoryParkingRepository struct {
	floors         int
	rows           int
	columns        int
	spots          [][][]*ParkingSpot
	gates          int
	mutex          sync.RWMutex
	vehicleMap     map[string]string // vehicleNumber -> current spotID
	vehicleHistory map[string]string // vehicleNumber -> last spotID
}

func NewParkingRepository() ParkingRepository {
	return &InMemoryParkingRepository{
		vehicleMap:     make(map[string]string),
		vehicleHistory: make(map[string]string),
	}
}

// InitializeParkingLot creates a new parking lot with the specified dimensions
func (r *InMemoryParkingRepository) InitializeParkingLot(floors, rows, columns, gates int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.floors = floors
	r.rows = rows
	r.columns = columns
	r.gates = gates

	// Initialize parking spots
	r.spots = make([][][]*ParkingSpot, floors)
	for f := 0; f < floors; f++ {
		r.spots[f] = make([][]*ParkingSpot, rows)
		for row := 0; row < rows; row++ {
			r.spots[f][row] = make([]*ParkingSpot, columns)
			for col := 0; col < columns; col++ {
				r.spots[f][row][col] = &ParkingSpot{
					Floor:         f,
					Row:           row,
					Column:        col,
					IsOccupied:    false,
					VehicleType:   "",
					IsActive:      false,
					VehicleNumber: "",
				}
			}
		}
	}

	return nil
}

// ConfigureSpot sets the type and active status of a specific parking spot
func (r *InMemoryParkingRepository) ConfigureSpot(floor, row, column int, vehicleType string, isActive bool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isValidLocation(floor, row, column) {
		return errors.New(pkgerrors.ErrInvalidLocation)
	}

	spot := r.spots[floor][row][column]
	spot.VehicleType = vehicleType
	spot.IsActive = isActive

	return nil
}

// IsValidLocation checks if the location is valid
func (r *InMemoryParkingRepository) IsValidLocation(floor, row, column int) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.isValidLocation(floor, row, column)
}

// isValidLocation is a helper function to check location validity
func (r *InMemoryParkingRepository) isValidLocation(floor, row, column int) bool {
	return floor >= 0 && floor < r.floors &&
		row >= 0 && row < r.rows &&
		column >= 0 && column < r.columns
}

// IsSpotOccupied checks if a spot is occupied
func (r *InMemoryParkingRepository) IsSpotOccupied(floor, row, column int) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.isValidLocation(floor, row, column) {
		return false, errors.New(pkgerrors.ErrInvalidLocation)
	}

	return r.spots[floor][row][column].IsOccupied, nil
}

// FindAvailableSpot finds an available spot for the specified vehicle type
func (r *InMemoryParkingRepository) FindAvailableSpot(vehicleType string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for f := 0; f < r.floors; f++ {
		for row := 0; row < r.rows; row++ {
			for col := 0; col < r.columns; col++ {
				spot := r.spots[f][row][col]

				if spot.IsActive && spot.VehicleType == vehicleType && !spot.IsOccupied {
					// Found an available spot
					return fmt.Sprintf("%d-%d-%d", f, row, col), nil
				}
			}
		}
	}

	return "", errors.New(pkgerrors.ErrNoAvailableSpot)
}

// ParkVehicle parks a vehicle at the specified spot
func (r *InMemoryParkingRepository) ParkVehicle(spotID string, vehicleNumber string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	floor, row, col, err := r.parseSpotID(spotID)
	if err != nil {
		return err
	}

	spot := r.spots[floor][row][col]
	spot.IsOccupied = true
	spot.VehicleNumber = vehicleNumber
	r.vehicleMap[vehicleNumber] = spotID

	return nil
}

// UnparkVehicle removes a vehicle from the specified spot
func (r *InMemoryParkingRepository) UnparkVehicle(floor, row, column int, vehicleNumber string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isValidLocation(floor, row, column) {
		return errors.New(pkgerrors.ErrInvalidLocation)
	}

	spot := r.spots[floor][row][column]

	// Check if the spot is occupied by the specified vehicle
	if !spot.IsOccupied || spot.VehicleNumber != vehicleNumber {
		return fmt.Errorf("%s: %s at spot %d-%d-%d",
			pkgerrors.ErrVehicleNotAtSpot, vehicleNumber, floor, row, column)
	}

	// Unpark the vehicle
	spot.IsOccupied = false
	spot.VehicleNumber = ""

	// Update the vehicle history and remove from current map
	spotID := fmt.Sprintf("%d-%d-%d", floor, row, column)
	r.vehicleHistory[vehicleNumber] = spotID
	delete(r.vehicleMap, vehicleNumber)

	return nil
}

// IsVehicleParked checks if a vehicle is currently parked
func (r *InMemoryParkingRepository) IsVehicleParked(vehicleNumber string) (bool, string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	spotID, exists := r.vehicleMap[vehicleNumber]
	return exists, spotID, nil
}

// GetAvailableSpots returns the list of available spots for a vehicle type
func (r *InMemoryParkingRepository) GetAvailableSpots(vehicleType string) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	availableSpots := []string{}

	for f := 0; f < r.floors; f++ {
		for row := 0; row < r.rows; row++ {
			for col := 0; col < r.columns; col++ {
				spot := r.spots[f][row][col]

				if spot.IsActive && spot.VehicleType == vehicleType && !spot.IsOccupied {
					availableSpots = append(availableSpots, fmt.Sprintf("%d-%d-%d", f, row, col))
				}
			}
		}
	}

	if len(availableSpots) == 0 {
		return nil, fmt.Errorf("%s: %s", pkgerrors.ErrNoAvailableSpot, vehicleType)
	}

	return availableSpots, nil
}

// SearchVehicle returns the current or last known spot ID for a vehicle
func (r *InMemoryParkingRepository) SearchVehicle(vehicleNumber string) (string, bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Check if the vehicle is currently parked
	if spotID, exists := r.vehicleMap[vehicleNumber]; exists {
		return spotID, true, nil
	}

	// Check if we have a history for this vehicle
	if lastSpotID, exists := r.vehicleHistory[vehicleNumber]; exists {
		return lastSpotID, false, nil
	}

	return "", false, fmt.Errorf("vehicle %s has never been parked in this parking lot", vehicleNumber)
}

// ParseSpotID parses a spot ID string into floor, row, column
func (r *InMemoryParkingRepository) ParseSpotID(spotID string) (int, int, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.parseSpotID(spotID)
}

// parseSpotID is a helper function to parse a spot ID
func (r *InMemoryParkingRepository) parseSpotID(spotID string) (int, int, int, error) {
	var floor, row, column int
	_, err := fmt.Sscanf(spotID, "%d-%d-%d", &floor, &row, &column)
	if err != nil {
		return 0, 0, 0, errors.New(pkgerrors.ErrInvalidSpotID)
	}

	// Check if the indices are within bounds
	if !r.isValidLocation(floor, row, column) {
		return 0, 0, 0, errors.New(pkgerrors.ErrInvalidLocation)
	}

	return floor, row, column, nil
}
