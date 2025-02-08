package malak

import (
	"fmt"
	"time"
)

func ordinalSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func GetTodayFormatted() string {

	today := time.Now()
	day := today.Day()
	formattedDate := fmt.Sprintf("%s %d%s, %d",
		today.Format("January"), day, ordinalSuffix(day), today.Year())

	return formattedDate
}
