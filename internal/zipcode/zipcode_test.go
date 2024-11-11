package zipcode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchCityByZipcode_Success(t *testing.T) {
	mockResponse := LocationResponse{Location: "São Paulo"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	originalURL := zipcodeAPIBaseURL
	zipcodeAPIBaseURL = server.URL
	defer func() { zipcodeAPIBaseURL = originalURL }()

	loc, err := FetchCityByZipcode("01001000")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if loc.Location != "São Paulo" {
		t.Errorf("expected location to be São Paulo, got %v", loc.Location)
	}
}

func TestFetchCityByZipcode_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := zipcodeAPIBaseURL
	zipcodeAPIBaseURL = server.URL
	defer func() { zipcodeAPIBaseURL = originalURL }()

	_, err := FetchCityByZipcode("01001000")
	if err == nil || err.Error() != "zipcode api returned an error" {
		t.Errorf("expected 'zipcode api returned an error', got %v", err)
	}
}

func TestFetchCityByZipcode_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	originalURL := zipcodeAPIBaseURL
	zipcodeAPIBaseURL = server.URL
	defer func() { zipcodeAPIBaseURL = originalURL }()

	_, err := FetchCityByZipcode("01001000")
	if err == nil || err.Error() != "could not decode zipcode response" {
		t.Errorf("expected 'could not decode zipcode response', got %v", err)
	}
}
