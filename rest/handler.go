// Tideland GoREST - REST - Handlers
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
	"fmt"
	"net/http"

	"github.com/tideland/golib/errors"
)

//--------------------
// RESOURCE HANDLER INTERFACES
//--------------------

// ResourceHandler is the base interface for all resource
// handlers understanding the REST verbs. It allows the
// initialization and returns an id that has to be unique
// for the combination of domain and resource. So it can
// later be removed again.
type ResourceHandler interface {
	// ID returns the deployment ID of the handler.
	ID() string

	// Init initializes the resource handler after registrations.
	Init(env Environment, domain, resource string) error
}

// GetResourceHandler is the additional interface for
// handlers understanding the verb GET.
type GetResourceHandler interface {
	Get(job Job) (bool, error)
}

// ReadResourceHandler is the additional interface for
// handlers understanding the verb GET but mapping it
// to method Read() according to the REST conventions.
type ReadResourceHandler interface {
	Read(job Job) (bool, error)
}

// HeadResourceHandler is the additional interface for
// handlers understanding the verb HEAD.
type HeadResourceHandler interface {
	Head(job Job) (bool, error)
}

// PutResourceHandler is the additional interface for
// handlers understanding the verb PUT.
type PutResourceHandler interface {
	Put(job Job) (bool, error)
}

// UpdateResourceHandler is the additional interface for
// handlers understanding the verb PUT but mapping it
// to method Update() according to the REST conventions.
type UpdateResourceHandler interface {
	Update(job Job) (bool, error)
}

// PostResourceHandler is the additional interface for
// handlers understanding the verb POST.
type PostResourceHandler interface {
	Post(job Job) (bool, error)
}

// CreateResourceHandler is the additional interface for
// handlers understanding the verb POST but mapping it
// to method Create() according to the REST conventions.
type CreateResourceHandler interface {
	Create(job Job) (bool, error)
}

// PatchResourceHandler is the additional interface for
// handlers understanding the verb PATCH.
type PatchResourceHandler interface {
	Patch(job Job) (bool, error)
}

// ModifyResourceHandler is the additional interface for
// handlers understanding the verb PATCH but mapping it
// to method Modify() according to the REST conventions.
type ModifyResourceHandler interface {
	Modify(job Job) (bool, error)
}

// DeleteResourceHandler is the additional interface for
// handlers understanding the verb DELETE.
type DeleteResourceHandler interface {
	Delete(job Job) (bool, error)
}

// OptionsResourceHandler is the additional interface for
// handlers understanding the verb OPTION.
type OptionsResourceHandler interface {
	Options(job Job) (bool, error)
}

// InfoResourceHandler is the additional interface for
// handlers understanding the verb OPTION but mapping it
// to method Info() according to the REST conventions.
type InfoResourceHandler interface {
	Info(job Job) (bool, error)
}

// handleJob dispatches the passed job to the right method of the
// passed handler. It always tries the nativ method first, then
// the alias method according to the REST conventions.
func handleJob(handler ResourceHandler, job Job) (bool, error) {
	switch job.Request().Method {
	case http.MethodGet:
		return handleGetJob(handler, job)
	case http.MethodHead:
		return handleHeadJob(handler, job)
	case http.MethodPut:
		return handlePutJob(handler, job)
	case http.MethodPost:
		return handlePostJob(handler, job)
	case http.MethodPatch:
		return handlePatchJob(handler, job)
	case http.MethodDelete:
		return handleDeleteJob(handler, job)
	case http.MethodOptions:
		return handleOptionsJob(handler, job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, job.Request().Method)
}

// handleGetJob handles a job containing a GET request.
func handleGetJob(handler ResourceHandler, job Job) (bool, error) {
	grh, ok := handler.(GetResourceHandler)
	if ok {
		return grh.Get(job)
	}
	rrh, ok := handler.(ReadResourceHandler)
	if ok {
		return rrh.Read(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// handleHeadJob handles a job containing a HEAD request.
func handleHeadJob(handler ResourceHandler, job Job) (bool, error) {
	hrh, ok := handler.(HeadResourceHandler)
	if ok {
		return hrh.Head(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// handlePutJob handles a job containing a PUT request.
func handlePutJob(handler ResourceHandler, job Job) (bool, error) {
	prh, ok := handler.(PutResourceHandler)
	if ok {
		return prh.Put(job)
	}
	urh, ok := handler.(UpdateResourceHandler)
	if ok {
		return urh.Update(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// handlePostJob handles a job containing a POST request.
func handlePostJob(handler ResourceHandler, job Job) (bool, error) {
	prh, ok := handler.(PostResourceHandler)
	if ok {
		return prh.Post(job)
	}
	crh, ok := handler.(CreateResourceHandler)
	if ok {
		return crh.Create(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// handlePatchJob handles a job containing a PATCH request.
func handlePatchJob(handler ResourceHandler, job Job) (bool, error) {
	prh, ok := handler.(PatchResourceHandler)
	if ok {
		return prh.Patch(job)
	}
	mrh, ok := handler.(ModifyResourceHandler)
	if ok {
		return mrh.Modify(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// handleDeleteJob handles a job containing a DELETE request.
func handleDeleteJob(handler ResourceHandler, job Job) (bool, error) {
	drh, ok := handler.(DeleteResourceHandler)
	if ok {
		return drh.Delete(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// handleOptionsJob handles a job containing an OPTIONS request.
func handleOptionsJob(handler ResourceHandler, job Job) (bool, error) {
	orh, ok := handler.(OptionsResourceHandler)
	if ok {
		return orh.Options(job)
	}
	irh, ok := handler.(InfoResourceHandler)
	if ok {
		return irh.Info(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, jobDescription(handler, job))
}

// jobDescription returns a description for possible errors.
func jobDescription(handler ResourceHandler, job Job) string {
	return fmt.Sprintf("%s %s@%s/%s", job.Request().Method, handler.ID(), job.Domain(), job.Resource())
}

// EOF
