package utils

import (
	"fmt"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
)

const bcryptCost = 12

// BcryptHash TODO
func BcryptHash(src ...string) (string, error) {
	in := strings.Join(src, "")
	out, err := bcrypt.GenerateFromPassword([]byte(in), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hashing has failed")
	}
	return string(out), nil
}

// BcryptVerify TODO
func BcryptVerify(hashed, plain string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil {
		return true
	}
	return false
}
