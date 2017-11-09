// Tideland GoREST - REST - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the rest package.
const (
	ErrDuplicateHandler = iota + 1
	ErrInitHandler
	ErrIllegalRequest
	ErrNoHandler
	ErrNoGetHandler
	ErrNoHeadHandler
	ErrNoPutHandler
	ErrNoPostHandler
	ErrNoPatchHandler
	ErrNoDeleteHandler
	ErrNoOptionsHandler
	ErrMethodNotSupported
	ErrUploadingFile
	ErrInvalidContentType
	ErrNoCachedTemplate
	ErrQueryValueNotFound
	ErrNoServerDefined
	ErrCannotPrepareRequest
	ErrHTTPRequestFailed
	ErrProcessingRequestContent
	ErrContentNotKeyValue
	ErrReadingResponse
)

var errorMessages = errors.Messages{
	ErrDuplicateHandler:         "cannot register handler %q, it is already registered",
	ErrInitHandler:              "error during initialization of handler %q",
	ErrIllegalRequest:           "illegal request containing too many parts",
	ErrNoHandler:                "found no handler with ID %q",
	ErrNoGetHandler:             "handler %q is no handler for GET requests",
	ErrNoHeadHandler:            "handler %q is no handler for HEAD requests",
	ErrNoPutHandler:             "handler %q is no handler for PUT requests",
	ErrNoPostHandler:            "handler %q is no handler for POST requests",
	ErrNoPatchHandler:           "handler %q is no handler for PATCH requests",
	ErrNoDeleteHandler:          "handler %q is no handler for DELETE requests",
	ErrNoOptionsHandler:         "handler %q is no handler for OPTIONS requests",
	ErrMethodNotSupported:       "method %q is not supported",
	ErrUploadingFile:            "uploaded file cannot be handled by %q",
	ErrInvalidContentType:       "content type is not %q",
	ErrNoCachedTemplate:         "template %q is not cached",
	ErrQueryValueNotFound:       "query value not found",
	ErrNoServerDefined:          "no server for domain '%s' configured",
	ErrCannotPrepareRequest:     "cannot prepare request",
	ErrHTTPRequestFailed:        "HTTP request failed",
	ErrProcessingRequestContent: "cannot process request content",
	ErrContentNotKeyValue:       "content is not key/value",
	ErrReadingResponse:          "cannot read the HTTP response",
}

// EOF
