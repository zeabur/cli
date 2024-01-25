package util

import (
	"fmt"
	"time"
)

func ConvertTimeAgoString(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration.Hours() > 24:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d day(s) ago", days)
	case duration.Minutes() > 60:
		hours := int(duration.Minutes() / 60)
		return fmt.Sprintf("%d hour(s) ago", hours)
	default:
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d minute(s) ago", minutes)
	}
}
