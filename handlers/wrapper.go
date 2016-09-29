// Tideland Go REST Server Library - Handlers - Wrapper
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"

	"github.com/tideland/gorest/rest"
)

//--------------------
// WRAPPER HANDLER
//--------------------

// wrapperHandler wraps existing handler functions for a usage inside
// the rest library.
type wrapperHandler struct {
	id     string
	handle http.HandlerFunc
}

// NewWrapperHandler creates a new wrapper around a standard
// handler function.
func NewWrapperHandler(id string, hf http.HandlerFunc) rest.ResourceHandler {
	return &wrapperHandler{id, hf}
}

// ID is specified on the ResourceHandler interface.
func (h *wrapperHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *wrapperHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *wrapperHandler) Get(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// Head is specified on the HeadResourceHandler interface.
func (h *wrapperHandler) Head(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// Put is specified on the PutResourceHandler interface.
func (h *wrapperHandler) Put(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// Post is specified on the PostResourceHandler interface.
func (h *wrapperHandler) Post(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// Patch is specified on the PatchResourceHandler interface.
func (h *wrapperHandler) Patch(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// Delete is specified on the DeleteResourceHandler interface.
func (h *wrapperHandler) Delete(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// Options is specified on the OptionsResourceHandler interface.
func (h *wrapperHandler) Options(job rest.Job) (bool, error) {
	h.handle(job.ResponseWriter(), job.Request())
	return true, nil
}

// EOF
