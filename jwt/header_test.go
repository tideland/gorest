// Tideland Go REST Server Library - JSON Web Token - Unit Tests
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
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

// TestDecodeInvalidRequest tests the decoding of requests
// without a header or an invalid one.
func TestDecodeInvalidRequest(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing decode invalid requests")
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	asserter := newHeaderAsserter(assert, ".* request contains no authorization header")
	err := mux.Register("test", "jwt", newTestHandler("jwt", asserter))
	assert.Nil(err)
	// Perform request without authorization.
	req := restaudit.NewRequest("GET", "/test/jwt/1234567890")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	resp := ts.DoRequest(req)
	ok := ""
	resp.AssertUnmarshalledBody(&ok)
	assert.Equal(ok, "OK")
}

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
	asserter := newDecodeAsserter(assert, false)
	err = mux.Register("test", "jwt", newTestHandler("jwt", asserter))
	assert.Nil(err)
	// Perform test request.
	req := restaudit.NewRequest("GET", "/test/jwt/1234567890")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	req.SetRequestProcessor(func(req *http.Request) *http.Request {
		return jwt.AddToRequest(req, jwtIn)
	})
	resp := ts.DoRequest(req)
	claimsOut := jwt.Claims{}
	resp.AssertUnmarshalledBody(&claimsOut)
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
	asserter := newDecodeAsserter(assert, true)
	err = mux.Register("test", "jwt", newTestHandler("jwt", asserter))
	assert.Nil(err)
	// Perform first test request.
	req := restaudit.NewRequest("GET", "/test/jwt/1234567890")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	req.SetRequestProcessor(func(req *http.Request) *http.Request {
		return jwt.AddToRequest(req, jwtIn)
	})
	resp := ts.DoRequest(req)
	claimsOut := jwt.Claims{}
	resp.AssertUnmarshalledBody(&claimsOut)
	assert.Equal(claimsOut, claimsIn)
	// Perform second test request.
	resp = ts.DoRequest(req)
	claimsOut = jwt.Claims{}
	resp.AssertUnmarshalledBody(&claimsOut)
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
	asserter := newVerifyAsserter(assert, key, false)
	err = mux.Register("test", "jwt", newTestHandler("jwt", asserter))
	assert.Nil(err)
	// Perform test request.
	req := restaudit.NewRequest("GET", "/test/jwt/1234567890")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	req.SetRequestProcessor(func(req *http.Request) *http.Request {
		return jwt.AddToRequest(req, jwtIn)
	})
	resp := ts.DoRequest(req)
	claimsOut := jwt.Claims{}
	resp.AssertUnmarshalledBody(&claimsOut)
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
	asserter := newVerifyAsserter(assert, key, true)
	err = mux.Register("test", "jwt", newTestHandler("jwt", asserter))
	assert.Nil(err)
	// Perform first test request.
	req := restaudit.NewRequest("GET", "/test/jwt/1234567890")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	req.SetRequestProcessor(func(req *http.Request) *http.Request {
		return jwt.AddToRequest(req, jwtIn)
	})
	resp := ts.DoRequest(req)
	claimsOut := jwt.Claims{}
	resp.AssertUnmarshalledBody(&claimsOut)
	assert.Equal(claimsOut, claimsIn)
	// Perform second test request.
	resp = ts.DoRequest(req)
	resp.AssertUnmarshalledBody(&claimsOut)
	assert.Equal(claimsOut, claimsIn)
}

//--------------------
// HANDLER
//--------------------

// testAsserter instances will handle the assertions in the testHandler.
type testAsserter func(job rest.Job) (bool, error)

func newHeaderAsserter(assert audit.Assertion, pattern string) testAsserter {
	return func(job rest.Job) (bool, error) {
		token, err := jwt.DecodeFromJob(job)
		assert.Nil(token)
		assert.ErrorMatch(err, pattern)
		job.JSON(true).Write(rest.StatusOK, "OK")
		return true, nil
	}
}

func newDecodeAsserter(assert audit.Assertion, cached bool) testAsserter {
	var cache jwt.Cache
	if cached {
		cache = jwt.NewCache(time.Minute, time.Minute, time.Minute, 10)
	}
	return func(job rest.Job) (bool, error) {
		var token jwt.JWT
		var err error
		if cached {
			token, err = jwt.DecodeCachedFromJob(job, cache)
		} else {
			token, err = jwt.DecodeFromJob(job)
		}
		assert.Nil(err)
		assert.True(token.IsValid(time.Minute))
		subject, ok := token.Claims().Subject()
		assert.True(ok)
		assert.Equal(subject, job.ResourceID())
		job.JSON(true).Write(rest.StatusOK, token.Claims())
		return true, nil
	}
}

func newVerifyAsserter(assert audit.Assertion, key jwt.Key, cached bool) testAsserter {
	var cache jwt.Cache
	if cached {
		cache = jwt.NewCache(time.Minute, time.Minute, time.Minute, 10)
	}
	return func(job rest.Job) (bool, error) {
		var token jwt.JWT
		var err error
		if cached {
			token, err = jwt.VerifyCachedFromJob(job, cache, key)
		} else {
			token, err = jwt.VerifyFromJob(job, key)
		}
		assert.Nil(err)
		assert.True(token.IsValid(time.Minute))
		subject, ok := token.Claims().Subject()
		assert.True(ok)
		assert.Equal(subject, job.ResourceID())
		job.JSON(true).Write(rest.StatusOK, token.Claims())
		return true, nil
	}
}

// testHandler is used in the test scenarios.
type testHandler struct {
	id       string
	asserter testAsserter
}

func newTestHandler(id string, asserter testAsserter) rest.ResourceHandler {
	return &testHandler{
		id:       id,
		asserter: asserter,
	}
}

func (th *testHandler) ID() string {
	return th.id
}

func (th *testHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

func (th *testHandler) Get(job rest.Job) (bool, error) {
	return th.asserter(job)
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
