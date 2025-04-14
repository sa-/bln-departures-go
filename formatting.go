package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func intsToStringLimited(nums []int, limit int) string {
	count := min(len(nums), limit)
	if count <= 0 {
		return ""
	}

	var sb strings.Builder

	for i := range count {
		sb.WriteString(strconv.Itoa(nums[i]))

		// Add a comma and space separator if it's not the last item
		if i < count-1 {
			sb.WriteString(", ")
			if nums[i] < 10 && nums[i] >= 0 {
				sb.WriteString(" ")
			}
		}
	}

	return sb.String()
}

// Function to calculate time difference in minutes
func getDiff(now time.Time, departureTimeStr string) int {
	parsedTime, _ := time.Parse("15:04:05", departureTimeStr)

	// Create a datetime with today's date and the parsed time
	depTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
		0, now.Location(),
	)

	// If the departure time is earlier than now, it's for tomorrow
	if depTime.Sub(now).Minutes() < -60 {
		depTime = depTime.Add(24 * time.Hour)
	}

	// Calculate difference in minutes
	return int(depTime.Sub(now).Minutes())
}

func sortTable(rows []table.Row, columns []table.Column) table.Model {
	sort.Slice(rows, func(i, j int) bool {
		// 1. Compare by Stop (index 0)
		if rows[i][0] != rows[j][0] {
			return rows[i][0] < rows[j][0] // Sort alphabetically by stop name
		}

		// 2. If Stops are equal, compare by Line (index 1)
		if rows[i][1] != rows[j][1] {
			// You might want custom logic for line sorting (e.g., U1 before U10)
			// Simple string comparison works for basic cases like U1, U2, U3...
			return rows[i][1] < rows[j][1] // Sort alphabetically/numerically by line name
		}

		// 3. If Stops and Lines are equal, compare by Direction (index 2)
		return rows[i][2] < rows[j][2] // Sort alphabetically by direction
	})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	return t
}

func padOrTruncate(s string, length int) string {
	if len(s) > length {
		return s[:length-1] + "…"
	}
	return s + strings.Repeat(" ", length-len(s))
}
