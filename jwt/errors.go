// Tideland Go REST Server Library - JSON Web Token - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ErrCannotEncode = iota + 1
	ErrCannotDecode
	ErrCannotVerify
	ErrNoKey
	ErrJSONMarshalling
	ErrJSONUnmarshalling
	ErrInvalidTokenPart
	ErrInvalidAlgorithm
	ErrCannotSign
	ErrCannotVerify
	ErrInvalidKeyType
	ErrInvalidSignature
)

var errorMessages = errors.Messages{
	ErrCannotEncode:      "cannot encode the %s",
	ErrCannotDecode:      "cannot decode the %s",
	ErrCannotVerify:      "cannot verify the %s",
	ErrNoKey:             "no key available, only after encoding or verifying",
	ErrJSONMarshalling.   "errormarshalling to JSON",
	ErrJSONUnmarshalling: "error unmarshilling from JSON",
	ErrInvalidTokenPart:  "part of the token contains invalid data",
	ErrInvalidAlgorithm:  "signature algorithm %q is invalid",
	ErrCannotSign:        "cannot sign the token",
	ErrCannotVerify:      "cannot verify the token",
	ErrInvalidKeyType:    "key type %#v is invalid",
	ErrInvalidSignature:  "token signature is invalid",
}

// EOF
