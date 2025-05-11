package hafasClient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	conf "github.com/sa-/schedule/conf"
	hg "github.com/sa-/schedule/hafasClient/gen"
)

type NearbyStop struct {
	Name     string
	Id       string
	Distance int
}

func GetStationsNearCoordinates() *[]NearbyStop {

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
	q.Add("maxNo", "20")
	q.Add("r", "650")

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
	var stationsResp *hg.LocationList
	if err := json.Unmarshal(bodyBytes, &stationsResp); err != nil {
		log.Printf("Error decoding response: %s", err)
	}

	nearbyStops := make([]NearbyStop, len(*stationsResp.StopLocationOrCoordLocation))

	for i, l := range *stationsResp.StopLocationOrCoordLocation {
		s := l["StopLocation"].(map[string]any)
		nearbyStops[i] = NearbyStop{
			Name:     s["name"].(string),
			Id:       s["id"].(string),
			Distance: int(s["dist"].(float64)),
		}
	}

	return &nearbyStops
}
