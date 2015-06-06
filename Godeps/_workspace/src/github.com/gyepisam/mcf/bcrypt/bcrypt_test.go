// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bcrypt

// Test vector copied from https://code.google.com/p/bcryptnet/source/browse/src/BCryptNET.Tests/BCryptTest.cs
// and edited for Go.

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/gyepisam/mcf"
	"testing"
)

var testVectors = []struct {
	plain  string
	salt   string
	passwd string
}{
	{"",
		"$2a$06$DCq7YPn5Rq63x1Lad4cll.",
		"$2a$06$DCq7YPn5Rq63x1Lad4cll.TV4S6ytwfsfvkgY8jIucDrjc8deX1s."},
	{"",
		"$2a$08$HqWuK6/Ng6sg9gQzbLrgb.",
		"$2a$08$HqWuK6/Ng6sg9gQzbLrgb.Tl.ZHfXLhvt/SgVyWhQqgqcZ7ZuUtye"},
	{"",
		"$2a$10$k1wbIrmNyFAPwPVPSVa/ze",
		"$2a$10$k1wbIrmNyFAPwPVPSVa/zecw2BCEnBwVS2GbrmgzxFUOqW9dk4TCW"},
	{"",
		"$2a$12$k42ZFHFWqBp3vWli.nIn8u",
		"$2a$12$k42ZFHFWqBp3vWli.nIn8uYyIkbvYRvodzbfbK18SSsY.CsIQPlxO"},
	{"a",
		"$2a$06$m0CrhHm10qJ3lXRY.5zDGO",
		"$2a$06$m0CrhHm10qJ3lXRY.5zDGO3rS2KdeeWLuGmsfGlMfOxih58VYVfxe"},
	{"a",
		"$2a$08$cfcvVd2aQ8CMvoMpP2EBfe",
		"$2a$08$cfcvVd2aQ8CMvoMpP2EBfeodLEkkFJ9umNEfPD18.hUF62qqlC/V."},
	{"a",
		"$2a$10$k87L/MF28Q673VKh8/cPi.",
		"$2a$10$k87L/MF28Q673VKh8/cPi.SUl7MU/rWuSiIDDFayrKk/1tBsSQu4u"},
	{"a",
		"$2a$12$8NJH3LsPrANStV6XtBakCe",
		"$2a$12$8NJH3LsPrANStV6XtBakCez0cKHXVxmvxIlcz785vxAIZrihHZpeS"},
	{"abc",
		"$2a$06$If6bvum7DFjUnE9p2uDeDu",
		"$2a$06$If6bvum7DFjUnE9p2uDeDu0YHzrHM6tf.iqN8.yx.jNN1ILEf7h0i"},
	{"abc",
		"$2a$08$Ro0CUfOqk6cXEKf3dyaM7O",
		"$2a$08$Ro0CUfOqk6cXEKf3dyaM7OhSCvnwM9s4wIX9JeLapehKK5YdLxKcm"},
	{"abc",
		"$2a$10$WvvTPHKwdBJ3uk0Z37EMR.",
		"$2a$10$WvvTPHKwdBJ3uk0Z37EMR.hLA2W6N9AEBhEgrAOljy2Ae5MtaSIUi"},
	{"abc",
		"$2a$12$EXRkfkdmXn2gzds2SSitu.",
		"$2a$12$EXRkfkdmXn2gzds2SSitu.MW9.gAVqa9eLS1//RYtYCmB1eLHg.9q"},
	{"abcdefghijklmnopqrstuvwxyz",
		"$2a$06$.rCVZVOThsIa97pEDOxvGu",
		"$2a$06$.rCVZVOThsIa97pEDOxvGuRRgzG64bvtJ0938xuqzv18d3ZpQhstC"},
	{"abcdefghijklmnopqrstuvwxyz",
		"$2a$08$aTsUwsyowQuzRrDqFflhge",
		"$2a$08$aTsUwsyowQuzRrDqFflhgekJ8d9/7Z3GV3UcgvzQW3J5zMyrTvlz."},
	{"abcdefghijklmnopqrstuvwxyz",
		"$2a$10$fVH8e28OQRj9tqiDXs1e1u",
		"$2a$10$fVH8e28OQRj9tqiDXs1e1uxpsjN0c7II7YPKXua2NAKYvM6iQk7dq"},
	{"abcdefghijklmnopqrstuvwxyz",
		"$2a$12$D4G5f18o7aMMfwasBL7Gpu",
		"$2a$12$D4G5f18o7aMMfwasBL7GpuQWuP3pkrZrOAnqP.bmezbMng.QwJ/pG"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD",
		"$2a$06$fPIsBO8qRqkjj273rfaOI.",
		"$2a$06$fPIsBO8qRqkjj273rfaOI.HtSV9jLDpTbZn782DC6/t7qT67P6FfO"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD",
		"$2a$08$Eq2r4G/76Wv39MzSX262hu",
		"$2a$08$Eq2r4G/76Wv39MzSX262huzPz612MZiYHVUJe/OcOql2jo4.9UxTW"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD",
		"$2a$10$LgfYWkbzEvQ4JakH7rOvHe",
		"$2a$10$LgfYWkbzEvQ4JakH7rOvHe0y8pHKF9OaFgwUZ2q7W2FFZmZzJYlfS"},
	{"~!@#$%^&*()      ~!@#$%^&*()PNBFRD",
		"$2a$12$WApznUOJfkEGSmYRfnkrPO",
		"$2a$12$WApznUOJfkEGSmYRfnkrPOr466oFDCaj4b6HY3EXGvfxm43seyhgC"},
}

func TestVectors(t *testing.T) {
	for i, v := range testVectors {

		x := strings.Split(v.salt[1:], "$")
		if len(x) != 3 {
			t.Fatalf("%d: invalid field structure: %s", i, v.salt)
		}
		costIn, saltIn := x[1], x[2] //index 0 is not interesting.

		salt, err := base64Decode([]byte(saltIn))

		if err != nil {
			t.Fatalf("%d: error decoding salt: %s", i, err)
		}
		var cost int
		if _, err := fmt.Sscanf(costIn, "%02d", &cost); err != nil {
			t.Fatalf("%d: error decoding param [%+v]: %s", i, costIn, err)
		}

		// Hijack rand.Reader to feed salt to crypto/bcrypt.
		// Need extra salt because SetCost generates a digest
		// to test the new cost, and each of the 3 calls to it uses
		// up one portion of salt.
		rand.Reader = bytes.NewReader(bytes.Repeat(salt, 4))

		err = SetCost(cost)
		if err != nil {
			t.Errorf("%d: SetCost: unexpected error: %s", i, err)
		}

		encoded, err := mcf.Create(v.plain)
		if err != nil {
			// password must be at least 2 bytes otherwise crypto/blowfish complains.
			// This really should be handled in crypto/bcrypt
			if len(v.plain) < 3 {
				continue
			}
			t.Errorf("%d: unexpected error: %s", i, err)
			continue
		}

		if want, got := v.passwd, encoded; want != got {
			t.Errorf("%d: output mismatch. want: %s, got %s", i, want, got)
			continue
		}

		isValid, err := mcf.Verify(v.plain, encoded)
		if err != nil {
			t.Errorf("%d: verify: unexpected failure on %q: %s", i, encoded, err)
			continue
		}
		if !isValid {
			t.Errorf("%d: IsValid: expecting true got false", i)
			continue
		}

		for j, pair := range []struct {
			cost   int
			answer bool
		}{{cost, true}, {cost + 1, false}} {

			err := SetCost(pair.cost)
			if err != nil {
				t.Errorf("%d-%d: SetCost: unexpected error: %s", i, j, err)
			}

			isCurrent, err := mcf.IsCurrent(encoded)
			if err != nil {
				t.Errorf("%d-%d: IsCurrent: unexpected failure: %", i, j, err)
				continue
			}
			if isCurrent != pair.answer {
				t.Errorf("%d-%d: IsCurrent: expecting %t got %t", i, j, pair.answer, isCurrent)
				continue
			}
		}
	}
}
