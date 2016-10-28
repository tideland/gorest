// Tideland Go REST Server Library - Request - Errors
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

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
	ErrNoServerDefined = iota + 1
	ErrCannotPrepareRequest
	ErrHTTPRequestFailed
	ErrProcessingRequestContent
	ErrContentNotKeyValue
	ErrReadingResponse
)

var errorMessages = errors.Messages{
	ErrNoServerDefined:          "no server for domain '%s' configured",
	ErrCannotPrepareRequest:     "cannot prepare request",
	ErrHTTPRequestFailed:        "HTTP request failed",
	ErrProcessingRequestContent: "cannot process request content",
	ErrContentNotKeyValue:       "content is not key/value",
	ErrReadingResponse:          "cannot read the HTTP response",
}

// EOF
