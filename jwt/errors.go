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
	ErrCannotSign
	ErrCannotVerify
	ErrNoKey
	ErrJSONMarshalling
	ErrJSONUnmarshalling
	ErrInvalidTokenPart
	ErrInvalidAlgorithm
	ErrInvalidKeyType
	ErrInvalidSignature
	ErrCannotReadPEM
	ErrCannotDecodePEM
	ErrCannotParseECDSA
	ErrNoECDSAKey
	ErrCannotParseRSA
	ErrNoRSAKey
)

var errorMessages = errors.Messages{
	ErrCannotEncode:      "cannot encode the %s",
	ErrCannotDecode:      "cannot decode the %s",
	ErrCannotSign:        "cannot sign the token",
	ErrCannotVerify:      "cannot verify the %s",
	ErrNoKey:             "no key available, only after encoding or verifying",
	ErrJSONMarshalling:   "error marshalling to JSON",
	ErrJSONUnmarshalling: "error unmarshilling from JSON",
	ErrInvalidTokenPart:  "part of the token contains invalid data",
	ErrInvalidAlgorithm:  "signature algorithm %q is invalid",
	ErrInvalidKeyType:    "key type %#v is invalid",
	ErrInvalidSignature:  "token signature is invalid",
	ErrCannotReadPEM:     "cannot read the PEM",
	ErrCannotDecodePEM:   "cannot decode the PEM",
	ErrCannotParseECDSA:  "cannot parse the ECDSA",
	ErrNoECDSAKey:        "passed key is no ECDSA key",
	ErrCannotParseRSA:    "cannot parse the RSA",
	ErrNoRSAKey:          "passed key is no RSA key",
}

// EOF
