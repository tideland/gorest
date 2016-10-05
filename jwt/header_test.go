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
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"

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
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "jwt", NewTestHandler("jwt", assert, nil, false))
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

// TestDecodeCachedRequest tests the decoding of a token
// in a handler including usage of the cache.
func TestDecodeCachedRequest(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing decode a request token using a cache")
	key := []byte("secret")
	claimsIn := initClaims()
	jwtIn, err := jwt.Encode(claimsIn, key, jwt.HS512)
	assert.Nil(err)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "jwt", NewTestHandler("jwt", assert, nil, true))
	assert.Nil(err)
	// Perform first test request.
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
	// Perform second test request.
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/jwt/1234567890",
		Header: restaudit.KeyValues{"Accept": "application/json"},
		RequestProcessor: func(req *http.Request) *http.Request {
			return jwt.AddTokenToRequest(req, jwtIn)
		},
	})
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
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "jwt", NewTestHandler("jwt", assert, key, false))
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

// TestVerifyCachedRequest tests the verification of a token
// in a handler including usage of the cache.
func TestVerifyCachedRequest(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing verify a request token using a cache")
	key := []byte("secret")
	claimsIn := initClaims()
	jwtIn, err := jwt.Encode(claimsIn, key, jwt.HS512)
	assert.Nil(err)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "jwt", NewTestHandler("jwt", assert, key, true))
	assert.Nil(err)
	// Perform first test request.
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
	// Perform second test request.
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/jwt/1234567890",
		Header: restaudit.KeyValues{"Accept": "application/json"},
		RequestProcessor: func(req *http.Request) *http.Request {
			return jwt.AddTokenToRequest(req, jwtIn)
		},
	})
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
	cache  jwt.Cache
}

func NewTestHandler(id string, assert audit.Assertion, key jwt.Key, useCache bool) rest.ResourceHandler {
	var cache jwt.Cache
	if useCache {
		cache = jwt.NewCache(time.Minute, time.Minute, time.Minute, 10)
	}
	return &testHandler{id, assert, key, cache}
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
	decode := func() (jwt.JWT, error) {
		if th.cache == nil {
			return jwt.DecodeFromJob(job)
		}
		return jwt.DecodeCachedFromJob(job, th.cache)
	}
	jwtOut, err := decode()
	th.assert.Nil(err)
	th.assert.True(jwtOut.IsValid(time.Minute))
	subject, ok := jwtOut.Claims().Subject()
	th.assert.True(ok)
	th.assert.Equal(subject, job.ResourceID())
	job.JSON(true).Write(rest.StatusOK, jwtOut.Claims())
	return true, nil
}

func (th *testHandler) testVerify(job rest.Job) (bool, error) {
	verify := func() (jwt.JWT, error) {
		if th.cache == nil {
			return jwt.VerifyFromJob(job, th.key)
		}
		return jwt.VerifyCachedFromJob(job, th.cache, th.key)
	}
	jwtOut, err := verify()
	th.assert.Nil(err)
	th.assert.True(jwtOut.IsValid(time.Minute))
	subject, ok := jwtOut.Claims().Subject()
	th.assert.True(ok)
	th.assert.Equal(subject, job.ResourceID())
	job.JSON(true).Write(rest.StatusOK, jwtOut.Claims())
	return true, nil
}

//--------------------
// HELPERS
//--------------------

// newMultiplexer creates a new multiplexer with a testing context
// and a testing configuration.
func newMultiplexer(assert audit.Assertion) rest.Multiplexer {
	cfgStr := "{etc {basepath /}{default-domain default}{default-resource default}}"
	cfg, err := etc.ReadString(cfgStr)
	assert.Nil(err)
	return rest.NewMultiplexer(context.Background(), cfg)
}

// EOF
