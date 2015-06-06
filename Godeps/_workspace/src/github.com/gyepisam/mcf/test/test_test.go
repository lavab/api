// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//package test exists to test the overall system interaction.
package test // cannot be mcf or encoder

import (
	"strings"
	"testing"

	//import all encoders
	"github.com/gyepisam/mcf"
	_ "github.com/gyepisam/mcf/bcrypt"
	_ "github.com/gyepisam/mcf/pbkdf2"
	_ "github.com/gyepisam/mcf/scrypt"
)

var plain = "password"

var encodings = []struct {
	id       string
	encoding mcf.Encoding
}{
	// order should not match import list.
	{"$pbkdf2$", mcf.PBKDF2},
	{"$scrypt$", mcf.SCRYPT},
	{"$2a$", mcf.BCRYPT},
}

func TestEncoderInteraction(t *testing.T) {

	var list []string

	// Test that SetDefault works
	for i, r := range encodings {
		err := mcf.SetDefault(r.encoding)
		if err != nil {
			t.Errorf("%d-1: SetDefault %s: unexpected error: %s", i, r.encoding, err)
			continue
		}

		encoded, err := mcf.Create(plain)
		if err != nil {
			t.Errorf("%d-1: Create: unexpected error: %s", i, err)
			continue
		}

		isValid, err := mcf.Verify(plain, encoded)
		if err != nil {
			t.Errorf("%d-1: Verify: unexpected error: %s", i, err)
			continue
		}

		if !isValid {
			t.Errorf("%d-1: Verify: unexpected failure on plain=%q encoded=%q, encoding=%q", i, plain, encoded, r.encoding)
			continue
		}

		if s := encoded; !strings.HasPrefix(s, r.id) {
			t.Errorf("%d-1: Create: encoding prefix mismatch: s=%s, prefix=%s", i, s, r.id)
			continue
		}

		list = append(list, encoded)
	}

	// Verification is independent of default setting
	// and should work as long as the encoder exists.
	for i, r := range encodings {
		err := mcf.SetDefault(r.encoding)
		if err != nil {
			t.Errorf("%d-2: SetDefault %s: unexpected error: %s", i, r.encoding, err)
			continue
		}

		for _, encoded := range list {
			isValid, err := mcf.Verify(plain, encoded)
			if err != nil {
				t.Errorf("%d-2: Verify: unexpected error: %s", i, err)
				continue
			}

			if !isValid {
				t.Errorf("%d-2: Verify: unexpected failure on plain=%q encoded=%q, encoding=%q", i, plain, encoded, r.encoding)
				continue
			}
		}
	}
}
