package dto

type ParkRequest struct {
	VehicleType   string `json:"vehicleType"`
	VehicleNumber string `json:"vehicleNumber"`
}

type ParkResponse struct {
	SpotID string `json:"spotId,omitempty"`
	Error  string `json:"error,omitempty"`
}

type UnparkRequest struct {
	SpotID        string `json:"spotId"`
	VehicleNumber string `json:"vehicleNumber"`
}

type UnparkResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type AvailableSpotRequest struct {
	VehicleType string `json:"vehicleType"`
}

type AvailableSpotResponse struct {
	Spots []string `json:"spots,omitempty"`
	Error string   `json:"error,omitempty"`
}

type SearchVehicleRequest struct {
	VehicleNumber string `json:"vehicleNumber"`
}

type SearchVehicleResponse struct {
	SpotID    string `json:"spotId,omitempty"`
	IsParked  bool   `json:"isParked"`
	WasParked bool   `json:"wasParked"`
	Error     string `json:"error,omitempty"`
}
