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
	"context"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
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
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/test/json/4711",
		Header: restaudit.KeyValues{"Accept": "application/json"},
	})
	var data TestRequestData
	err = json.Unmarshal(resp.Body, &data)
	assert.Nil(err)
	assert.Equal(data.ResourceID, "4711")
	assert.Equal(data.Context, "foo")
}

// TestPutJSON tests the PUT command with a JSON payload and result.
func TestPutJSON(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	reqData := TestRequestData{"foo", "bar", "4711", ""}
	reqBuf, _ := json.Marshal(reqData)
	resp := ts.DoRequest(&restaudit.Request{
		Method: "PUT",
		Path:   "/base/test/json/4711",
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
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "xml", NewTestHandler("xml", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/test/xml/4711",
		Header: restaudit.KeyValues{"Accept": "application/xml"},
	})
	assert.Substring("<ResourceID>4711</ResourceID>", string(resp.Body))
}

// TestPutXML tests the PUT command with a XML payload and result.
func TestPutXML(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "xml", NewTestHandler("xml", assert))
	assert.Nil(err)
	// Perform test requests.
	reqData := TestRequestData{"foo", "bar", "4711", ""}
	reqBuf, _ := xml.Marshal(reqData)
	resp := ts.DoRequest(&restaudit.Request{
		Method: "PUT",
		Path:   "/base/test/xml/4711",
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
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "gob", NewTestHandler("putgob", assert))
	assert.Nil(err)
	// Perform test requests.
	reqData := TestCounterData{"test", 4711}
	reqBuf := new(bytes.Buffer)
	err = gob.NewEncoder(reqBuf).Encode(reqData)
	assert.Nil(err, "GOB encode.")
	resp := ts.DoRequest(&restaudit.Request{
		Method: "POST",
		Path:   "/base/test/gob",
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
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("content", "blog", NewTestHandler("default", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/content/blog/2014/09/30/just-a-test",
	})
	assert.Substring("<li>Resource ID: 2014/09/30/just-a-test</li>", string(resp.Body))
}

// TestFallbackDefault tests the fallback to default.
func TestFallbackDefault(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("testing", "index", NewTestHandler("default", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/x/y",
	})
	assert.Substring("<li>Resource: y</li>", string(resp.Body))
}

// TestHandlerStack tests a complete handler stack.
func TestHandlerStack(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.RegisterAll(rest.Registrations{
		{"authentication", "token", NewTestHandler("auth:token", assert)},
		{"test", "stack", NewAuthHandler("stack:auth", assert)},
		{"test", "stack", NewTestHandler("stack:test", assert)},
	})
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/test/stack",
	})
	token := resp.Header["Token"]
	assert.Equal(token, "foo")
	assert.Substring("<li>Resource: token</li>", string(resp.Body))
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/test/stack",
		Header: restaudit.KeyValues{"token": "foo"},
	})
	assert.Substring("<li>Resource: stack</li>", string(resp.Body))
	resp = ts.DoRequest(&restaudit.Request{
		Method: "GET",
		Path:   "/base/test/stack",
		Header: restaudit.KeyValues{"token": "foo"},
	})
	assert.Substring("<li>Resource: stack</li>", string(resp.Body))
}

// TestMethodNotSupported tests the handling of a not support HTTP method.
func TestMethodNotSupported(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "method", NewTestHandler("method", assert))
	assert.Nil(err)
	// Perform test requests.
	resp := ts.DoRequest(&restaudit.Request{
		Method: "OPTION",
		Path:   "/base/test/method",
	})
	assert.Substring("OPTION", string(resp.Body))
}

//--------------------
// AUTHENTICATION HANDLER
//--------------------

type AuthHandler struct {
	id     string
	assert audit.Assertion
}

func NewAuthHandler(id string, assert audit.Assertion) rest.ResourceHandler {
	return &AuthHandler{id, assert}
}

func (ah *AuthHandler) ID() string {
	return ah.id
}

func (ah *AuthHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

func (ah *AuthHandler) Get(job rest.Job) (bool, error) {
	token := job.Request().Header.Get("Token")
	if token != "foo" {
		job.Redirect("authentication", "token", "")
		return false, nil
	}
	job.ExtentContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "Token", "foo")
	})
	return true, nil
}

//--------------------
// TEST HANDLER
//--------------------

type TestRequestData struct {
	Domain     string
	Resource   string
	ResourceID string
	Context    string
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
<li>Context: {{.Context}}</li>
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
	env.TemplatesCache().Parse("test:context:html", testTemplateHTML, "text/html")
	return nil
}

func (th *TestHandler) Get(job rest.Job) (bool, error) {
	if th.id == "auth:token" {
		job.ResponseWriter().Header().Add("Token", "foo")
	}
	if th.id == "stack:test" {
		ctxToken := job.Context().Value("Token")
		th.assert.Equal(ctxToken, "foo")
	}
	ctxTest := job.Context().Value("test")
	data := TestRequestData{job.Domain(), job.Resource(), job.ResourceID(), ctxTest.(string)}
	switch {
	case job.AcceptsContentType(rest.ContentTypeXML):
		th.assert.Logf("GET XML")
		job.XML().Write(rest.StatusOK, data)
	case job.AcceptsContentType(rest.ContentTypeJSON):
		th.assert.Logf("GET JSON")
		job.JSON(true).Write(rest.StatusOK, data)
	default:
		th.assert.Logf("GET HTML")
		job.Renderer().Render("test:context:html", data)
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
			job.JSON(true).Write(rest.StatusBadRequest, TestErrorData{err.Error()})
		} else {
			job.JSON(true).Write(rest.StatusOK, data)
		}
	case job.HasContentType(rest.ContentTypeXML):
		err := job.XML().Read(&data)
		if err != nil {
			job.XML().Write(rest.StatusBadRequest, TestErrorData{err.Error()})
		} else {
			job.XML().Write(rest.StatusOK, data)
		}
	}

	return true, nil
}

func (th *TestHandler) Post(job rest.Job) (bool, error) {
	var data TestCounterData
	err := job.GOB().Read(&data)
	if err != nil {
		job.GOB().Write(rest.StatusBadRequest, err)
	} else {
		job.GOB().Write(rest.StatusOK, data)
	}
	return true, nil
}

func (th *TestHandler) Delete(job rest.Job) (bool, error) {
	return false, nil
}

//--------------------
// HELPERS
//--------------------

// newMultiplexer creates a new multiplexer with a testing context
// and a testing configuration.
func newMultiplexer(assert audit.Assertion) rest.Multiplexer {
	ctx := context.WithValue(context.Background(), "test", "foo")
	cfgStr := "{etc {basepath /base/}{default-domain testing}{default-resource index}}"
	cfg, err := etc.ReadString(cfgStr)
	assert.Nil(err)
	return rest.NewMultiplexer(ctx, cfg)
}

// EOF
