package zipcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type LocationResponse struct {
	Location string `json:"localidade"`
}

var zipcodeAPIBaseURL = "https://viacep.com.br"

func FetchCityByZipcode(zipCode string) (*LocationResponse, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", zipcodeAPIBaseURL, zipCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not call zipcode api")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("zipcode api returned an error")
	}

	var loc LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, errors.New("could not decode zipcode response")
	}

	return &loc, nil
}
