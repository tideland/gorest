// Tideland Go REST Server Library - Handlers - Unit Tests
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers_test

//--------------------
// IMPORTS
//--------------------

import (
	"bufio"
	"context"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"

	"github.com/tideland/gorest/handlers"
	"github.com/tideland/gorest/jwt"
	"github.com/tideland/gorest/rest"
	"github.com/tideland/gorest/restaudit"
)

//--------------------
// TESTS
//--------------------

// TestWrapperHandler tests the usage of standard handler funcs
// wrapped to be used inside the package context.
func TestWrapperHandler(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	data := "Been there, done that!"
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	handler := func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(data))
	}
	err := mux.Register("test", "wrapper", handlers.NewWrapperHandler("wrapper", handler))
	assert.Nil(err)
	// Perform test requests.
	req := restaudit.NewRequest("GET", "/test/wrapper")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(data)
	punctuation := resp.AssertBodyGrep("[,!]")
	assert.Length(punctuation, 2)
}

// TestFileServeHandler tests the serving of files.
func TestFileServeHandler(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	data := "Been there, done that!"
	// Setup the test file.
	dir, err := ioutil.TempDir("", "gorest")
	assert.Nil(err)
	defer os.RemoveAll(dir)
	filename := filepath.Join(dir, "foo.txt")
	f, err := os.Create(filename)
	assert.Nil(err)
	_, err = f.WriteString(data)
	assert.Nil(err)
	assert.Logf("written %s", f.Name())
	err = f.Close()
	assert.Nil(err)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "files", handlers.NewFileServeHandler("files", dir))
	assert.Nil(err)
	// Perform test requests.
	req := restaudit.NewRequest("GET", "/test/files/foo.txt")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(data)
	req = restaudit.NewRequest("GET", "/test/files/does.not.exist")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("404 page not found")
}

// TestFileUploadHandler tests the uploading of files.
func TestFileUploadHandler(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	data := "Been there, done that!"
	// Setup the file upload processor.
	processor := func(job rest.Job, header *multipart.FileHeader, file multipart.File) error {
		assert.Equal(header.Filename, "test.txt")
		scanner := bufio.NewScanner(file)
		assert.True(scanner.Scan())
		text := scanner.Text()
		assert.Equal(text, data)
		return nil
	}
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "files", handlers.NewFileUploadHandler("files", processor))
	assert.Nil(err)
	// Perform test requests.
	ts.DoUpload("/test/files", "testfile", "test.txt", data)
}

// TestJWTAuthorizationHandler tests the authorization process
// using JSON Web Tokens.
func TestJWTAuthorizationHandler(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	key := []byte("secret")
	tests := []struct {
		id      string
		tokener func() jwt.JWT
		config  *handlers.JWTAuthorizationConfig
		runs    int
		status  int
		body    string
		auditf  handlers.AuditHandlerFunc
	}{
		{
			id:     "no-token",
			status: 401,
		}, {
			id: "token-decode-no-gatekeeper",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			status: 200,
			auditf: func(assert audit.Assertion, job rest.Job) (bool, error) {
				token, ok := jwt.FromContext(job.Context())
				assert.True(ok)
				assert.NotNil(token)
				subject, ok := token.Claims().Subject()
				assert.True(ok)
				assert.Equal(subject, "test")
				return true, nil
			},
		}, {
			id: "token-verify-no-gatekeeper",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Key: key,
			},
			status: 200,
		}, {
			id: "cached-token-verify-no-gatekeeper",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Cache: jwt.NewCache(time.Minute, time.Minute, time.Minute, 10),
				Key:   key,
			},
			runs:   5,
			status: 200,
		}, {
			id: "cached-token-verify-positive-gatekeeper",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Cache: jwt.NewCache(time.Minute, time.Minute, time.Minute, 10),
				Key:   key,
				Gatekeeper: func(job rest.Job, claims jwt.Claims) error {
					subject, ok := claims.Subject()
					assert.True(ok)
					assert.Equal(subject, "test")
					return nil
				},
			},
			runs:   5,
			status: 200,
		}, {
			id: "cached-token-verify-negative-gatekeeper",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Cache: jwt.NewCache(time.Minute, time.Minute, time.Minute, 10),
				Key:   key,
				Gatekeeper: func(job rest.Job, claims jwt.Claims) error {
					_, ok := claims.Subject()
					assert.True(ok)
					return errors.New("subject is test")
				},
			},
			runs:   1,
			status: 401,
		}, {
			id: "token-expired",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				claims.SetExpiration(time.Now().Add(-time.Hour))
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			status: 403,
		},
	}
	// Run defined tests.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	for i, test := range tests {
		// Prepare one test.
		assert.Logf("JWT test #%d: %s", i, test.id)
		err := mux.Register("jwt", test.id, handlers.NewJWTAuthorizationHandler(test.id, test.config))
		assert.Nil(err)
		if test.auditf != nil {
			err := mux.Register("jwt", test.id, handlers.NewAuditHandler("audit", assert, test.auditf))
			assert.Nil(err)
		}
		// Create request.
		req := restaudit.NewRequest("GET", "/jwt/"+test.id+"/1234567890")
		if test.tokener != nil {
			req.SetRequestProcessor(func(req *http.Request) *http.Request {
				return jwt.AddToRequest(req, test.tokener())
			})
		}
		// Make request(s).
		runs := 1
		if test.runs != 0 {
			runs = test.runs
		}
		for i := 0; i < runs; i++ {
			resp := ts.DoRequest(req)
			resp.AssertStatusEquals(test.status)
		}
	}
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
