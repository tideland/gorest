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
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
	"github.com/tideland/gorest/rest"
	"github.com/tideland/gorest/restaudit"
)

//--------------------
// TESTS
//--------------------

// TestDecodeRequest tests the decoding of a token
// in a handler.
func TestDecodeRequest(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing decode a request token")
	key := []byte("secret")
	claimsIn := initClaims()
	jwtIn, err := jwt.Encode(claimsIn, key, jwt.HS512)
	assert.Nil(err)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "jwt", NewTestHandler("jwt", assert, nil))
	assert.Nil(err)
	// Perform test request.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/jwt/1234567890",
		Header: restaudit.KeyValues{"Accept": "application/json"},
		RequestProcessor: func(req *http.Request) *http.Request {
			return jwt.AddTokenToRequest(req, jwtIn)
		},
	})
	var claimsOut jwt.Claims
	err = json.Unmarshal(resp.Body, &claimsOut)
	assert.Nil(err)
	assert.Equal(claimsOut, claimsIn)
}

// TestVerifyRequest tests the verification of a token
// in a handler.
func TestVerifyRequest(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing verify a request token")
	key := []byte("secret")
	claimsIn := initClaims()
	jwtIn, err := jwt.Encode(claimsIn, key, jwt.HS512)
	assert.Nil(err)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "jwt", NewTestHandler("jwt", assert, key))
	assert.Nil(err)
	// Perform test request.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/jwt/1234567890",
		Header: restaudit.KeyValues{"Accept": "application/json"},
		RequestProcessor: func(req *http.Request) *http.Request {
			return jwt.AddTokenToRequest(req, jwtIn)
		},
	})
	var claimsOut jwt.Claims
	err = json.Unmarshal(resp.Body, &claimsOut)
	assert.Nil(err)
	assert.Equal(claimsOut, claimsIn)
}

//--------------------
// HANDLER
//--------------------

// testHandler is used in the test scenarios.
type testHandler struct {
	id     string
	assert audit.Assertion
	key    jwt.Key
}

func NewTestHandler(id string, assert audit.Assertion, key jwt.Key) rest.ResourceHandler {
	return &testHandler{id, assert, key}
}

func (th *testHandler) ID() string {
	return th.id
}

func (th *testHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

func (th *testHandler) Get(job rest.Job) (bool, error) {
	if th.key == nil {
		return th.testDecode(job)
	} else {
		return th.testVerify(job)
	}
}

func (th *testHandler) testDecode(job rest.Job) (bool, error) {
	jwtOut, err := jwt.DecodeFromJob(job)
	th.assert.Nil(err)
	th.assert.True(jwtOut.IsValid(time.Minute))
	subject, ok := jwtOut.Claims().Subject()
	th.assert.True(ok)
	th.assert.Equal(subject, job.ResourceID())
	job.JSON(true).Write(jwtOut.Claims())
	return true, nil
}

func (th *testHandler) testVerify(job rest.Job) (bool, error) {
	jwtOut, err := jwt.VerifyFromJob(job, th.key)
	th.assert.Nil(err)
	th.assert.True(jwtOut.IsValid(time.Minute))
	subject, ok := jwtOut.Claims().Subject()
	th.assert.True(ok)
	th.assert.Equal(subject, job.ResourceID())
	job.JSON(true).Write(jwtOut.Claims())
	return true, nil
}

// EOF
