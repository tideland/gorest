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
// REQUEST AND JOB HANDLING
//--------------------

// AddToRequest adds a token as header to a request for
// usage by a client.
func AddToRequest(req *http.Request, jwt JWT) *http.Request {
	req.Header.Add("Authorization", "Bearer "+jwt.String())
	return req
}

// DecodeFromRequest tries to retrieve a token from a request
// header. 
func DecodeFromRequest(req *http.Request) (JWT, error) {
	return nil, nil
}

// DecodeFromJob retrieves a possible JWT from
// the request inside a REST job. The JWT is only decoded.
func DecodeFromJob(job rest.Job) (JWT, error) {
	return retrieveFromJob(job, nil, nil)
}

// DecodeCachedFromJob retrieves a possible JWT from the request
// inside a REST job and checks if it already is cached. The JWT is
// only decoded. In case of no error the token is added to the cache.
func DecodeCachedFromJob(job rest.Job, cache Cache) (JWT, error) {
	return retrieveFromJob(job, cache, nil)
}

// VerifyFromJob retrieves a possible JWT from
// the request inside a REST job. The JWT is verified.
func VerifyFromJob(job rest.Job, key Key) (JWT, error) {
	return retrieveFromJob(job, nil, key)
}

// VerifyCachedFromJob retrieves a possible JWT from the request
// inside a REST job and checks if it already is cached. The JWT is
// verified. In case of no error the token is added to the cache.
func VerifyCachedFromJob(job rest.Job, cache Cache, key Key) (JWT, error) {
	return retrieveFromJob(job, cache, key)
}

//--------------------
// PRIVATE HELPERS
//--------------------

// retrieveFromRequest is the generic retrieval function with possible
// caching and verification.
func retrieveFromRequest(req *http.Request, cache Cache, key Key) (JWT, error) {
	// Retrieve token from header.
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		// TODO(mue): Add error. 
		return nil, nil
	}
	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "Bearer" {
		// TODO(mue): Add error. 
		return nil, nil
	}
	// Check cache.
	if cache != nil {
		jwt, ok := cache.Get(fields[1])
		if ok {
			return jwt, nil
		}
	}
	// Decode or verify.
	var jwt JWT
	var err error
	if key == nil {
		jwt, err = Decode(fields[1])
	} else {
		jwt, err = Verify(fields[1], key)
	}
	if err != nil {
		return nil, err
	}
	// Add to cache and return.
	if cache != nil {
		cache.Put(jwt)
	}
	return jwt, nil
}

// EOF
