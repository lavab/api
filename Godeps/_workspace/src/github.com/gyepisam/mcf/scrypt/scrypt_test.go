package scrypt

import (
	"bytes"
	"github.com/gyepisam/mcf"
	"github.com/gyepisam/mcf/password"

	"testing"
)

var testData = []struct {
	plaintext string
	encoded   string
}{
	{"yssIi0AL",
		"$scrypt$KeyLen=128,N=65536,R=10,P=2$PmxwHoNHjIILwrdOG8vA+A==$KRYMgbJr4vrYutrEjtueDDylXHQ2EoePyPoqtrDnil0jm1RfuyxT90/3gce5hw0/DTVmpcDnzkt1MWiK+zfS7+hh1EONspTZl8nGFLCsXcGiarrNKSSyRnsJN0DmSe20cfxodAB+DN1f84hxZbdmF20A2uFj36kE2ZKgTdlAYFE="},
	{"yssIi0AL", "$scrypt$KeyLen=128,N=16384,R=8,P=1$643873597251626754$e6b2da790d99bef794d6feb8ab7fda61f44251303da936ad75162454fe3017e80302ed0fdf654ade552906819aee1370278e852aa2ee9ca9b91e934c6337d6607c94a98e6c06b7279cbfc30d9c783d30e9958ae601dd90d57dcc162eebc513164bd717abe23e7b8281cc13865e0d7453ccd36ff9dbfcb7aadf2451da926c8413"},
}

var plaintext = "g5Dr58dvyD"
var salt = []byte("d8sYrQbgT")

func roundTrip(t *testing.T, plaintext string) {
	encoded, err := mcf.Create(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := mcf.Verify(plaintext, encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatalf("Verify(%q, %q) failed", plaintext, encoded)
	}
}

func TestVerificationExisting(t *testing.T) {
	for _, slot := range testData {
		valid, err := mcf.Verify(slot.plaintext, slot.encoded)
		if err != nil {
			panic(err)
		}
		if !valid {
			t.Fatalf("#1 Verify(%q, %q) failed", slot.plaintext, slot.encoded)
		}
	}
}

func TestRoundtrip(t *testing.T) {
	roundTrip(t, plaintext)
}

func setConfig(k, s, N, r, p int) error {
	conf := GetConfig()
	conf.KeyLen = k
	conf.SaltLen = s
	conf.N = N
	conf.R = r
	conf.P = p

	return SetConfig(conf)
}

func TestCustomParameters(t *testing.T) {

	defaultConf := GetConfig()
	conf := defaultConf

	conf.KeyLen = 128
	conf.N = 1 << 14
	conf.R = 8
	conf.P = 1

	defaultParams := (&defaultConf).Params()
	newParams := (&conf).Params()

	if defaultParams == newParams {
		t.Fatalf("Expected different params, not: %s", newParams)
	}

	err := SetConfig(conf)
	if err != nil {
		t.Fatal(err)
	}

	roundTrip(t, plaintext)
}

func setSalt(salt string) {
	SaltMine = func(n int) ([]byte, error) { return []byte(salt), nil }
}

func TestKey(t *testing.T) {
	for i, v := range good {

		err := setConfig(len(v.output), len(v.salt), v.N, v.r, v.p)
		if err != nil {
			t.Errorf("%d: unexpected error setting config: %s", i, err)
		}

		setSalt(v.salt)

		encoded, err := mcf.Create(v.password)
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}

		passwd := password.New([]byte("scrypt"))
		err = passwd.Parse([]byte(encoded))
		if err != nil {
			t.Errorf("%d: unexpected error creating password instance: %s", err)
		}

		if !bytes.Equal(passwd.Salt, []byte(v.salt)) {
			t.Errorf("%d: salt: expected %s, got %s", i, v.salt, string(passwd.Salt))
		}

		if !bytes.Equal(passwd.Key, v.output) {
			t.Errorf("%d: expected %x, got %x", i, v.output, passwd.Key)
		}
	}

	for i, v := range bad {
		err := setConfig(32, len(v.salt), v.N, v.r, v.p)
		if err == nil {
			t.Errorf("%d: expected error, got nil", i)
		}

		setSalt(v.salt)

		_, err = mcf.Create(v.password)
		if err == nil {
			t.Errorf("%d: expected error, got nil", i)
		}
	}
}
