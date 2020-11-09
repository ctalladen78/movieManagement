package util

import (
	"strings"
)

// FormatDate returns date format
func FormatDate(dateTime string) string {
	dateString := ""
	dateArray := strings.Split(dateTime, " ")
	if len(dateArray) > 1 {
		dateString = dateArray[0]
	}
	return dateString
}
