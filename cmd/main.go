package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alcimerio/gopos-cloudrun/internal/weather"
	"github.com/alcimerio/gopos-cloudrun/internal/zipcode"
)

func weatherByZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("zipcode")
	if len(zipCode) != 8 || !isNumeric(zipCode) {
		log.Printf("Invalid zipcode: %s", zipCode)
		http.Error(w, "Invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	loc, err := zipcode.FetchCityByZipcode(zipCode)
	if err != nil {
		log.Printf("Error fetching zipcode: %v", err)
		http.Error(w, "Zipcode not found", http.StatusNotFound)
		return
	}

	tempResponse, err := weather.FetchWeatherByCity(loc.Location)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		http.Error(w, "Not possible to find weather", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tempResponse)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func main() {
	http.HandleFunc("/", weatherByZipCodeHandler)
	http.ListenAndServe(":8080", nil)
}
