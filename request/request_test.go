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
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gorest/request"
	"github.com/tideland/gorest/rest"
)

//--------------------
// TESTS
//--------------------

// data is used for transfer data in the tests.
type data struct {
	Index int
	Name  string
}

// tests defines requests and asserts.
var tests = []struct {
	name     string
	method   string
	resource string
	id       string
	params   *request.Parameters
	expected *data
}{
	{
		name:     "GET for one item formatted in JSON",
		method:   "GET",
		resource: "item",
		id:       "foo",
		params: &request.Parameters{
			ContentType: rest.ContentTypeJSON,
		},
		expected: &data{
			Index: 0,
			Name:  "foo",
		},
	}, {
		name:     "GET for one item formatted in XML",
		method:   "GET",
		resource: "item",
		id:       "foo",
		params: &request.Parameters{
			ContentType: rest.ContentTypeXML,
		},
		expected: &data{
			Index: 0,
			Name:  "foo",
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
		}
		assert.Nil(err)
		assert.True(response.HasContentType(test.params.ContentType))
		var content data
		err = response.Read(&content)
		if test.expected != nil {
			assert.Equal(content.Name, test.expected.Name)
		}
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
	th.assert.Logf("CT: %v", job.Request().Header.Get("Content-Type"))
	switch {
	case job.HasContentType(rest.ContentTypeJSON):
		th.assert.Logf("CT: JSON")
		job.JSON(true).Write(rest.StatusOK, th.data(job.ResourceID()))
	case job.HasContentType(rest.ContentTypeXML):
		th.assert.Logf("CT: XML")
		job.XML().Write(rest.StatusOK, th.data(job.ResourceID()))
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

func (th *TestHandler) data(name string) *data {
	return &data{
		Index: th.index,
		Name:  name,
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
	// addresses := []string{":12345", ":12346", ":12347", ":12348", ":12349"}
	addresses := []string{":12345"}
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
	time.Sleep(5 * time.Millisecond)
	return servers
}

// EOF
