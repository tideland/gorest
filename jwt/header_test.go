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
	assert.Logf("testing decoding a token")
	key := []byte("secret")
	claimsIn := initClaims()
	jwt, err := jwt.Encode(claimsIn, key, jwt.HS512)
	assert.Nil(err)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "jwt", NewTestHandler("jwt", assert))
	assert.Nil(err)
	// Perform test request.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/jwt/1234567890",
		Header: restaudit.KeyValues{"Accept": "application/json"},
		ProcessRequest: func(req *http.Request) *http.Request {
			return jwt.AddTokenToRequest(req, jwt)
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
	key	   jwt.Key
}

func NewTestHandler(id string, assert audit.Assertion, key jwt.Key) rest.ResourceHandler {
	return &TestHandler{id, assert, key}
}

func (th *TestHandler) ID() string {
	return th.id
}

func (th *TestHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

func (th *TestHandler) Get(job rest.Job) (bool, error) {
	if th.Key == nil {
		return th.testDecode(job)
	} else {
		return th.testVerify(job)
	}
}

func (th *TestHandler) testDecode(job rest.Job) (bool, error) {
	jwt, err := jwt.DecodeTokenFromJob(job)
	th.assert.Nil(err)
	th.assert.True(jwt.IsValid(time.Minute))
	subject, ok := jwt.Subject()
	th.assert.True(ok)
	th.assert.Equal(subject, job.ResourceID())
	job.JSON(true).Write(jwt.Claims())
	return true, nil
}

func (th *TestHandler) testVerify(job rest.Job) (bool, error) {
	jwt, err := jwt.VerifyTokenFromJob(job, th.key)
	th.assert.Nil(err)
	th.assert.True(jwt.IsValid(time.Minute))
	subject, ok := jwt.Subject()
	th.assert.True(ok)
	th.assert.Equal(subject, job.ResourceID())
	job.JSON(true).Write(jwt.Claims())
	return true, nil
}

// EOF