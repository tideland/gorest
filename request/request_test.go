// Tideland Go REST Server Library - Request - Unit Tests
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"net/http"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gorest/request"
	"github.com/tideland/gorest/rest"
)

//--------------------
// TESTS
//--------------------

// tests defines requests and asserts.
var tests = []struct {
	name     string
	method   string
	resource string
	id       string
	params   *request.Parameters
}{
	{
		name:     "GET for one item formatted in JSON",
		method:   "GET",
		resource: "item",
		id:       "foo",
		params:   &request.Parameters{
			ContentType: rest.ContentTypeJSON,
		},
	},
}

// TestRequests runs the different request tests.
func TestRequests(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	servers := newServers(assert)
	// Run the tests.
	for i, test := range tests {
		assert.Logf("test #%d: %s", i, test.name)
		caller, err := servers.Caller("testing")
		assert.Nil(err)
		var response *request.Response
		switch test.method {
		case "GET":
			response, err = caller.Get(test.resource, test.id, test.params)
			assert.Nil(err)
		}
		assert.Logf("response: %+v", response)
	}
}

//--------------------
// TEST HANDLER
//--------------------

type TestHandler struct {
	index  int
	assert audit.Assertion
}

func NewTestHandler(index int, assert audit.Assertion) rest.ResourceHandler {
	return &TestHandler{index, assert}
}

func (th *TestHandler) ID() string {
	return "test"
}

func (th *TestHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

func (th *TestHandler) Get(job rest.Job) (bool, error) {
	switch {
	case job.HasContentType(rest.ContentTypeJSON):
		job.JSON(true).Write(rest.StatusOK, th.item())
	case job.HasContentType(rest.ContentTypeXML):
		job.XML().Write(rest.StatusOK, th.item())
	}
	return true, nil
}

func (th *TestHandler) Head(job rest.Job) (bool, error) {
	return true, nil
}

func (th *TestHandler) Put(job rest.Job) (bool, error) {
	return true, nil
}

func (th *TestHandler) Post(job rest.Job) (bool, error) {
	return true, nil
}

func (th *TestHandler) Patch(job rest.Job) (bool, error) {
	return true, nil
}

func (th *TestHandler) Delete(job rest.Job) (bool, error) {
	return true, nil
}

func (th *TestHandler) Options(job rest.Job) (bool, error) {
	return true, nil
}

func (th *TestHandler) item() map[string]interface{} {
	return map[string]interface{}{
		"index": th.index,
		"name":  "Item",
	}
}

//--------------------
// SERVER
//--------------------

// newServers starts the server map for the requests.
func newServers(assert audit.Assertion) request.Servers {
	// Preparation.
	logger.SetLevel(logger.LevelDebug)
	cfgStr := "{etc {basepath /}{default-domain testing}{default-resource item}}"
	cfg, err := etc.ReadString(cfgStr)
	assert.Nil(err)
	addresses := []string{":12345", ":12346", ":12347", ":12348", ":12349"}
	servers := request.NewServers()
	// Start and register each server.
	for i, address := range addresses {
		mux := rest.NewMultiplexer(context.Background(), cfg)
		h := NewTestHandler(i, assert)
		err = mux.Register("testing", "item", h)
		assert.Nil(err)
		err = mux.Register("testing", "items", h)
		assert.Nil(err)
		go func() {
			http.ListenAndServe(address, mux)
		}()
		servers.Add("testing", "http://localhost"+address, nil)
	}
	return servers
}

// EOF
