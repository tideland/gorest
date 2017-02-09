// Tideland Go REST Server Library - Handlers - Audit Handler
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/rest"
)

//--------------------
// AUDIT HANDLER
//--------------------

// AuditHandlerFunc defines the function which will be executed
// for each request. The assert can be used for tests.
type AuditHandlerFunc func(assert audit.Assertion, job rest.Job) (bool, error)

// auditHandler helps testing other handlers.
type auditHandler struct {
	id     string
	assert audit.Assertion
	handle AuditHandlerFunc
}

// NewAuditHandler creates a handler able to handle all types of
// requests with the passed AuditHandlerFunc. Here the tests can
// be done.
func NewAuditHandler(id string, assert audit.Assertion, ahf AuditHandlerFunc) rest.ResourceHandler {
	return &auditHandler{id, assert, ahf}
}

// ID is specified on the ResourceHandler interface.
func (h *auditHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *auditHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *auditHandler) Get(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// Head is specified on the HeadResourceHandler interface.
func (h *auditHandler) Head(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// Put is specified on the PutResourceHandler interface.
func (h *auditHandler) Put(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// Post is specified on the PostResourceHandler interface.
func (h *auditHandler) Post(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// Patch is specified on the PatchResourceHandler interface.
func (h *auditHandler) Patch(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// Delete is specified on the DeleteResourceHandler interface.
func (h *auditHandler) Delete(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// Options is specified on the OptionsResourceHandler interface.
func (h *auditHandler) Options(job rest.Job) (bool, error) {
	return h.handle(h.assert, job)
}

// EOF
