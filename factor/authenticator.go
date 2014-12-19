package factor

import (
	"github.com/gokyle/hotp"
)

type Authenticator struct {
	length int
}

func NewAuthenticator(length int) *Authenticator {
	return &Authenticator{
		length: length,
	}
}

func (a *Authenticator) Type() string {
	return "authenticator"
}

func (a *Authenticator) Request(data string) (string, error) {
	otp, err := hotp.GenerateHOTP(a.length, false)
	if err != nil {
		return "", err
	}

	return otp.URL(data), nil
}

func (a *Authenticator) Verify(data []string, input string) (bool, error) {
	// obviously broken
	hotp, err := hotp.Unmarshal([]byte(data[0]))
	if err != nil {
		return false, err
	}

	return hotp.Check(input), nil
}
