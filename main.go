package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func getData() table.Model {
	stopId := "A=1@O=U Hallesches Tor (Berlin)@X=13391761@Y=52497777@U=86@L=900012103@"
	departureResp := getDepartureBoardForStop(stopId)

	// Get Berlin timezone
	berlinLoc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal("Error loading timezone:", err)
	}

	// Current time in Berlin timezone, zeros out seconds and nanoseconds
	now := time.Now().In(berlinLoc)
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, berlinLoc)

	// Process departures
	nameDirectionMap := make(map[string]*GroupedDeparture)

	for _, dep := range *departureResp.Departure {
		key := dep.Name + "|" + *dep.Direction

		timeDiff := getDiff(now, dep.Time)

		if group, exists := nameDirectionMap[key]; exists {
			group.TimeDiffMins = append(group.TimeDiffMins, timeDiff)
		} else {
			nameDirectionMap[key] = &GroupedDeparture{
				Name:         dep.Name,
				Direction:    *dep.Direction,
				TimeDiffMins: []int{timeDiff},
			}
		}
	}

	// Sort and deduplicate time differences
	for _, group := range nameDirectionMap {
		sort.Ints(group.TimeDiffMins)

		// Deduplicate times (equivalent to 'unique' in pandas)
		if len(group.TimeDiffMins) > 1 {
			uniqueTimes := []int{group.TimeDiffMins[0]}
			for i := 1; i < len(group.TimeDiffMins); i++ {
				if group.TimeDiffMins[i] != group.TimeDiffMins[i-1] {
					uniqueTimes = append(uniqueTimes, group.TimeDiffMins[i])
				}
			}
			group.TimeDiffMins = uniqueTimes
		}
	}

	columns := []table.Column{
		{Title: "Stop", Width: 20},
		{Title: "Line", Width: 4},
		{Title: "Direction", Width: 30},
		{Title: "Departures (min)", Width: 18},
	}

	rows := []table.Row{}
	for _, group := range nameDirectionMap {
		rows = append(rows, table.Row{
			"U Hallesches Tor",
			group.Name,
			group.Direction,
			intsToStringLimited(group.TimeDiffMins, 5),
		})
	}

	return sortTable(rows, columns)
}

var conf map[string]string

const (
	CONF_API_URL     = "HAFAS_API_URL"
	CONF_API_KEY     = "HAFAS_API_KEY"
	CONF_COORDINATES = "COORDINATES"
)

func loadConfig(conf *map[string]string) {
	path := ".conf.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("Config file not found:", path)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&conf); err != nil {
		log.Fatal("Error decoding config file:", err)
	}

	for _, k := range []string{CONF_API_URL, CONF_API_KEY, CONF_COORDINATES} {
		if _, ok := (*conf)[k]; !ok {
			log.Fatal("Missing required config key:", k)
		}
	}
}

func main() {
	conf = map[string]string{}
	loadConfig(&conf)

	m := model{getData(), 20, 20}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
