// Tideland Go REST Server Library - JSON Web Token - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt_test

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
)

//--------------------
// TESTS
//--------------------

var (
	payload = testPayload{
		Sub:   "1234567890",
		Name:  "John Doe",
		Admin: true,
	}

	hsTests = []struct {
		algorithm jwt.Algorithm
		key       string
		signature string
	}{{
		algorithm: jwt.HS256,
		key:       "secret",
		signature: "TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
	}, {
		algorithm: jwt.HS384,
		key:       "secret",
		signature: "DtVnCyiYCsCbg8gUP-579IC2GJ7P3CtFw6nfTTPw-0lZUzqgWAo9QIQElyxOpoRm",
	}, {
		algorithm: jwt.HS512,
		key:       "secret",
		signature: "YI0rUGDq5XdRw8vW2sDLRNFMN8Waol03iSFH8I4iLzuYK7FKHaQYWzPt0BJFGrAmKJ6SjY0mJIMZqNQJFVpkuw",
	}}
)

// TestHSAlgorithms tests the HMAC algorithms for the
// JWT signature.
func TestHSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	for _, test := range hsTests {
		assert.Logf("testing algorithm %q", test.algorithm)
		// Encode.
		jwtEncode, err := jwt.Encode(payload, []byte(test.key), test.algorithm)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		assert.Equal(parts[2], test.signature)
		// Verify.
		var verifyPayload testPayload
		jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, []byte(test.key))
		assert.Nil(err)
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		assert.Equal(payload.Sub, verifyPayload.Sub)
		assert.Equal(payload.Name, verifyPayload.Name)
		assert.Equal(payload.Admin, verifyPayload.Admin)
	}
}

//--------------------
// HELPERS
//--------------------

type testPayload struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

// EOF
