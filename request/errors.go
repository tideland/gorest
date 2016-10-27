// Tideland Go REST Server Library - Request - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ErrServiceNotConfigured = iota + 1
	ErrCannotPrepareRequest
	ErrHTTPRequestFailed
	ErrProcessingRequestContent
	ErrReadingResponse
)

var errorMessages = errors.Messages{
	ErrServiceNotConfigured:     "service '%s' is not configured",
	ErrCannotPrepareRequest:     "cannot prepare request",
	ErrHTTPRequestFailed:        "HTTP request failed",
	ErrProcessingRequestContent: "cannot process request content",
	ErrReadingResponse:          "cannot read the HTTP response",
}

// EOF
