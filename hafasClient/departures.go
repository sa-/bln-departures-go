package hafasClient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	conf "github.com/sa-/schedule/conf"
)

func headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + conf.Conf.VbbAPIKey,
		"Accept":        "application/json",
	}
}

func route(path string) string {
	base, _ := url.Parse(conf.Conf.VbbApiUrl)
	ref, _ := url.Parse(path)
	return base.ResolveReference(ref).String()
}

var client = &http.Client{}

func GetDepartureBoardForStop(stopID string) *DepartureBoard {

	// Create request
	req, err := http.NewRequest("GET", route("departureBoard"), nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Add headers
	for key, value := range headers() {
		req.Header.Add(key, value)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("id", stopID)
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	var departureResp *DepartureBoard
	if err := json.NewDecoder(resp.Body).Decode(&departureResp); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	return departureResp
}

func GetStationsNearCoordinates() *LocationList {

	// Create request
	req, err := http.NewRequest("GET", route("location.nearbystops"), nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Add headers
	for key, value := range headers() {
		req.Header.Add(key, value)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("originCoordLat", conf.Conf.Latitude)
	q.Add("originCoordLong", conf.Conf.Longitude)
	q.Add("accessId", conf.Conf.VbbAPIKey)
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}
	var stationsResp *LocationList
	if err := json.Unmarshal(bodyBytes, &stationsResp); err != nil {
		log.Printf("Error decoding response: %s", err)
	}

	return stationsResp
}
