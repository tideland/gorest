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
	
	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
)

//--------------------
// TESTS
//--------------------

var (
	hsTests = []struct{
		payload	  testPayload
		algorithm jwt.Algorithm
		key		  jwt.Key
		signature string
	}{
		payload: testPayload{
			Sub:	"1234567890",
			Name:	"John Doe",
			Admin:	true,
		},
		algorithm: jwt.HS256,
		key:		[]byte("secret"),
		signature: "TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
	}
}

// TestHSAlgorithms tests the HMAC algorithms for the
// JWT signature.
func TestHSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	for _, test := range hsTests {
		assert.Logf("testing algorithm %q", test.algorithm)
		// Encode.
		jwtEncode, err := jwt.Encode(test.payload, test.key, test.algorithm)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		assert.Equal(part[2], signature)
		// Verify.
		var payload testPayload
		jwtVerify, err := jwt.Verify(jwtEncode.String(), &payload, test.key)
		assert.Nil(err)
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		assert.Equal(payload.Sub, test.payload.Sub)
		assert.Equal(payload.Name, test.payload.Name)
		assert.Equal(payload.Admin, test.payload.Admin)
	}
}

//--------------------
// HELPERS
//--------------------

type testPayload struct {
	Sub		string	`json:"sub"`
	Name	string	`json:"name"`
	Admin 	bool	`json:"admin"`
}


// EOF