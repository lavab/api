// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package password manipulates passwords stored in Modular Crypt Format.
package password

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

//This separates password fields
const separator byte = '$'

// Passwd is a password separated into components
type Passwd struct {
	Name   []byte
	Params []byte
	Salt   []byte
	Key    []byte

	// Encoder is a function that encodes Salt and Key into serialized form
	// Defaults to EncoderBase64, but can be changed to EncoderHex or other function.
	Encoder func([]byte) []byte

	// Decoder deserializes the Salt and Key fields.
	// The default recognizes Base64 or Hex.
	// If you replace Encoder with something else, replace Decoder too, possibly with a fallback to the
	// default.
	Decoder func([]byte) ([]byte, error)
}

// count of "interesting" fields.
const fieldCount = 4

//ErrorInputPassword is returned for input passwords that fail validation.
//The struct can be examined for a possible solution.
type ErrorInputPassword struct {
	Msg      string //Error message
	Password string //Input password
}

func (e ErrorInputPassword) Error() string {
	return e.Msg
}

// New returns a Passwd struct initialized with the default encoders.
func New(name []byte) *Passwd {
	return &Passwd{Name: name, Decoder: decode, Encoder: EncodeBase64}
}

// Parse extracts an encoded password in Modular Crypt Format into a Passwd structure.
// The input is validated and should match what String() produces.
func (p *Passwd) Parse(encoded []byte) (err error) {

	inputErr := func(format string, args ...interface{}) error {
		return ErrorInputPassword{
			fmt.Sprintf("%s: "+format, append([]interface{}{string(p.Name)}, args...)...),
			string(encoded),
		}
	}

	if len(encoded) == 0 {
		return inputErr("empty password")
	}

	if encoded[0] != separator {
		return inputErr("password does not begin with separator")
	}

	parts := bytes.Split(encoded[1:], []byte{separator})
	if len(parts) < fieldCount {
		return inputErr("password has too few fields")
	}
	if len(parts) > fieldCount {
		return inputErr("password has too many fields")
	}

	if b := parts[0]; subtle.ConstantTimeCompare(b, p.Name) == 0 {
		return inputErr("unexpected password type: %s", string(b))
	}

	p.Params = parts[1]

	p.Salt, err = p.Decoder(parts[2])
	if err != nil {
		return
	}

	p.Key, err = p.Decoder(parts[3])

	return
}

// Bytes produces an encoded password in Modular Crypt Format.
// The output can be stored and later used to verify the password.
func (p *Passwd) Bytes() []byte {

	in := [][]byte{p.Name, p.Params, p.Encoder(p.Salt), p.Encoder(p.Key)}

	n := 0
	for _, b := range in {
		n += len(b)
	}

	n += len(in) //add separator for each item

	out := make([]byte, n)

	n = 0
	for _, b := range in {
		out[n] = separator
		n++
		n += copy(out[n:], b)
	}

	return out
}

// determine input type and decode accordingly
func decode(encoded []byte) (dst []byte, err error) {
	var b64NotHex = []byte("GHIJKLMNOPQRSTUVWXYZghijklmnopqrstuvwxyz+/-_")
	if bytes.IndexAny(b64NotHex, string(encoded)) < 0 {
		dst = make([]byte, hex.DecodedLen(len(encoded)))
		n, err := hex.Decode(dst, encoded)
		return dst[:n], err
	}

	enc := base64.StdEncoding
	dst = make([]byte, enc.DecodedLen(len(encoded)))
	n, err := enc.Decode(dst, encoded)
	return dst[:n], err
}

// EncodeHex encodes the input bytes into hex format.
func EncodeHex(in []byte) (out []byte) {
	out = make([]byte, hex.EncodedLen(len(in)))
	hex.Encode(out, in)
	return
}

// EncodeBase64 encodes the input bytes into standard base64 format.
func EncodeBase64(in []byte) (out []byte) {
	enc := base64.StdEncoding
	out = make([]byte, enc.EncodedLen(len(in)))
	enc.Encode(out, in)
	return out
}
