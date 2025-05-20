package errors

// Error constants used throughout the parking lot system
const (
	// Location related errors
	ErrInvalidLocation = "invalid parking spot location: index out of bounds"
	ErrInvalidSpotID   = "invalid spot ID format: must be floor-row-column"

	// Configuration related errors
	ErrInvalidSpotType = "invalid spot type: must be B-1, M-1, A-1, or X-0"

	// Vehicle related errors
	ErrInvalidVehicleType   = "invalid vehicle type: must be Bicycle, Motorcycle, or Automobile"
	ErrVehicleAlreadyParked = "vehicle is already parked"
	ErrVehicleNotParked     = "vehicle is not currently parked"
	ErrVehicleNotAtSpot     = "vehicle is not parked at the specified spot"

	// Availability related errors
	ErrNoAvailableSpot = "no available parking spot for the specified vehicle type"
)
