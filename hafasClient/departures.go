package hafasClient

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	hg "github.com/sa-/schedule/hafasClient/gen"
)

type Departure struct {
	StopName  string
	Name      string
	Direction string
	Platform  string
	Time      string
}

func getDepartureBoardForStop(stopID string) *hg.DepartureBoard {

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
	var departureResp *hg.DepartureBoard
	if err := json.NewDecoder(resp.Body).Decode(&departureResp); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	return departureResp
}

func GetDeparturesForStop(stopID, stopName string) []Departure {
	departureBoard := getDepartureBoardForStop(stopID)
	if departureBoard == nil {
		return []Departure{}
	}

	if departureBoard.Departure == nil {
		return []Departure{}
	}

	departures := make([]Departure, len(*departureBoard.Departure))

	for i, dep := range *departureBoard.Departure {
		platform := "-"
		if dep.Platform != nil {
			platform = *dep.Platform.Text
		}
		departures[i] = Departure{
			StopName:  tuncateParenths(stopName),
			Name:      dep.Name,
			Direction: *dep.Direction,
			Platform:  tuncateParenths(platform),
			Time:      dep.Time,
		}
	}

	return departures
}

func tuncateParenths(s string) string {
	// Remove everything after the first parenthesis
	if i := len(s); i > 2 {
		for j := i - 1; j >= 0; j-- {
			if s[j] == '(' {
				return s[:j]
			}
		}
	}
	return strings.TrimSpace(s)
}
