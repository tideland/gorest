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
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tideland/golib/audit"
)

//--------------------
// CONSTENTS
//--------------------

const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
)

//--------------------
// TEST TYPES
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

// SetJSONContent sets the content of a request to JSON.
func (r *Request) SetJSONContent(assert audit.Assertion, data interface{}) {
	body, err := json.Marshal(data)
	assert.Nil(err)
	r.Body = body
	r.Header = KeyValues{
		HeaderContentType: ApplicationJSON,
		HeaderAccept:      ApplicationJSON,
	}
}

// SetXMLContent sets the content of a request to XML.
func (r *Request) SetXMLContent(assert audit.Assertion, data interface{}) {
	body, err := xml.Marshal(data)
	assert.Nil(err)
	r.Body = body
	r.Header = KeyValues{
		HeaderContentType: ApplicationXML,
		HeaderAccept:      ApplicationXML,
	}
}

// Response wraps all infos of a test response.
type Response struct {
	Status  int
	Header  KeyValues
	Cookies KeyValues
	Body    []byte
}

// JSONContent retrieves the JSON content and unmarshals it.
func (r *Response) JSONContent(assert audit.Assertion, data interface{}) {
	contentType, ok := r.Header[HeaderContentType]
	assert.True(ok)
	assert.Equal(contentType, ApplicationJSON)
	err := json.Unmarshal(r.Body, data)
	assert.Nil(err)
}

// XMLContent retrieves the XML content and unmarshals it.
func (r *Response) XMLContent(assert audit.Assertion, data interface{}) {
	contentType, ok := r.Header[HeaderContentType]
	assert.True(ok)
	assert.Equal(contentType, ApplicationJSON)
	err := xml.Unmarshal(r.Body, data)
	assert.Nil(err)
}

//--------------------
// TEST SERVER
//--------------------

// TestServer defines the test server with methods for requests
// and uploads.
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

// Close implements the TestServer interface.
func (ts *testServer) Close() {
	ts.server.Close()
}

// DoRequest implements the TestServer interface.
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

// DoUpload implements the TestServer interface.
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
		Status:  hr.StatusCode,
		Header:  respHeader,
		Cookies: respCookies,
		Body:    respBody,
	}
}

// EOF
