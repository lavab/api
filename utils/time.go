package utils

import "time"

// TimeNowString returns time.Now + UTC + Format(RFC3339). This makes the strings comparable.
func TimeNowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// HoursFromNowString returns time.Now + n hours in string format. The result is RFC3339
// encoded and in UTC. This makes it possible to compare the times with a simple
// string comparison. Example: "2006-12-30T15:34:45Z00:00"
func HoursFromNowString(n int) string {
	return time.Now().UTC().Add(time.Hour * time.Duration(n)).Format(time.RFC3339)
}

// StringToTime returns a time.Time object from a RFC3339-encoded string.
func StringToTime(in string) (time.Time, error) {
	return time.Parse(time.RFC3339, in)
}
