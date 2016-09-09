// Tideland Go REST Server Library - JSON Web Token - Header
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"

	"github.com/tideland/gorest/rest"
)

//--------------------
// REQUEST HANDLING
//--------------------

// TokenFromJob retrieves a possible JWT from the request
// inside a REST job. The JWT is only decoded.
func TokenFromJob(job rest.Job) (JWT, error) {
	return TokenFromRequest(job.Request())
}

// TokenFromRequest retrieves a possible JWT from a
// HTTP request. The JWT is only decoded.
func TokenFromRequest(req *http.Request) (JWT, error) {
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return nil, nil
	}
	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, nil
	}
	return Decode(fields[1])
}

// EOF
