package resources

import (
	"strconv"
	"time"
)

func YesterdaysMonthAndYear(location *time.Location) (int, time.Month) {
	year, month, day := time.Now().In(location).Date()
	yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, location)
	return yesterday.Year(), yesterday.Month()
}

func PadMonth(month time.Month) string {
	m := strconv.Itoa(int(month))
	if month < 10 {
		return "0" + m
	}
	return m
}
