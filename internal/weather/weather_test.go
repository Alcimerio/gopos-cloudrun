package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchWeatherByCity_Success(t *testing.T) {
	mockResponse := WeatherResponse{
		Current: struct {
			TempC float64 `json:"temp_c"`
		}{
			TempC: 25.0,
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	originalURL := weatherAPIBaseURL
	weatherAPIBaseURL = server.URL
	defer func() { weatherAPIBaseURL = originalURL }()

	location := "Sao Paulo"
	tempResponse, err := FetchWeatherByCity(location)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expectedTempF := 77.0
	expectedTempK := 298.15
	if tempResponse.TempF != expectedTempF {
		t.Errorf("expected TempF to be %v, got %v", expectedTempF, tempResponse.TempF)
	}
	if tempResponse.TempK != expectedTempK {
		t.Errorf("expected TempK to be %v, got %v", expectedTempK, tempResponse.TempK)
	}
}

func TestFetchWeatherByCity_MissingAPIKey(t *testing.T) {
	os.Unsetenv("WEATHER_API_KEY")
	defer os.Setenv("WEATHER_API_KEY", "testkey")
	_, err := FetchWeatherByCity("Sao Paulo")
	if err == nil || err.Error() != "env WEATHER_API_KEY is not set" {
		t.Errorf("expected 'env WEATHER_API_KEY is not set' error, got %v", err)
	}
}

func TestFetchWeatherByCity_APIErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := weatherAPIBaseURL
	weatherAPIBaseURL = server.URL
	defer func() { weatherAPIBaseURL = originalURL }()

	_, err := FetchWeatherByCity("Sao Paulo")
	if err == nil || err.Error() != "weather api returned an error" {
		t.Errorf("expected 'weather api returned an error', got %v", err)
	}
}

func TestFetchWeatherByCity_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	originalURL := weatherAPIBaseURL
	weatherAPIBaseURL = server.URL
	defer func() { weatherAPIBaseURL = originalURL }()

	_, err := FetchWeatherByCity("Sao Paulo")
	if err == nil || err.Error() != "could not decode weather response" {
		t.Errorf("expected 'could not decode weather response', got %v", err)
	}
}
