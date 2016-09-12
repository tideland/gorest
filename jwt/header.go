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

// DecodeFromJob retrieves a possible JWT from the request
// inside a REST job. The JWT is only decoded.
func DecodeFromJob(job rest.Job) (JWT, error) {
	return DecodeFromRequest(job.Request())
}

// DecodeFromRequest retrieves a possible JWT from a
// HTTP request. The JWT is only decoded.
func DecodeFromRequest(req *http.Request) (JWT, error) {
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

// VerifyFromJob retrieves a possible JWT from
// the request inside a REST job. The JWT is verified.
func VerifyFromJob(job rest.Job, key Key) (JWT, error) {
	return VerifyFromRequest(job.Request(), key)
}

// VerifyFromRequest retrieves a possible JWT from a
// HTTP request. The JWT is verified.
func VerifyFromRequest(req *http.Request, key Key) (JWT, error) {
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return nil, nil
	}
	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, nil
	}
	return Verify(fields[1], key)
}

// AddTokenToRequest adds a token as header to a request for
// usage by a client.
func AddTokenToRequest(req *http.Request, jwt JWT) *http.Request {
	req.Header.Add("Authorization", "Bearer "+jwt.String())
	return req
}

// EOF
