package parking

import (
	"errors"
	"fmt"
	"parking-lot-system/internal/repository"
	pkgerrors "parking-lot-system/pkg/errors"
)

type ParkingService struct {
	repo repository.ParkingRepository
}

func NewParkingService(repo repository.ParkingRepository) *ParkingService {
	return &ParkingService{
		repo: repo,
	}
}

// InitializeParkingLot creates a new parking lot with the specified dimensions
func (s *ParkingService) InitializeParkingLot(floors, rows, columns, gates int) error {
	// Validate inputs
	if floors < 1 || floors > 8 {
		return errors.New("floors must be between 1 and 8")
	}
	if rows < 1 || rows > 1000 {
		return errors.New("rows must be between 1 and 1000")
	}
	if columns < 1 || columns > 1000 {
		return errors.New("columns must be between 1 and 1000")
	}
	if gates < 1 {
		return errors.New("gates must be at least 1")
	}

	return s.repo.InitializeParkingLot(floors, rows, columns, gates)
}

// ConfigureSpot sets the type and active status of a specific parking spot
func (s *ParkingService) ConfigureSpot(floor, row, column int, spotType string) error {
	// Validate location indices
	if !s.repo.IsValidLocation(floor, row, column) {
		return errors.New(pkgerrors.ErrInvalidLocation)
	}

	// Check if spot is occupied
	isOccupied, err := s.repo.IsSpotOccupied(floor, row, column)
	if err != nil {
		return err
	}

	if isOccupied {
		return errors.New("cannot reconfigure an occupied parking spot")
	}

	// Validate and set spot type
	var vehicleType string
	var isActive bool

	switch spotType {
	case "B-1":
		vehicleType = Bicycle
		isActive = true
	case "M-1":
		vehicleType = Motorcycle
		isActive = true
	case "A-1":
		vehicleType = Automobile
		isActive = true
	case "X-0":
		vehicleType = ""
		isActive = false
	default:
		return errors.New(pkgerrors.ErrInvalidSpotType)
	}

	return s.repo.ConfigureSpot(floor, row, column, vehicleType, isActive)
}

// Park assigns a parking spot to a vehicle
func (s *ParkingService) Park(vehicleType, vehicleNumber string) (string, error) {
	// Validate inputs
	if err := s.validateVehicleType(vehicleType); err != nil {
		return "", err
	}

	if err := s.validateVehicleNumber(vehicleNumber); err != nil {
		return "", err
	}

	// Check if vehicle is already parked
	isParked, currentSpotID, _ := s.repo.IsVehicleParked(vehicleNumber)
	if isParked {
		return "", fmt.Errorf("%s: %s at spot %s", pkgerrors.ErrVehicleAlreadyParked, vehicleNumber, currentSpotID)
	}

	// Find an available spot
	spotID, err := s.repo.FindAvailableSpot(vehicleType)
	if err != nil {
		return "", errors.New(pkgerrors.ErrNoAvailableSpot)
	}

	// Park the vehicle
	err = s.repo.ParkVehicle(spotID, vehicleNumber)
	if err != nil {
		return "", err
	}

	return spotID, nil
}

// Unpark removes a vehicle from its parking spot
func (s *ParkingService) Unpark(spotID, vehicleNumber string) error {
	// Validate inputs
	if err := s.validateVehicleNumber(vehicleNumber); err != nil {
		return err
	}

	// Check if the vehicle is currently parked
	isParked, currentSpotID, err := s.repo.IsVehicleParked(vehicleNumber)
	if err != nil {
		return err
	}

	if !isParked {
		return fmt.Errorf("%s: %s", pkgerrors.ErrVehicleNotParked, vehicleNumber)
	}

	// Check if the vehicle is at the specified spot
	if currentSpotID != spotID {
		return fmt.Errorf("%s: %s (expected: %s, actual: %s)",
			pkgerrors.ErrVehicleNotAtSpot, vehicleNumber, spotID, currentSpotID)
	}

	// Parse and validate spotID
	floor, row, column, err := s.repo.ParseSpotID(spotID)
	if err != nil {
		return err
	}

	// Unpark the vehicle
	return s.repo.UnparkVehicle(floor, row, column, vehicleNumber)
}

// GetAvailableSpots returns the list of available spots for a vehicle type
func (s *ParkingService) GetAvailableSpots(vehicleType string) ([]string, error) {
	// Validate inputs
	if err := s.validateVehicleType(vehicleType); err != nil {
		return nil, err
	}

	return s.repo.GetAvailableSpots(vehicleType)
}

// SearchVehicle returns the current or last known spot ID for a vehicle
func (s *ParkingService) SearchVehicle(vehicleNumber string) (string, bool, error) {
	// Validate inputs
	if err := s.validateVehicleNumber(vehicleNumber); err != nil {
		return "", false, err
	}

	return s.repo.SearchVehicle(vehicleNumber)
}

// validateVehicleType checks if the vehicle type is valid
func (s *ParkingService) validateVehicleType(vehicleType string) error {
	switch vehicleType {
	case Bicycle, Motorcycle, Automobile:
		return nil
	default:
		return errors.New(pkgerrors.ErrInvalidVehicleType)
	}
}

// validateVehicleNumber checks if the vehicle number is valid
func (s *ParkingService) validateVehicleNumber(vehicleNumber string) error {
	if vehicleNumber == "" {
		return errors.New("vehicle number cannot be empty")
	}
	return nil
}
