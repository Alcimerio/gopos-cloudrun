package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/alcimerio/gopos-cloudrun/pkg/textutils"
)

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type TempResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var weatherAPIBaseURL = "https://api.weatherapi.com/v1/current.json"

func FetchWeatherByCity(location string) (*TempResponse, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("env WEATHER_API_KEY is not set")
	}
	parsedLocation := url.QueryEscape(textutils.RemoveAccents(location))
	log.Printf("Trying to fetch location: %s", parsedLocation)

	url := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", weatherAPIBaseURL, apiKey, parsedLocation)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not call weather api")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("weather api returned an error")
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, errors.New("could not decode weather response")
	}

	tempC := weather.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	return &TempResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
