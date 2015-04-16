package pbkdf2

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/gyepisam/mcf"
	"github.com/gyepisam/mcf/password"
)

var testVectors = []struct {
	plain      string
	salt       string
	iterations int
	key        string
}{
	{
		plain:      "password",
		salt:       "salt",
		iterations: 1,
		key:        "0c60c80f961f0e71f3a9b524af6012062fe037a6",
	},
	{
		plain:      "password",
		salt:       "salt",
		iterations: 2,
		key:        "ea6c014dc72d6f8ccd1ed92ace1d41f0d8de8957",
	},
	{
		plain:      "password",
		salt:       "salt",
		iterations: 4096,
		key:        "4b007901b765489abead49d926f721d065a429c1",
	},
	{
		plain:      "password",
		salt:       "salt",
		iterations: 16777216,
		key:        "eefe3d61cd4da4e4e9945b3d6ba2158c2634e984",
	},
	{
		plain:      "passwordPASSWORDpassword",
		salt:       "saltSALTsaltSALTsaltSALTsaltSALTsalt",
		iterations: 4096,
		key:        "3d2eec4fe41c849b80c8d83662c0e44a8b291a964cf2f07038",
	},
	{
		plain:      "pass\x00word",
		salt:       "sa\x00lt",
		iterations: 4096,
		key:        "56fa6aa75548099dcc37d7f03425e0c3",
	},
	{
		plain:      "password",
		salt:       "ATHENA.MIT.EDUraeburn",
		iterations: 1,
		key:        "cdedb5281bb2f801565a1122b2563515",
	},
	{
		plain:      "password",
		salt:       "ATHENA.MIT.EDUraeburn",
		iterations: 2,
		key:        "01dbee7f4a9e243e988b62c73cda935d",
	},
	{
		plain:      "password",
		salt:       "ATHENA.MIT.EDUraeburn",
		iterations: 1200,
		key:        "5c08eb61fdf71e4e4ec3cf6ba1f5512b",
	},
	{
		plain:      "password",
		salt:       "\x12\x34\x56\x78\x78\x56\x34\x12",
		iterations: 5,
		key:        "d1daa78615f287e6a1c8b120d7062a49",
	},
	{
		plain:      "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		salt:       "pass phrase equals block size",
		iterations: 1200,
		key:        "139c30c0966bc32ba55fdbf212530ac9",
	},
	{
		plain:      "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		salt:       "pass phrase exceeds block size",
		iterations: 1200,
		key:        "9ccad6d468770cd51b10e6a68721be61",
	},
	{
		plain:      "\xf0\x9d\x84\x9e", /* g-clef (0xf09d849e) */
		salt:       "EXAMPLE.COMpianist",
		iterations: 50,
		key:        "6b9cf26d45455a43a5b8bb276a403b39",
	},
}

func setConfig(keylen, iterations, saltlen int) (*Config, error) {
	conf := GetConfig()
	conf.KeyLen = keylen
	conf.Iterations = iterations
	conf.SaltLen = saltlen
	return &conf, SetConfig(conf)
}

func setSalt(salt string) {
	SaltMine = func(n int) ([]byte, error) { return []byte(salt), nil }
}

func TestVectors(t *testing.T) {
	for i, v := range testVectors {

		key, err := hex.DecodeString(v.key)
		if err != nil {
			t.Errorf("%d: could not decode key: %s: %s", i, v.key, err)
			continue
		}

		config, err := setConfig(len(key), v.iterations, len(v.salt))
		if err != nil {
			t.Errorf("%d: unexpected error setting config: %s", i, err)
		}

		setSalt(v.salt)

		encoded, err := mcf.Create(v.plain)
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}

		passwd := password.New([]byte("pbkdf2"))
		err = passwd.Parse([]byte(encoded))
		if err != nil {
			t.Errorf("%d: unexpected error creating password instance: %s", i, err)
		}

		if p, q := []byte(config.Params()), passwd.Params; !bytes.Equal(p, q) {
			t.Errorf("%d: params: expected %s, got %s", i, string(p), string(q))
		}

		if p, q := []byte(v.salt), passwd.Salt; !bytes.Equal(p, q) {
			t.Errorf("%d: salt: expected %s, got %s", i, string(p), string(q))
		}

		if p, q := key, passwd.Key; !bytes.Equal(p, q) {
			t.Errorf("%d: key: expected %x, got %x", i, p, q)
		}

		isValid, err := mcf.Verify(v.plain, encoded)
		if err != nil {
			t.Errorf("%d: verify: unexpected failure on %q: %s", i, encoded, err)
			continue
		}
		if !isValid {
			t.Errorf("%d: verify - unexpectedly returned false", i)
			continue
		}

		// perturb configuration...
		newConfig := *config
		newConfig.KeyLen += 1

		for j, c := range []*Config{config, &newConfig} {
			setConfig(c.KeyLen, c.Iterations, c.SaltLen)
			isCurrent, err := mcf.IsCurrent(encoded)
			if err != nil {
				t.Errorf("%d-%d: IsCurrent: unexpected failure: %", i, j, err)
				continue
			}
			//old configuration says yes, new configuration says no
			if answer := c == config; isCurrent != answer {
				t.Errorf("%d-%d: IsCurrent: expecting %t got %t", i, j, answer, isCurrent)
				continue
			}
		}
	}
}
