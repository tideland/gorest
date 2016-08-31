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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

	esTests = []jwt.Algorithm{jwt.ES256, jwt.ES384, jwt.ES512}
	hsTests = []jwt.Algorithm{jwt.HS256, jwt.HS384, jwt.HS512}
	psTests = []jwt.Algorithm{jwt.PS256, jwt.PS384, jwt.PS512}
	rsTests = []jwt.Algorithm{jwt.RS256, jwt.RS384, jwt.RS512}
)

// TestESAlgorithms tests the ECDSA algorithms for the
// JWT signature.
func TestESAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	for _, test := range esTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(payload, privateKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		var verifyPayload testPayload
		jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, privateKey.Public())
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		assert.Equal(payload.Sub, verifyPayload.Sub)
		assert.Equal(payload.Name, verifyPayload.Name)
		assert.Equal(payload.Admin, verifyPayload.Admin)
	}
}

// TestHSAlgorithms tests the HMAC algorithms for the
// JWT signature.
func TestHSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	testHSKey := []byte("secret")
	for _, test := range hsTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(payload, testHSKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		var verifyPayload testPayload
		jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, testHSKey)
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		assert.Equal(payload.Sub, verifyPayload.Sub)
		assert.Equal(payload.Name, verifyPayload.Name)
		assert.Equal(payload.Admin, verifyPayload.Admin)
	}
}

// TestPSAlgorithms tests the RSAPSS algorithms for the
// JWT signature.
func TestPSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	for _, test := range psTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(payload, privateKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		var verifyPayload testPayload
		jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, privateKey.Public())
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		assert.Equal(payload.Sub, verifyPayload.Sub)
		assert.Equal(payload.Name, verifyPayload.Name)
		assert.Equal(payload.Admin, verifyPayload.Admin)
	}
}

// TestRSAlgorithms tests the RSA algorithms for the
// JWT signature.
func TestRSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	for _, test := range rsTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(payload, privateKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		var verifyPayload testPayload
		jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, privateKey.Public())
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		assert.Equal(payload.Sub, verifyPayload.Sub)
		assert.Equal(payload.Name, verifyPayload.Name)
		assert.Equal(payload.Admin, verifyPayload.Admin)
	}
}

// TestNoneAlgorithm tests the none algorithm for the
// JWT signature.
func TestNoneAlgorithm(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing algorithm \"none\"")
	// Encode.
	jwtEncode, err := jwt.Encode(payload, "", jwt.NONE)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	assert.Equal(parts[2], "")
	// Verify.
	var verifyPayload testPayload
	jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, "")
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
	assert.Equal(jwtEncode.String(), jwtVerify.String())
	assert.Equal(payload.Sub, verifyPayload.Sub)
	assert.Equal(payload.Name, verifyPayload.Name)
	assert.Equal(payload.Admin, verifyPayload.Admin)
}

// TestESTools tests the tools for the reading of PEM encoded
// ECDSA keys.
func TestESTools(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing \"ECDSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	privateBytes, err := x509.MarshalECPrivateKey(privateKeyIn)
	assert.Nil(err)
	privateBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := jwt.ReadECPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := jwt.ReadECPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	jwtEncode, err := jwt.Encode(payload, privateKeyOut, jwt.ES512)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	var verifyPayload testPayload
	jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, publicKeyOut)
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
	assert.Equal(jwtEncode.String(), jwtVerify.String())
	assert.Equal(payload.Sub, verifyPayload.Sub)
	assert.Equal(payload.Name, verifyPayload.Name)
	assert.Equal(payload.Admin, verifyPayload.Admin)
}

// TestRSTools tests the tools for the reading of PEM encoded
// RSA keys.
func TestRSTools(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing \"RSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKeyIn)
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := jwt.ReadRSAPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := jwt.ReadRSAPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	jwtEncode, err := jwt.Encode(payload, privateKeyOut, jwt.RS512)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	var verifyPayload testPayload
	jwtVerify, err := jwt.Verify(jwtEncode.String(), &verifyPayload, publicKeyOut)
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
	assert.Equal(jwtEncode.String(), jwtVerify.String())
	assert.Equal(payload.Sub, verifyPayload.Sub)
	assert.Equal(payload.Name, verifyPayload.Name)
	assert.Equal(payload.Admin, verifyPayload.Admin)
}

//--------------------
// HELPERS
//--------------------

// testPayload is used as payload instead of claims
// for stable mashalling.
type testPayload struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

// EOF
