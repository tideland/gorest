// Tideland Go REST Server Library - Request - Errors
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
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

// Error codes of the request package.
const (
	ErrNoServerDefined = iota + 1
	ErrCannotPrepareRequest
	ErrHTTPRequestFailed
	ErrProcessingRequestContent
	ErrInvalidContent
	ErrAnalyzingResponse
	ErrDecodingResponse
	ErrInvalidContentType
)

var errorMessages = errors.Messages{
	ErrNoServerDefined:          "no server for domain '%s' configured",
	ErrCannotPrepareRequest:     "cannot prepare request",
	ErrHTTPRequestFailed:        "HTTP request failed",
	ErrProcessingRequestContent: "cannot process request content",
	ErrInvalidContent:           "content invalid for URL encoding",
	ErrAnalyzingResponse:        "cannot analyze the HTTP response",
	ErrDecodingResponse:         "cannot decode the HTTP response",
	ErrInvalidContentType:       "invalid content type '%s'",
}

// EOF
