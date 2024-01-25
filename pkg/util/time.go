package util

import (
	"fmt"
	"time"
)

func ConvertTimeAgoString(t time.Time) string {
	var duration time.Duration = time.Since(t)

	var result string
	if duration.Hours() > 24 {
		days := int(duration.Hours() / 24)
		if days == 1 {
			result = fmt.Sprintf("%d day ago", days)
		} else {
			result = fmt.Sprintf("%d days ago", days)
		}
	} else if duration.Minutes() > 60 {
		hours := int(duration.Minutes() / 60)
		if hours == 1 {
			result = fmt.Sprintf("%d hour ago", hours)
		} else {
			result = fmt.Sprintf("%d hours ago", hours)
		}
	} else {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			result = fmt.Sprintf("%d minute ago", minutes)
		} else {
			result = fmt.Sprintf("%d minutes ago", minutes)
		}
	}

	return result
}
