package utils

/*
import (
	"fmt"
	"time"
)

// GetJWT TODO
func GetJWT(username string) (string, error) {
	if _, err := FetchUser(username); err == nil {
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["username"] = username
		token.Claims["expires"] = time.Now().Add(time.Hour * 80).Unix()
	} else {
		return fmt.Errorf("User not found in database or connection broken")
	}
	return ""
}
*/
