package calendar

import (
	"strconv"
	"time"
)

func YesterdaysDate(location *time.Location) (int, time.Month, int) {
	year, month, day := time.Now().In(location).Date()
	yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, location)
	return yesterday.Year(), yesterday.Month(), yesterday.Day()
}

func PadMonth(month time.Month) string {
	m := strconv.Itoa(int(month))
	if month < 10 {
		return "0" + m
	}
	return m
}
