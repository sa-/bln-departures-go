package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sa-/schedule/conf"
	"github.com/sa-/schedule/hafasClient"
	ms "github.com/sa-/schedule/meteoSource"
)

// Define a struct to hold our grouped data
type GroupedDeparture struct {
	Name         string
	Direction    string
	Platform     string
	TimeDiffMins []int
}

func makeV1Table(departureBoard *hafasClient.DepartureBoard) table.Model {
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

	for _, dep := range *departureBoard.Departure {
		key := dep.Name + "|" + *dep.Direction

		timeDiff := getDiff(now, dep.Time)

		if group, exists := nameDirectionMap[key]; exists {
			group.TimeDiffMins = append(group.TimeDiffMins, timeDiff)
		} else {
			nameDirectionMap[key] = &GroupedDeparture{
				Name:         dep.Name,
				Direction:    *dep.Direction,
				Platform:     *dep.Platform.Text,
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
		{Title: "Platform", Width: 9},
		{Title: "Direction", Width: 25},
		{Title: "Departures (min)", Width: 18},
	}

	rows := []table.Row{}
	for _, group := range nameDirectionMap {
		rows = append(rows, table.Row{
			"U Hallesches Tor",
			group.Name,
			group.Platform,
			group.Direction,
			intsToStringLimited(group.TimeDiffMins, 5),
		})
	}

	return sortTable(rows, columns)
}

func makeWeather(weatherData *ms.PointPointData) (viewport.Model, viewport.Model) {
	hourlySb := strings.Builder{}
	current := fmt.Sprintf(
		"Now: %d°C %s",
		int(*weatherData.Current.Temperature),
		*weatherData.Current.Summary,
	)
	hourlySb.WriteString(current + "\n")

	for _, d := range weatherData.Hourly.Data {
		hour := (*d.Date)[11:13]

		hourlySb.WriteString(fmt.Sprintf(
			"%sh %d°C %s\n",
			hour,
			int(*d.Temperature),
			ms.GetEmoji(*d.Icon),
		))
	}

	dailySb := strings.Builder{}
	for _, d := range weatherData.Daily.Data {
		minTempInt := int(*d.AllDay.TemperatureMin)
		minTempStr := strconv.Itoa(minTempInt)
		if minTempInt >= 0 && minTempInt < 10 {
			minTempStr = " " + minTempStr
		}

		maxTempInt := int(*d.AllDay.TemperatureMax)
		maxTempStr := strconv.Itoa(maxTempInt)
		if maxTempInt >= 0 && maxTempInt < 10 {
			maxTempStr = " " + maxTempStr
		}

		bold := lipgloss.NewStyle().Bold(true)

		dailySb.WriteString(fmt.Sprintf(
			"%s\n\t%s°C %s°C %.2f %s\n\t%s\n",
			bold.Render((*d.Day)[5:10]),
			minTempStr,
			maxTempStr,
			*d.AllDay.Precipitation.Total,
			*d.AllDay.Precipitation.Type,
			padOrTruncate(*d.Summary, 45),
		))
	}

	hourlyVp := viewport.New(50, 20)
	hourlyVp.SetContent(hourlySb.String())

	dailyVp := viewport.New(50, 20)
	dailyVp.SetContent(dailySb.String())
	return hourlyVp, dailyVp
}

type NearbyStop struct {
	Name     string
	Id       string
	Distance int
}

func main() {
	conf.LoadConfig()

	stopId := "A=1@O=U Hallesches Tor (Berlin)@X=13391761@Y=52497777@U=86@L=900012103@"
	appStateDepartureBoard = hafasClient.GetDepartureBoardForStop(stopId)

	departureTable := makeV1Table(appStateDepartureBoard)

	weatherData := ms.GetResponse()
	hourlyVp, dailyVp := makeWeather(weatherData)

	locs := hafasClient.GetStationsNearCoordinates()
	stops := strings.Builder{}
	for _, l := range *locs.StopLocationOrCoordLocation {
		s := l["StopLocation"].(map[string]interface{})
		ns := NearbyStop{
			Name:     s["name"].(string),
			Id:       s["id"].(string),
			Distance: int(s["dist"].(float64)),
		}
		stops.WriteString(fmt.Sprintf("Name: %s\tDist: %d\tId: %s", ns.Name, ns.Distance, ns.Id) + "\n")
	}

	altVp := viewport.New(50, 20)
	altVp.SetContent(stops.String())

	m := model{20, 20, 0, departureTable, hourlyVp, dailyVp, altVp}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
