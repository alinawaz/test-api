package server

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type postOrder struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type updateOrder struct {
	Status string `json:"status"`
}

// Index to list all orders
func Index(w http.ResponseWriter, r *http.Request) {

	// Getting params
	params := r.URL.Query()

	results := get(params["page"][0], params["limit"][0])
	json.NewEncoder(w).Encode(results)
}

// Store new order
func Store(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	// Parsing Request
	decoder := json.NewDecoder(r.Body)
	var po postOrder
	err := decoder.Decode(&po)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Request Body"})
		return
	}

	// Calculating Distance
	lat1, _ := strconv.ParseFloat(po.Origin[0], 64)
	long1, _ := strconv.ParseFloat(po.Origin[1], 64)
	lat2, _ := strconv.ParseFloat(po.Destination[0], 64)
	long2, _ := strconv.ParseFloat(po.Destination[1], 64)
	var distance = calculateDistance(lat1, long1, lat2, long2) / 1000 //KM

	// Creating Order
	create(distance, "UNASSIGN")
	var row = getLast()

	// Output Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(row[0])
}

// Update order & throw error if already exists
func Update(w http.ResponseWriter, r *http.Request) {

	// Getting params
	params := mux.Vars(r)

	w.Header().Set("content-type", "application/json")

	// Parsing Request
	decoder := json.NewDecoder(r.Body)
	var uo updateOrder
	err := decoder.Decode(&uo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Request Body"})
		return
	}

	// Check for existing order
	var row = getByID(params["id"])
	if len(row) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "ORDER_DOES_NOT_EXISTS"})
	} else {
		if row[0]["status"] == "UNASSIGN" {
			// take order
			updateStatus(params["id"])
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "SUCCESS"})
		} else {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "ORDER_ALREADY_BEEN_TAKEN"})
		}
	}
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
