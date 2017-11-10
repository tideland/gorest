// Tideland GoREST - REST - Multiplexer
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/etc"
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

	// RegisteredHandlers returns the ID stack of registered handlers
	// for a domain and resource.
	RegisteredHandlers(domain, resource string) []string

	// Deregister removes one, more, or all resource handler for a
	// given domain and resource.
	Deregister(domain, resource string, ids ...string)
}

// multiplexer implements the Multiplexer interface.
type multiplexer struct {
	mutex       sync.RWMutex
	environment *environment
	mapping     *mapping
}

// NewMultiplexer creates a new HTTP multiplexer. The passed context
// will be  used if a handler requests a context from a job, the
// configuration allows to configure the multiplexer. The allowed
// parameters are
//
//     {etc
//         {basepath /}
//         {default-domain default}
//         {default-resource default}
//         {ignore-favicon true}
//     }
//
// The values shown here are the default values if the configuration
// is nil or missing these settings.
func NewMultiplexer(ctx context.Context, cfg etc.Etc) Multiplexer {
	return &multiplexer{
		environment: newEnvironment(ctx, cfg),
		mapping:     newMapping(cfg),
	}
}

// Register implements the Multiplexer interface.
func (mux *multiplexer) Register(domain, resource string, handler ResourceHandler) error {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	err := handler.Init(mux.environment, domain, resource)
	if err != nil {
		return err
	}
	return mux.mapping.register(domain, resource, handler)
}

// RegisterAll implements the Multiplexer interface.
func (mux *multiplexer) RegisterAll(registrations Registrations) error {
	for _, registration := range registrations {
		err := mux.Register(registration.Domain, registration.Resource, registration.Handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisteredHandlers implements the Multiplexer interface.
func (mux *multiplexer) RegisteredHandlers(domain, resource string) []string {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	return mux.mapping.registeredHandlers(domain, resource)
}

// Deregister implements the Multiplexer interface.
func (mux *multiplexer) Deregister(domain, resource string, ids ...string) {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	mux.mapping.deregister(domain, resource, ids...)
}

// ServeHTTP implements the http.Handler interface.
func (mux *multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mutex.RLock()
	defer mux.mutex.RUnlock()
	job := newJob(mux.environment, r, w)
	measuring := monitoring.BeginMeasuring(job.String())
	defer measuring.EndMeasuring()
	if err := mux.mapping.handle(job); err != nil {
		mux.handleError("error handling request", job, err)
	}
}

// handleError logs an error and returns it to the user.
func (mux *multiplexer) handleError(format string, job Job, err error) {
	code := http.StatusInternalServerError
	msg := fmt.Sprintf(format+" %q: %v", job, err)
	logger.Errorf(msg)
	if errors.IsError(err, ErrMethodNotSupported) {
		code = http.StatusMethodNotAllowed
	}
	http.Error(job.ResponseWriter(), msg, code)
}

// EOF
