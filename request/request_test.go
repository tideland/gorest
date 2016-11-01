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
	"fmt"
	"net/http"
	"net/url"
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
	check    func(assert audit.Assertion, response *request.Response)
}{
	{
		name:     "GET for one item formatted in JSON",
		method:   "GET",
		resource: "item",
		id:       "foo",
		params: &request.Parameters{
			Accept: rest.ContentTypeJSON,
		},
		check: func(assert audit.Assertion, response *request.Response) {
			assert.True(response.HasContentType(rest.ContentTypeJSON))
			var content data
			err := response.Read(&content)
			assert.Nil(err)
			assert.Equal(content.Name, "foo")
		},
	}, {
		name:     "GET for one item formatted in XML",
		method:   "GET",
		resource: "item",
		id:       "foo",
		params: &request.Parameters{
			Accept: rest.ContentTypeXML,
		},
		check: func(assert audit.Assertion, response *request.Response) {
			assert.True(response.HasContentType(rest.ContentTypeXML))
			var content data
			err := response.Read(&content)
			assert.Nil(err)
			assert.Equal(content.Name, "foo")
		},
	}, {
		name:     "GET for one item formatted URL encoded",
		method:   "GET",
		resource: "item",
		id:       "foo",
		params: &request.Parameters{
			Accept: rest.ContentTypeURLEncoded,
		},
		check: func(assert audit.Assertion, response *request.Response) {
			assert.True(response.HasContentType(rest.ContentTypeURLEncoded))
			values := url.Values{}
			err := response.Read(values)
			assert.Nil(err)
			assert.Equal(values["name"][0], "foo")
		},
	}, {
		name:     "HEAD returning the resource ID as header",
		method:   "HEAD",
		resource: "item",
		id:       "foo",
		check: func(assert audit.Assertion, response *request.Response) {
			assert.Equal(response.Header["Resource-Id"], "foo")
		},
	},
}

// TestRequests runs the different request tests.
func TestRequests(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	servers := newServers(assert, 12345, 12346)
	// Run the tests.
	for i, test := range tests {
		assert.Logf("test #%d: %s", i, test.name)
		caller, err := servers.Caller("testing")
		assert.Nil(err)
		var response *request.Response
		switch test.method {
		case "GET":
			response, err = caller.Get(test.resource, test.id, test.params)
		case "HEAD":
			response, err = caller.Head(test.resource, test.id, test.params)
		case "PUT":
			response, err = caller.Put(test.resource, test.id, test.params)
		case "POST":
			response, err = caller.Post(test.resource, test.id, test.params)
		case "PATCH":
			response, err = caller.Patch(test.resource, test.id, test.params)
		case "DELETE":
			response, err = caller.Delete(test.resource, test.id, test.params)
		case "OPTIONS":
			response, err = caller.Options(test.resource, test.id, test.params)
		}
		assert.Nil(err)
		test.check(assert, response)
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
	th.assert.Logf("HANDLER #%d: GET", th.index)
	switch {
	case job.AcceptsContentType(rest.ContentTypeJSON):
		job.JSON(true).Write(rest.StatusOK, th.data(job.ResourceID()))
	case job.AcceptsContentType(rest.ContentTypeXML):
		job.XML().Write(rest.StatusOK, th.data(job.ResourceID()))
	case job.AcceptsContentType(rest.ContentTypeURLEncoded):
		values := url.Values{}
		values.Set("name", job.ResourceID())
		values.Set("index", fmt.Sprintf("%d", th.index))
		job.ResponseWriter().Header().Set("Content-Type", rest.ContentTypeURLEncoded)
		job.ResponseWriter().Write([]byte(values.Encode()))
	}
	return true, nil
}

func (th *TestHandler) Head(job rest.Job) (bool, error) {
	th.assert.Logf("HANDLER #%d: HEAD", th.index)
	job.ResponseWriter().Header().Set("Resource-Id", job.ResourceID())
	job.ResponseWriter().WriteHeader(rest.StatusOK)
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
func newServers(assert audit.Assertion, ports ...int) request.Servers {
	// Preparation.
	logger.SetLevel(logger.LevelDebug)
	cfgStr := "{etc {basepath /}{default-domain testing}{default-resource item}}"
	cfg, err := etc.ReadString(cfgStr)
	assert.Nil(err)
	servers := request.NewServers()
	// Start and register each server.
	for i, port := range ports {
		mux := rest.NewMultiplexer(context.Background(), cfg)
		h := NewTestHandler(i, assert)
		err = mux.Register("testing", "item", h)
		assert.Nil(err)
		err = mux.Register("testing", "items", h)
		assert.Nil(err)
		address := fmt.Sprintf(":%d", port)
		go func() {
			http.ListenAndServe(address, mux)
		}()
		servers.Add("testing", "http://localhost"+address, nil)
	}
	time.Sleep(5 * time.Millisecond)
	return servers
}

// EOF
