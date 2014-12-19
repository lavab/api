package factor

import "github.com/GeertJohan/yubigo"

// YubiCloud is an implementation of Factor to authenticate with YubiCloud
type YubiCloud struct {
	client *yubigo.YubiAuth
}

// NewYubiCloud set ups a new Factor that supports authing using YubiCloud
func NewYubiCloud(id string, key string) (*YubiCloud, error) {
	client, err := yubigo.NewYubiAuth(id, key)
	if err != nil {
		return nil, err
	}

	return &YubiCloud{
		client: client,
	}, nil
}

// Type returns factor's type
func (y *YubiCloud) Type() string {
	return "yubicloud"
}

// Request does nothing in this driver
func (y *YubiCloud) Request(data string) (string, error) {
	return "", nil
}

// Verify checks if the token is valid
func (y *YubiCloud) Verify(data []string, input string) (bool, error) {
	publicKey := input[:12]

	found := false
	for _, prefix := range data {
		if publicKey == prefix {
			found = true
			break
		}
	}

	if !found {
		return false, nil
	}

	_, ok, err := y.client.Verify(input)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}
