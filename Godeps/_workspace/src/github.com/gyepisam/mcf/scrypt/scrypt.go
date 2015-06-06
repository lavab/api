// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scrypt implements a password encoding mechanism for the mcf framework
package scrypt

import (
	"code.google.com/p/go.crypto/scrypt"

	"fmt"

	"github.com/gyepisam/mcf"
	"github.com/gyepisam/mcf/bridge"
)

// Circa 2014 work factors.
// These are exported to show default values.
// See GetConfig and SetConfig(...) to change them.
const (
	DefaultKeyLen  = 32
	DefaultSaltLen = 16
	DefaultN       = 1 << 16
	DefaultR       = 10
	DefaultP       = 2
)

// Config contains the scrypt algorithm parameters and other associated values.
// Use the GetConfig() and SetConfig() combination to change any desired parameters.
type Config struct {
	KeyLen  int //Key output size in bytes.
	SaltLen int // Length of salt in bytes.

	N int // CPU/Memory cost. Must be a power of two.
	R int // block size parameter.
	P int // parallelization parameter.
}

// Custom source of salt, normally unset.
// Set this if you need to override the user of rand.Reader and
// use a custom salt producer.
// Also useful for testing.
var SaltMine mcf.SaltMiner = nil

// ErrInvalidParameter is returned by SetConfig if any of the provided parameters
// fail validation. The error message contains the name and value of the faulty
// parameter to aid in resolving the problem.
type ErrInvalidParameter struct {
	Name  string
	Value int
}

func (e ErrInvalidParameter) Error() string {
	return fmt.Sprintf("parameter %s has invalid value: %d", e.Name, e.Value)
}

// Config returns the default configuration used to create new scrypt password hashes.
// The return value can be modified and used as a parameter to SetConfig
func GetConfig() Config {
	return Config{
		KeyLen:  DefaultKeyLen,
		SaltLen: DefaultSaltLen,
		N:       DefaultN,
		R:       DefaultR,
		P:       DefaultP,
	}
}

/*
SetConfig sets the default encoding parameters, salt length or key length.
It is best to modify a copy of the default configuration unless all parameters are changed.

Here is an example that doubles the default work factor.

	config := scrypt.GetConfig()
	config.N *= 2
	scrypt.SetConfig(config)

*/
func SetConfig(config Config) error {
	c := &config
	err := c.validate()
	if err != nil {
		return err
	}

	return register(config)
}

func register(config Config) error {
	// Constructor function. Provide fresh copy each time.
	fn := func() bridge.Implementer {
		c := config
		return &c
	}

	enc := bridge.New([]byte("scrypt"), fn)

	return mcf.Register(mcf.SCRYPT, enc)
}

func init() {
	err := register(GetConfig())
	if err != nil {
		panic(err)
	}
}

func (c *Config) validate() error {
	//punt, cheat and see if the underlying algorithm complains!
	_, err := c.Key([]byte("password"), []byte("salt"))
	return err
}

// Keep these together.
var format = "KeyLen=%d,N=%d,R=%d,P=%d"

// Params returns the current digest algorithm parameters.
func (c *Config) Params() string {
	return fmt.Sprintf(format, c.KeyLen, c.N, c.R, c.P)
}

// SetParams sets the parameters for the digest algorithm.
func (c *Config) SetParams(s string) error {
	_, err := fmt.Sscanf(s, format, &c.KeyLen, &c.N, &c.R, &c.P)
	if err != nil {
		return err
	}
	return c.validate()
}

// Salt produces SaltLen bytes of random data.
func (c *Config) Salt() ([]byte, error) {
	return mcf.Salt(c.SaltLen, SaltMine)
}

// Key returns an scrypt digest of password and salt using the algorithm parameters: N, r, and p.
// The returned value is of length KeyLen.
func (c *Config) Key(plaintext []byte, salt []byte) (b []byte, err error) {
	return scrypt.Key(plaintext, salt, c.N, c.R, c.P, c.KeyLen)
}

// AtLeast returns true if the parameters used to generate the encoded password
// are at least as good as those currently in use.
func (c *Config) AtLeast(current_imp bridge.Implementer) bool {
	current := current_imp.(*Config) // ok to panic
	return !(c.N < current.N || c.R < current.R || c.P < current.P || c.KeyLen < current.KeyLen)
}
