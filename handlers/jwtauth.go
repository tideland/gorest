// Tideland Go REST Server Library - Handlers - JWT Authorization
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
	"time"

	"github.com/tideland/gorest/jwt"
	"github.com/tideland/gorest/rest"
)

//--------------------
// JWT AUTHORIZATION HANDLER
//--------------------

// GatekeeperFunc has to be defined by the handler user to check if the
// claims contain the required authorization information. In case it doesn't
// the function has to return false.
type GatekeeperFunc func(job rest.Job, claims jwt.Claims) (bool, error)

// jwtAuthorizationHandler checks for a valid token and then runs
// a gatekeeper function.
type jwtAuthorizationHandler struct {
	id         string
	key        jwt.Key
	gatekeeper GatekeeperFunc
}

// NewJWTAuthorizationHandler creates a handler checking for a valid JSON
// Web Token in each request. In case the request has one the configured
// gatekeeper function will be called with job and claims for further
// validation.
func NewjwtAuthorizationHandler(id string, key jwt.Key, gf GatekeeperFunc) rest.ResourceHandler {
	return &jwtAuthorizationHandler{id, key, gf}
}

// ID is specified on the ResourceHandler interface.
func (h *jwtAuthorizationHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *jwtAuthorizationHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *jwtAuthorizationHandler) Get(job rest.Job) (bool, error) {
	return h.check(job)
}

// Head is specified on the HeadResourceHandler interface.
func (h *jwtAuthorizationHandler) Head(job rest.Job) (bool, error) {
	return h.check(job)
}

// Put is specified on the PutResourceHandler interface.
func (h *jwtAuthorizationHandler) Put(job rest.Job) (bool, error) {
	return h.check(job)
}

// Post is specified on the PostResourceHandler interface.
func (h *jwtAuthorizationHandler) Post(job rest.Job) (bool, error) {
	return h.check(job)
}

// Patch is specified on the PatchResourceHandler interface.
func (h *jwtAuthorizationHandler) Patch(job rest.Job) (bool, error) {
	return h.check(job)
}

// Delete is specified on the DeleteResourceHandler interface.
func (h *jwtAuthorizationHandler) Delete(job rest.Job) (bool, error) {
	return h.check(job)
}

// Options is specified on the OptionsResourceHandler interface.
func (h *jwtAuthorizationHandler) Options(job rest.Job) (bool, error) {
	return h.check(job)
}

// check is used by all methods to check the token.
func (h *jwtAuthorizationHandler) check(job rest.Job) (bool, error) {
	token, err := jwt.VerifyFromJob(job, h.key)
	if err != nil {
		return false, h.deny(job, err.Error())
	}
	if !token.IsValid(time.Minute) {
		// TODO Configurable leeway.
		return false, h.deny(job, "invalid JSON web token")
	}
	return h.gatekeeper(job, token.Claims())
}

// deny sends a negative feedback to the caller.
func (h *jwtAuthorizationHandler) deny(job rest.Job, msg string) error {
	var f rest.Formatter
	if job.AcceptsContentType(rest.ContentTypeJSON) {
		f = job.JSON(true)
	} else {
		f = job.XML()
	}
	return rest.NegativeFeedback(f, msg)
}

// EOF
