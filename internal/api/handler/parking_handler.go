package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"parking-lot-system/internal/api/dto"
	"parking-lot-system/internal/domain/parking"
)

type ParkingHandler struct {
	service *parking.ParkingService
}

func NewParkingHandler(service *parking.ParkingService) *ParkingHandler {
	return &ParkingHandler{service: service}
}

// Error response helper
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// handles the POST /park endpoint

/** cURL example
curl -X POST http://localhost:8080/park \
     -H "Content-Type: application/json" \
     -d '{"vehicleType": "Bicycle", "vehicleNumber": "BC001"}'
**/

func (h *ParkingHandler) handlePark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var req dto.ParkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	spotID, err := h.service.Park(req.VehicleType, req.VehicleNumber)
	resp := dto.ParkResponse{}

	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp.SpotID = spotID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handles the POST /unpark endpoint

/** cURL example
curl -X POST http://localhost:8080/unpark \
     -H "Content-Type: application/json" \
     -d '{"spotId": "0-0-1", "vehicleNumber": "BC001"}'
**/

func (h *ParkingHandler) handleUnpark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var req dto.UnparkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	err := h.service.Unpark(req.SpotID, req.VehicleNumber)
	resp := dto.UnparkResponse{}

	if err != nil {
		resp.Success = false
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp.Success = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handles the GET /available endpoint

/** cURL example
curl -X GET "http://localhost:8080/available?vehicleType=Bicycle"
**/

func (h *ParkingHandler) handleAvailableSpots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	vehicleType := r.URL.Query().Get("vehicleType")
	if vehicleType == "" {
		writeErrorResponse(w, http.StatusBadRequest, "vehicleType query parameter is required")
		return
	}

	spots, err := h.service.GetAvailableSpots(vehicleType)
	resp := dto.AvailableSpotResponse{}

	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp.Spots = spots
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handles the GET /search endpoint

/** cURL example
curl -X GET "http://localhost:8080/search?vehicleNumber=BC001"
**/

func (h *ParkingHandler) handleSearchVehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	vehicleNumber := r.URL.Query().Get("vehicleNumber")
	if vehicleNumber == "" {
		writeErrorResponse(w, http.StatusBadRequest, "vehicleNumber query parameter is required")
		return
	}

	spotID, isParked, err := h.service.SearchVehicle(vehicleNumber)
	resp := dto.SearchVehicleResponse{}

	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp.SpotID = spotID
		resp.IsParked = isParked
		resp.WasParked = spotID != ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// registers all the API routes
func (h *ParkingHandler) registerRoutes() {
	http.HandleFunc("/park", h.handlePark)
	http.HandleFunc("/unpark", h.handleUnpark)
	http.HandleFunc("/available", h.handleAvailableSpots)
	http.HandleFunc("/search", h.handleSearchVehicle)
}

// starts the HTTP server on the specified port
func (h *ParkingHandler) StartServer(port int) error {
	h.registerRoutes()

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting parking lot API server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
