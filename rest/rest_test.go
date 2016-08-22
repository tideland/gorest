// Tideland Go REST Server Library - REST - Unit Tests
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest_test

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gorest/rest"
	"github.com/tideland/gorest/restaudit"
)

//--------------------
// INIT
//--------------------

func init() {
	logger.SetLevel(logger.LevelDebug)
}

//--------------------
// TESTS
//--------------------

// TestGetJSON tests the GET command with a JSON result.
func TestGetJSON(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/json/4711",
		Header: restaudit.KeyValues{"Accept": "application/json"},
	})
	var data TestRequestData
	err = json.Unmarshal(resp.Body, &data)
	assert.Nil(err)
	assert.Equal(data.ResourceID, "4711")
}

// TestPutJSON tests the PUT command with a JSON payload and result.
func TestPutJSON(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	reqData := TestRequestData{"foo", "bar", "4711"}
	reqBuf, _ := json.Marshal(reqData)
	resp := ts.DoRequest(&restaudit.Request{
		Method: "PUT",
		Path:   "/test/json/4711",
		Header: restaudit.KeyValues{"Content-Type": "application/json", "Accept": "application/json"},
		Body:   reqBuf,
	})
	var recvData TestRequestData
	err = json.Unmarshal(resp.Body, &recvData)
	assert.Nil(err)
	assert.Equal(recvData, reqData)
}

// TestGetXML tests the GET command with an XML result.
func TestGetXML(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "xml", NewTestHandler("xml", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/xml/4711",
		Header: restaudit.KeyValues{"Accept": "application/xml"},
	})
	assert.Substring("<ResourceID>4711</ResourceID>", string(resp.Body))
}

// TestPutXML tests the PUT command with a XML payload and result.
func TestPutXML(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "xml", NewTestHandler("xml", assert))
	assert.Nil(err)
	// Perform test requests.
	reqData := TestRequestData{"foo", "bar", "4711"}
	reqBuf, _ := xml.Marshal(reqData)
	resp := ts.DoRequest(&restaudit.Request{
		Method: "PUT",
		Path:   "/test/xml/4711",
		Header: restaudit.KeyValues{"Content-Type": "application/xml", "Accept": "application/xml"},
		Body:   reqBuf,
	})
	var recvData TestRequestData
	err = xml.Unmarshal(resp.Body, &recvData)
	assert.Nil(err)
	assert.Equal(recvData, reqData)
}

// TestPutGOB tests the PUT command with a GOB payload and result.
func TestPutGOB(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "gob", NewTestHandler("putgob", assert))
	assert.Nil(err)
	// Perform test requests.
	reqData := TestCounterData{"test", 4711}
	reqBuf := new(bytes.Buffer)
	err = gob.NewEncoder(reqBuf).Encode(reqData)
	assert.Nil(err, "GOB encode.")
	assert.Logf("%q", reqBuf.String())
	resp := ts.DoRequest(&restaudit.Request{
		Method: "POST",
		Path:   "/test/gob",
		Header: restaudit.KeyValues{"Content-Type": "application/vnd.tideland.gob"},
		Body:   reqBuf.Bytes(),
	})
	var respData TestCounterData
	err = gob.NewDecoder(bytes.NewBuffer(resp.Body)).Decode(&respData)
	assert.Nil(err)
	assert.Equal(respData.ID, "test")
	assert.Equal(respData.Count, int64(4711))
}

// TestLongPath tests the setting of long path tail as resource ID.
func TestLongPath(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("content", "blog", NewTestHandler("default", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/content/blog/2014/09/30/just-a-test",
	})
	assert.Substring("<li>Resource ID: 2014/09/30/just-a-test</li>", string(resp.Body))
}

// TestFallbackDefault tests the fallback to default.
func TestFallbackDefault(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("default", "default", NewTestHandler("default", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/x/y",
	})
	assert.Substring("<li>Resource: y</li>", string(resp.Body))
}

// TestHandlerStack tests a complete handler stack.
func TestHandlerStack(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.RegisterAll(rest.Registrations{
		{"authentication", "login", NewTestHandler("login", assert)},
		{"test", "stack", NewAuthHandler("foo", assert)},
		{"test", "stack", NewTestHandler("stack", assert)},
	})
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/stack",
	})
	assert.Substring("<li>Resource: login</li>", string(resp.Body))
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/stack",
		Header: restaudit.KeyValues{"password": "foo"},
	})
	assert.Substring("<li>Resource: stack</li>", string(resp.Body))
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/test/stack",
		Header: restaudit.KeyValues{"password": "foo"},
	})
	assert.Substring("<li>Resource: stack</li>", string(resp.Body))
}

// TestMethodNotSupported tests the handling of a not support HTTP method.
func TestMethodNotSupported(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := rest.NewMultiplexer()
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "method", NewTestHandler("method", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "OPTION",
		Path:   "/test/method",
	})
	assert.Substring("OPTION", string(resp.Body))
}

//--------------------
// AUTHENTICATION HANDLER
//--------------------

type AuthHandler struct {
	password string
	assert   audit.Assertion
}

func NewAuthHandler(password string, assert audit.Assertion) rest.ResourceHandler {
	return &AuthHandler{password, assert}
}

func (ah *AuthHandler) ID() string {
	return ah.password
}

func (ah *AuthHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

func (ah *AuthHandler) Get(job rest.Job) (bool, error) {
	return true, nil
}

//--------------------
// TEST HANDLER
//--------------------

type TestRequestData struct {
	Domain     string
	Resource   string
	ResourceID string
}

type TestCounterData struct {
	ID    string
	Count int64
}

type TestErrorData struct {
	Error string
}

const testTemplateHTML = `
<?DOCTYPE html?>
<html>
<head><title>Test</title></head>
<body>
<ul>
<li>Domain: {{.Domain}}</li>
<li>Resource: {{.Resource}}</li>
<li>Resource ID: {{.ResourceID}}</li>
</ul>
</body>
</html>
`

type TestHandler struct {
	id     string
	assert audit.Assertion
}

func NewTestHandler(id string, assert audit.Assertion) rest.ResourceHandler {
	return &TestHandler{id, assert}
}

func (th *TestHandler) ID() string {
	return th.id
}

func (th *TestHandler) Init(env rest.Environment, domain, resource string) error {
	env.Templates().Parse("test:context:html", testTemplateHTML, "text/html")
	return nil
}

func (th *TestHandler) Get(job rest.Job) (bool, error) {
	data := TestRequestData{job.Domain(), job.Resource(), job.ResourceID()}
	switch {
	case job.AcceptsContentType(rest.ContentTypeXML):
		th.assert.Logf("GET XML")
		job.XML().Write(data)
	case job.AcceptsContentType(rest.ContentTypeJSON):
		th.assert.Logf("GET JSON")
		job.JSON(true).Write(data)
	default:
		th.assert.Logf("GET HTML")
		job.RenderTemplate("test:context:html", data)
	}
	return true, nil
}

func (th *TestHandler) Head(job rest.Job) (bool, error) {
	return false, nil
}

func (th *TestHandler) Put(job rest.Job) (bool, error) {
	var data TestRequestData
	switch {
	case job.HasContentType(rest.ContentTypeJSON):
		err := job.JSON(true).Read(&data)
		if err != nil {
			job.JSON(true).Write(TestErrorData{err.Error()})
		} else {
			job.JSON(true).Write(data)
		}
	case job.HasContentType(rest.ContentTypeXML):
		err := job.XML().Read(&data)
		if err != nil {
			job.XML().Write(TestErrorData{err.Error()})
		} else {
			job.XML().Write(data)
		}
	}

	return true, nil
}

func (th *TestHandler) Post(job rest.Job) (bool, error) {
	var data TestCounterData
	err := job.GOB().Read(&data)
	if err != nil {
		job.GOB().Write(err)
	} else {
		job.GOB().Write(data)
	}
	return true, nil
}

func (th *TestHandler) Delete(job rest.Job) (bool, error) {
	return false, nil
}

// EOF
