// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoder represents an interface that MCF password encoders must implement
package encoder

// An Encoder encodes a plaintext password into a hashed format.
type Encoder interface {
	// Id returns a set of bytes that identify this encoder's hashed passwords.
	// It should be the same as the id used in the encoded password and should not include a separator.
	Id() []byte

	// Create produces an encoded password from a plaintext password.
	// The application must store the encoded password for future use.
	Create(plaintext []byte) (encoded []byte, err error)

	// Verify returns true if the proffered plaintext password,
	// when encoded using the same parameters, matches the encoded password.
	Verify(plaintext, encoded []byte) (isValid bool, err error)

	// IsCurrent returns true if the parameters used to generate the encoded password
	// are at least as good as those the encoder would use to generate a new encoded password.
	// If IsCurrent returns false the encoding is out of date and should be regenerated,
	// the application should call mcf.Create() to produce a new encoding to replace the current one.
	IsCurrent(encoded []byte) (isCurrent bool, err error)
}
