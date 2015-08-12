// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pbkdf2 implements a password encoding mechanism for the mcf framework
package pbkdf2

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"code.google.com/p/go.crypto/pbkdf2"

	"github.com/gyepisam/mcf"
	"github.com/gyepisam/mcf/bridge"
)

// Hash represents the HMAC hash function that the PBKDF2 algorithm uses as a pseudorandom function.
type Hash string

// Hash implements the Stringer interface
func (h Hash) String() string { return string(h) }

// Size is the output length of the hash function.
func (h Hash) Size() int {
	hash, ok := hashes[h]
	if !ok {
		panic("unknown hash: " + string(h))
	}
	return hash().Size()
}

// Available hashes
const (
	SHA1   Hash = "SHA1"
	SHA224 Hash = "SHA224"
	SHA256 Hash = "SHA256"
	SHA384 Hash = "SHA384"
	SHA512 Hash = "SHA512"
)

var hashes = map[Hash]func() hash.Hash{
	SHA1:   sha1.New,
	SHA224: sha256.New224,
	SHA256: sha256.New,
	SHA384: sha512.New384,
	SHA512: sha512.New,
}

// Config contains the parameters for the PBKDF2 algorithm along with associated values.
type Config struct {
	// The Pseudo Random Function (prf) used by the PBKDF2 algorithm.
	// When this value is changed, KeyLen will most likely need to change as well
	// since hash functions generally produce output of differing lengths.
	Hash Hash

	// Number of iteration rounds in the PBKDF2 algorithm.
	// The RFC recommends at least 1000
	Iterations int

	// Length of key produced by algorithm in bytes.
	// Defaults to the output length of the HMAC Hash.
	KeyLen int

	// Size of salt in bytes.
	// The RFC recommends at least 8 bytes.
	SaltLen int
}

// Default values. These are exported for documentation purposes.
// See GetConfig() and SetConfig() on how to change them.
const (
	DefaultIterations = 2000

	DefaultSaltLen = 16

	DefaultPrf = SHA1
)

// This value is exported for documentation purposes.
// See GetConfig() and SetConfig() on how to change it.
// Defined as a var to allow for non constant initialization.
var DefaultKeyLen = DefaultPrf.Size()

// Returns the default configuration.
// To change default values, pass the modified output of GetConfig() to SetConfig().
// See SetConfig() for an example.
func GetConfig() Config {
	return Config{
		Iterations: DefaultIterations,
		SaltLen:    DefaultSaltLen,
		Hash:       DefaultPrf,
		KeyLen:     DefaultKeyLen,
	}
}

// SetConfig establishes a new default configuration for this algorithm.
// It is only necessary to set the configuration if this algorithm is the default algorithm.
// To change a configuration,
// Get a copy of the configuration
//
//      config := pbkdf2.GetConfig()
//
// Change as necessary
//
//      config.Iterations *= 1.5        // increase iterations by 50%
//      config.Hash = pbkdf2.SHA256  // change hmac
//      config.KeyLen= pbkdf2.SHA256.Size()
//
// then set
//
//      err := pbkdf2.SetConfig(config)
//      // error handling elided
func SetConfig(config Config) error {
	err := (&config).validate()
	if err != nil {
		return err
	}
	return register(config)
}

// SaltMine is a custom source of salt, which is normally unset.
// Change this to override the use of rand.Reader if you need to use a custom salt producer.
var SaltMine mcf.SaltMiner = nil

func register(config Config) error {

	// Constructor for Implementer. Always return a fresh copy.
	fn := func() bridge.Implementer {
		c := config
		return &c
	}

	// the bridge handles the generic parts of the interface
	enc := bridge.New([]byte("pbkdf2"), fn)

	return mcf.Register(mcf.PBKDF2, enc)
}

func init() {
	err := register(GetConfig())
	if err != nil {
		panic(err)
	}
}

// ErrInvalidHash is returned when an invalid Hash is encountered.
// The name of the hash is printed in the Error() string and is also exported.
type ErrInvalidHash struct {
	Hash Hash
}

// ErrInvalidHash implements the Error interface.
func (e *ErrInvalidHash) Error() string {
	return fmt.Sprintf("Invalid Hash: %s", e.Hash)
}

func (c *Config) validate() error {
	if _, ok := hashes[c.Hash]; !ok {
		return &ErrInvalidHash{c.Hash}
	}
	return nil
}

// Keep these together
// Note that Sscanf on %s breaks on space and must therefore be the last item (and the only string).
const format = "keylen=%d,iterations=%d,hmac=%s"

// Params encodes algorithm parameters in a string for later use.
func (c *Config) Params() string {
	return fmt.Sprintf(format, c.KeyLen, c.Iterations, c.Hash)
}

// SetParams extracts encoded algorithm parameters from the output of Params().
func (c *Config) SetParams(params string) error {
	_, err := fmt.Sscanf(params, format, &c.KeyLen, &c.Iterations, &c.Hash)
	if err != nil {
		return err
	}
	return c.validate()
}

// Salt produces SaltLen bytes of random data.
func (c *Config) Salt() ([]byte, error) {
	return mcf.Salt(c.SaltLen, SaltMine)
}

// Key generates a PBKDF2 digest from the password, salt and iteration count,
// using the Hash as a pseudorandom function.
func (c *Config) Key(password, salt []byte) ([]byte, error) {
	return pbkdf2.Key(password, salt, c.Iterations, c.KeyLen, hashes[c.Hash]), nil
}

// AtLeast compares the parameters for an encoded password to the current configuration
// and returns true if the encoded password configuration has the same or longer configuration
// parameter values.
func (c *Config) AtLeast(current_imp bridge.Implementer) bool {
	current := current_imp.(*Config) // ok to panic if this fails.
	return !(c.Iterations < current.Iterations || c.KeyLen < current.KeyLen || c.SaltLen < current.SaltLen)
}
