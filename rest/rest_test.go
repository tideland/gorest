// Tideland Go REST Server Library - REST - Unit Tests
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/logger"
	"github.com/tideland/golib/version"

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
	req := restaudit.NewRequest(assert, "GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertBodyContains(`"ResourceID":"4711"`)
	resp.AssertBodyContains(`"Query":"0815"`)
	resp.AssertBodyContains(`"Context":"foo"`)
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
	req := restaudit.NewRequest(assert, "PUT", "/base/test/json/4711")
	reqData := TestRequestData{"foo", "bar", "4711", "0815", ""}
	req.SetBody(restaudit.ApplicationJSON, reqData)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	respData := TestRequestData{}
	resp.AssertBody(&respData)
	assert.Equal(respData, reqData)
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
	req := restaudit.NewRequest(assert, "GET", "/base/test/xml/4711")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationXML)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertBodyContains(`<ResourceID>4711</ResourceID>`)
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
	req := restaudit.NewRequest(assert, "PUT", "/base/test/xml/4711")
	reqData := TestRequestData{"foo", "bar", "4711", "0815", ""}
	req.SetBody(restaudit.ApplicationXML, reqData)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	respData := TestRequestData{}
	resp.AssertBody(&respData)
	assert.Equal(respData, reqData)
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
	req := restaudit.NewRequest(assert, "POST", "/base/test/gob")
	reqData := TestCounterData{"test", 4711}
	req.SetBody(restaudit.ApplicationGOB, reqData)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	respData := TestCounterData{}
	resp.AssertBody(&respData)
	assert.Equal(respData, reqData)
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
	req := restaudit.NewRequest(assert, "GET", "/base/content/blog/2014/09/30/just-a-test")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(`Resource ID: 2014/09/30/just-a-test`)
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
	req := restaudit.NewRequest(assert, "GET", "/base/x/y")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(`Resource: y`)
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
	req := restaudit.NewRequest(assert, "GET", "/base/test/stack")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains("Resource: token")
	resp.AssertHeaderEquals("Token", "foo")
	req = restaudit.NewRequest(assert, "GET", "/base/test/stack")
	req.AddHeader("token", "foo")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("Resource: stack")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("Resource: stack")
}

// TestVersion tests request and response version.
func TestVersion(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	req := restaudit.NewRequest(assert, "GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertHeaderEquals("Version", "1.0.0")

	req = restaudit.NewRequest(assert, "GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	req.AddHeader("Version", "2")
	resp = ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertHeaderEquals("Version", "2.0.0")

	req = restaudit.NewRequest(assert, "GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(restaudit.HeaderAccept, restaudit.ApplicationJSON)
	req.AddHeader("Version", "3.0")
	resp = ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertHeaderEquals("Version", "4.0.0-alpha")
}

//  TestDeregister tests the different possibilities to stop handlers.
func TestDeregister(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := restaudit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.RegisterAll(rest.Registrations{
		{"deregister", "single", NewTestHandler("s1", assert)},
		{"deregister", "pair", NewTestHandler("p1", assert)},
		{"deregister", "pair", NewTestHandler("p2", assert)},
		{"deregister", "group", NewTestHandler("g1", assert)},
		{"deregister", "group", NewTestHandler("g2", assert)},
		{"deregister", "group", NewTestHandler("g3", assert)},
		{"deregister", "group", NewTestHandler("g4", assert)},
		{"deregister", "group", NewTestHandler("g5", assert)},
		{"deregister", "group", NewTestHandler("g6", assert)},
	})
	assert.Nil(err)
	// Perform tests.
	assert.Equal(mux.RegisteredHandlers("deregister", "single"), []string{"s1"})
	assert.Equal(mux.RegisteredHandlers("deregister", "pair"), []string{"p1", "p2"})
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g2", "g3", "g4", "g5", "g6"})

	mux.Deregister("deregister", "single", "s1")
	assert.Nil(mux.RegisteredHandlers("deregister", "single"))
	mux.Deregister("deregister", "single")
	assert.Nil(mux.RegisteredHandlers("deregister", "single"))

	mux.Deregister("deregister", "pair")
	assert.Nil(mux.RegisteredHandlers("deregister", "pair"))

	mux.Deregister("deregister", "group", "x99")
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g2", "g3", "g4", "g5", "g6"})
	mux.Deregister("deregister", "group", "g5")
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g2", "g3", "g4", "g6"})
	mux.Deregister("deregister", "group", "g4", "g2")
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g3", "g6"})
	mux.Deregister("deregister", "group")
	assert.Nil(mux.RegisteredHandlers("deregister", "group"))
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
	req := restaudit.NewRequest(assert, "OPTION", "/base/test/method")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains("OPTION")
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
	job.EnhanceContext(func(ctx context.Context) context.Context {
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
	Query      string
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
<li>Query {{.Query}}</li>
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
	query := job.Query().ValueAsString("foo", "bar")
	precedence, _ := job.Version().Compare(version.New(3, 0, 0))
	// Create response.
	data := TestRequestData{job.Domain(), job.Resource(), job.ResourceID(), query, ctxTest.(string)}
	if precedence == version.Equal {
		job.SetVersion(version.New(4, 0, 0, "alpha"))
	}
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
