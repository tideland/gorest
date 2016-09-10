// Tideland Go REST Server Library - REST Audit
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package restaudit

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tideland/golib/audit"
)

//--------------------
// TEST TOOLS
//--------------------

// KeyValues handles keys and values for request headers and cookies.
type KeyValues map[string]string

// Request wraps all infos for a test request.
type Request struct {
	Method           string
	Path             string
	Header           KeyValues
	Cookies          KeyValues
	Body             []byte
	RequestProcessor func(req *http.Request) *http.Request
}

// Response wraps all infos of a test response.
type Response struct {
	Header  KeyValues
	Cookies KeyValues
	Body    []byte
}

//--------------------
// TEST SERVER
//--------------------

type TestServer interface {
	// Close shuts down the server and blocks until all outstanding
	// requests have completed.
	Close()

	// DoRequest performs a request against the test server.
	DoRequest(req *Request) *Response

	// DoUpload is a special request for uploading a file.
	DoUpload(path, fieldname, filename, data string) *Response
}

// testServer implements the TestServer interface.
type testServer struct {
	server *httptest.Server
	assert audit.Assertion
}

// StartServer starts a test server using the passed handler
func StartServer(handler http.Handler, assert audit.Assertion) TestServer {
	return &testServer{
		server: httptest.NewServer(handler),
		assert: assert,
	}
}

// Close is specified on the TestServer interface.
func (ts *testServer) Close() {
	ts.server.Close()
}

// DoRequest is specified on the TestServer interface.
func (ts *testServer) DoRequest(req *Request) *Response {
	// First prepare it.
	transport := &http.Transport{}
	c := &http.Client{Transport: transport}
	url := ts.server.URL + req.Path
	var bodyReader io.Reader
	if req.Body != nil {
		bodyReader = ioutil.NopCloser(bytes.NewBuffer(req.Body))
	}
	httpReq, err := http.NewRequest(req.Method, url, bodyReader)
	ts.assert.Nil(err, "cannot prepare request")
	for key, value := range req.Header {
		httpReq.Header.Set(key, value)
	}
	for key, value := range req.Cookies {
		cookie := &http.Cookie{
			Name:  key,
			Value: value,
		}
		httpReq.AddCookie(cookie)
	}
	// Check if request shall be processed before performed.
	if req.RequestProcessor != nil {
		httpReq = req.RequestProcessor(httpReq)
	}
	// Now do it.
	resp, err := c.Do(httpReq)
	ts.assert.Nil(err, "cannot perform test request")
	return ts.response(resp)
}

// DoUpload is specified on the TestServer interface.
func (ts *testServer) DoUpload(path, fieldname, filename, data string) *Response {
	// Prepare request.
	transport := &http.Transport{}
	c := &http.Client{Transport: transport}
	url := ts.server.URL + path
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile(fieldname, filename)
	ts.assert.Nil(err, "cannot create form file")
	_, err = io.WriteString(part, data)
	ts.assert.Nil(err, "cannot write data")
	contentType := writer.FormDataContentType()
	err = writer.Close()
	ts.assert.Nil(err, "cannot close multipart writer")
	// And now do it.
	resp, err := c.Post(url, contentType, buffer)
	ts.assert.Nil(err, "cannot perform test upload")
	return ts.response(resp)
}

// response creates a Response instance out of the http.Response-
func (ts *testServer) response(hr *http.Response) *Response {
	respHeader := KeyValues{}
	for key, values := range hr.Header {
		respHeader[key] = strings.Join(values, ", ")
	}
	respCookies := KeyValues{}
	for _, cookie := range hr.Cookies() {
		respCookies[cookie.Name] = cookie.Value
	}
	respBody, err := ioutil.ReadAll(hr.Body)
	ts.assert.Nil(err, "cannot read response")
	defer hr.Body.Close()
	return &Response{
		Header:  respHeader,
		Cookies: respCookies,
		Body:    respBody,
	}
}

// EOF
