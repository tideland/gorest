// Tideland Go REST Server Library - REST - Handlers
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

// PostResourceHandler is the additional interface for
// handlers understanding the verb POST.
type PostResourceHandler interface {
	Post(job Job) (bool, error)
}

// PatchResourceHandler is the additional interface for
// handlers understanding the verb PATCH.
type PatchResourceHandler interface {
	Patch(job Job) (bool, error)
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

// handleJob dispatches the passed job to the right method of the
// passed handler.
func handleJob(handler ResourceHandler, job Job) (bool, error) {
	id := func() string {
		return fmt.Sprintf("%s@%s/%s", handler.ID(), job.Domain(), job.Resource())
	}
	switch job.Request().Method {
	case "GET":
		grh, ok := handler.(GetResourceHandler)
		if !ok {
			return false, errors.New(ErrNoGetHandler, errorMessages, id())
		}
		return grh.Get(job)
	case "HEAD":
		hrh, ok := handler.(HeadResourceHandler)
		if !ok {
			return false, errors.New(ErrNoHeadHandler, errorMessages, id())
		}
		return hrh.Head(job)
	case "PUT":
		prh, ok := handler.(PutResourceHandler)
		if !ok {
			return false, errors.New(ErrNoPutHandler, errorMessages, id())
		}
		return prh.Put(job)
	case "POST":
		prh, ok := handler.(PostResourceHandler)
		if !ok {
			return false, errors.New(ErrNoPostHandler, errorMessages, id())
		}
		return prh.Post(job)
	case "PATCH":
		prh, ok := handler.(PatchResourceHandler)
		if !ok {
			return false, errors.New(ErrNoPatchHandler, errorMessages, id())
		}
		return prh.Patch(job)
	case "DELETE":
		drh, ok := handler.(DeleteResourceHandler)
		if !ok {
			return false, errors.New(ErrNoDeleteHandler, errorMessages, id())
		}
		return drh.Delete(job)
	case "OPTIONS":
		orh, ok := handler.(OptionsResourceHandler)
		if !ok {
			return false, errors.New(ErrNoOptionsHandler, errorMessages, id())
		}
		return orh.Options(job)
	}
	return false, errors.New(ErrMethodNotSupported, errorMessages, job.Request().Method)
}

// EOF
