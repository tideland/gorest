// Tideland Go REST Server Library - Handlers - Unit Tests
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers_test

//--------------------
// IMPORTS
//--------------------

import (
	"bufio"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

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
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	handler := func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(data))
	}
	err := mux.Register("test", "wrapper", handlers.NewWrapperHandler("wrapper", handler))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/wrapper",
	})
	assert.Equal(string(resp.Body), data)
}

// TestFileServeHandler tests the serving of files.
func TestFileServeHandler(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	data := "Been there, done that!"
	// Setup the test file.
	dir, err := ioutil.TempDir("", "golib-rest")
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
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "files", handlers.NewFileServeHandler("files", dir))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/files/foo.txt",
	})
	assert.Equal(string(resp.Body), data)
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/files/does.not.exist",
	})
	assert.Equal(string(resp.Body), "404 page not found\n")
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
	mux := rest.NewMultiplexer()
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
		status  int
		body    string
	}{
		{
			id:     "no-token",
			status: 401,
		}, {
			id: "token-no-gatekeeper",
			tokener: func() jwt.JWT {
				claims := jwt.NewClaims()
				claims.SetSubject("test")
				out, err := jwt.Encode(claims, key, jwt.HS512)
				assert.Nil(err)
				return out
			},
			status: 200,
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
			status: 401,
		},
	}
	// Run defined tests.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	for i, test := range tests {
		// Prepare one test.
		assert.Logf("JWT test #%d: %s", i, test.id)
		err := mux.Register("jwt", test.id, handlers.NewJWTAuthorizationHandler(test.id, test.config))
		assert.Nil(err)
		// Make request.
		var requestProcessor func(req *http.Request) *http.Request
		if test.tokener != nil {
			requestProcessor = func(req *http.Request) *http.Request {
				return jwt.AddTokenToRequest(req, test.tokener())
			}
		}
		resp := ts.DoRequest(&restaudit.Request{
			Method:           "GET",
			Path:             "/jwt/" + test.id + "/1234567890",
			RequestProcessor: requestProcessor,
		})
		assert.Equal(resp.Status, test.status)
	}
}

// EOF
