package helpers

import "time"

func FormatShowTime(t time.Time) string {
	return t.Format("15:04") // 24-hour format (16:38)
	// For 12-hour format with AM/PM: return t.Format("03:04 PM")
}
