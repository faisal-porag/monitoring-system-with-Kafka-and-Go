package utils

import (
	"time"
)

func DateFormat(date time.Time) string {
	formattedDate := date.Format("02 Jan, 2006 03:04 PM")
	return formattedDate
}
