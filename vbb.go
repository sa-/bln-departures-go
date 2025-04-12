package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	hafas "github.com/sa-/schedule/hafasClient"
)

func headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + conf[CONF_API_KEY],
		"Accept":        "application/json",
	}
}

func route(path string) string {
	base, _ := url.Parse(conf[CONF_API_URL])
	ref, _ := url.Parse(path)
	return base.ResolveReference(ref).String()
}

func getDepartureBoardForStop(stopID string) *hafas.DepartureBoard {
	client := &http.Client{}

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
	var departureResp *hafas.DepartureBoard
	if err := json.NewDecoder(resp.Body).Decode(&departureResp); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	return departureResp
}

// Define a struct to hold our grouped data
type GroupedDeparture struct {
	Name         string
	Direction    string
	TimeDiffMins []int
}
