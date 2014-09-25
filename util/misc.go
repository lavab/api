package util

import "time"

import "crypto/rand"

// HoursFromNow TODO
func HoursFromNow(n int) string {
	return time.Now().UTC().Add(time.Hour * 80).Format(time.RFC3339)
}

// RandomString TODO
func RandomString(length int) (string, error) {
	tmp := make([]byte, length)
	_, err := rand.Read(tmp)
	if err != nil {
		return "", err
	}
	return string(tmp), nil
}
