// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bridge simplifies the creation of a password encoder by abstracting out the generic parts.

// Concrete implementations that use this package can be greatly simplified.
// See github.com/gyepisam/mcf/scrypt and github.com/gyepisam/mcf/pbkdf2

package bridge

import (
	"crypto/subtle"

	"github.com/gyepisam/mcf/encoder"
	"github.com/gyepisam/mcf/password"
)

// Implementer represents a concrete implementation such as scrypt or pbkdf2
type Implementer interface {
	// Params encodes the parameters used by Key.
	Params() string

	// Setparams restores parameters from string produced by Params.
	SetParams(string) error

	// Produces required amount of salt.
	Salt() ([]byte, error)

	// Key generates a key (hash, digest, key) using password, salt and implementer specific parameters.
	Key(password, salt []byte) ([]byte, error)

	// AtLeast compares the implementer to the current implementer and determines whether the
	// implementer parameters are at least as strong as the current. It is used to determine
	// whether an encoded password needs to be re-encoded.
	AtLeast(Implementer) bool
}

// Encoder implements the encoder.Encoder interface using an Implementer to
// abstract implementation specific parts.
type Encoder struct {
	name        []byte
	implementer func() Implementer
}

// New takes an implementation name and a function that produces an Implementer
// and produces an encoder.Encoder that uses the Implementer.
func New(name []byte, fn func() Implementer) encoder.Encoder {
	return &Encoder{name: name, implementer: fn}
}

// Id returns the name of the encoder, which is the type of passwords it can handle.
func (enc *Encoder) Id() []byte { return enc.name }

// Create produces an encoded password from a plaintext password using the current configuration.
// The application must store the encoded password for future use.
func (enc *Encoder) Create(plaintext []byte) (encoded []byte, err error) {

	imp := enc.implementer()

	passwd := password.New(enc.name)
	passwd.Params = []byte(imp.Params())

	passwd.Salt, err = imp.Salt()
	if err != nil {
		return
	}

	passwd.Key, err = imp.Key(plaintext, passwd.Salt)
	if err != nil {
		return
	}

	return passwd.Bytes(), nil
}

// Verify returns true if the proffered plaintext password,
// when encoded using the same parameters, matches the encoded password.
func (enc *Encoder) Verify(plaintext, encoded []byte) (isValid bool, err error) {

	passwd := password.New(enc.name)

	err = passwd.Parse(encoded)
	if err != nil {
		return
	}

	imp := enc.implementer()
	err = imp.SetParams(string(passwd.Params))
	if err != nil {
		return
	}

	testKey, err := imp.Key(plaintext, passwd.Salt)
	if err != nil {
		return
	}

	return subtle.ConstantTimeCompare(passwd.Key, testKey) == 1, nil
}

// IsCurrent returns true if the parameters used to generate the encoded password
// are at least as good as those in params.
// If IsCurrent returns false the encoding is out of date and should be regenerated,
// the application should call mcf.Create() to produce a new encoding to replace the current one.
func (enc *Encoder) IsCurrent(encoded []byte) (isCurrent bool, err error) {

	passwd := password.New(enc.name)

	err = passwd.Parse(encoded)
	if err != nil {
		return
	}

	imp := enc.implementer()
	err = imp.SetParams(string(passwd.Params))
	if err != nil {
		return
	}

	return imp.AtLeast(enc.implementer()), nil
}
