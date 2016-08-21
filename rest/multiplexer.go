// Tideland Go REST Server Library - REST - Multiplexer
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/tideland/golib/logger"
	"github.com/tideland/golib/monitoring"
)

//--------------------
// REGISTRATIONS
//--------------------

// Registration encapsulates one handler registration.
type Registration struct {
	Domain   string
	Resource string
	Handler  ResourceHandler
}

// Registrations is a number handler registratons.
type Registrations []Registration

//--------------------
// MULTIPLEXER
//--------------------

// Multiplexer enhances the http.Handler interface by registration
// an deregistration of handlers.
type Multiplexer interface {
	http.Handler

	// Register adds a resource handler for a given domain and resource.
	Register(domain, resource string, handler ResourceHandler) error

	// RegisterAll allows to register multiple handler in one run.
	RegisterAll(registrations Registrations) error

	// Deregister removes a resource handler for a given domain and resource.
	Deregister(domain, resource, id string)
}

// multiplexer implements the Multiplexer interface.
type multiplexer struct {
	mutex       sync.RWMutex
	environment Environment
	mapping     *mapping
}

// NewMultiplexer creates a new HTTP multiplexer.
func NewMultiplexer(options ...Option) Multiplexer {
	return &multiplexer{
		environment: newEnvironment(options...),
		mapping:     newMapping(),
	}
}

// Register is specified on the Multiplexer interface.
func (mux *multiplexer) Register(domain, resource string, handler ResourceHandler) error {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	err := handler.Init(mux.environment, domain, resource)
	if err != nil {
		return err
	}
	return mux.mapping.register(domain, resource, handler)
}

// RegisterAll is specified on the Multiplexer interface.
func (mux *multiplexer) RegisterAll(registrations Registrations) error {
	for _, registration := range registrations {
		err := mux.Register(registration.Domain, registration.Resource, registration.Handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deregister is specified on the Multiplexer interface.
func (mux *multiplexer) Deregister(domain, resource, id string) {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	mux.mapping.deregister(domain, resource, id)
}

// ServeHTTP is specified on the http.Handler interface.
func (mux *multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mutex.RLock()
	defer mux.mutex.RUnlock()
	job := newJob(mux.environment, r, w)
	measuring := monitoring.BeginMeasuring(job.String())
	defer measuring.EndMeasuring()
	if err := mux.mapping.handle(job); err != nil {
		mux.internalServerError("error handling request", job, err)
	}
}

// internalServerError logs an internal error and returns it to the user.
func (mux *multiplexer) internalServerError(format string, job Job, err error) {
	msg := fmt.Sprintf(format+" %q: %v", job, err)
	logger.Errorf(msg)
	http.Error(job.ResponseWriter(), msg, http.StatusInternalServerError)
}

// EOF
