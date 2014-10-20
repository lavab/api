package utils

import (
	"crypto/rand"
	"log"
	"os"

	"code.google.com/p/go-uuid/uuid"
)

// RandomString returns a secure random string of a certain length
func RandomString(length int) string {
	tmp := make([]byte, length)
	_, err := rand.Read(tmp)
	if err != nil {
		log.Fatalln("Secure random string generation has failed.", err)
	}
	return string(tmp)
}

// FileExists is a stupid little wrapper of os.Stat that checks whether a file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

// UUID returns a new Universally Unique IDentifier (UUID)
func UUID() string {
	return uuid.New()
}
