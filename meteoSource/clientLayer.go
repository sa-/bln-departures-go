package meteoSource

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sa-/schedule/conf"
)

func route(path string) string {
	base, _ := url.Parse("https://www.meteosource.com/api/v1/free/")
	ref, _ := url.Parse(path)
	return base.ResolveReference(ref).String()
}

func GetResponse() *PointPointData {
	client := &http.Client{}

	req, err := http.NewRequest("GET", route("point"), nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")

	// Add query parameters
	coords := strings.Split(conf.Conf.Coordinates, ",")
	q := req.URL.Query()
	q.Add("lat", coords[0])
	q.Add("lon", coords[1])
	q.Add("key", conf.Conf.MeteosourceApiKey)
	q.Add("sections", "all")
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}

	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	// Write bodyBytes to a file
	filePath := ".scratch/meteo.json"
	var formattedBytes *bytes.Buffer
	formattedBytes = new(bytes.Buffer)
	json.Indent(formattedBytes, bodyBytes, "", "  ")
	if err := os.WriteFile(filePath, formattedBytes.Bytes(), 0644); err != nil {
		log.Fatalf("Error writing to file %s: %v", filePath, err)
	}

	bodyString := string(formattedBytes.Bytes())

	// Decode JSON response
	var pointResp *PointPointData
	if err := json.Unmarshal([]byte(bodyString), &pointResp); err != nil {
		log.Printf("Error decoding response: %v", err)
	}

	return pointResp
}
