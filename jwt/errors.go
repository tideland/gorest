// Tideland GoREST - JSON Web Token - Errors
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
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

// Error codes of the JWT package.
const (
	ErrCannotEncode = iota + 1
	ErrCannotDecode
	ErrCannotSign
	ErrCannotVerify
	ErrNoKey
	ErrJSONMarshalling
	ErrJSONUnmarshalling
	ErrInvalidTokenPart
	ErrInvalidCombination
	ErrInvalidAlgorithm
	ErrInvalidKeyType
	ErrInvalidSignature
	ErrCannotReadPEM
	ErrCannotDecodePEM
	ErrCannotParseECDSA
	ErrNoECDSAKey
	ErrCannotParseRSA
	ErrNoRSAKey
	ErrNoAuthorizationHeader
	ErrInvalidAuthorizationHeader
)

var errorMessages = errors.Messages{
	ErrCannotEncode:               "cannot encode the %s",
	ErrCannotDecode:               "cannot decode the %s",
	ErrCannotSign:                 "cannot sign the token",
	ErrCannotVerify:               "cannot verify the %s",
	ErrNoKey:                      "no key available, only after encoding or verifying",
	ErrJSONMarshalling:            "error marshalling to JSON",
	ErrJSONUnmarshalling:          "error unmarshalling from JSON",
	ErrInvalidTokenPart:           "part of the token contains invalid data",
	ErrInvalidCombination:         "invalid combination of algorithm %q and key type %q",
	ErrInvalidAlgorithm:           "signature algorithm %q is invalid",
	ErrInvalidKeyType:             "key type %T is invalid",
	ErrInvalidSignature:           "token signature is invalid",
	ErrCannotReadPEM:              "cannot read the PEM",
	ErrCannotDecodePEM:            "cannot decode the PEM",
	ErrCannotParseECDSA:           "cannot parse the ECDSA",
	ErrNoECDSAKey:                 "passed key is no ECDSA key",
	ErrCannotParseRSA:             "cannot parse the RSA",
	ErrNoRSAKey:                   "passed key is no RSA key",
	ErrNoAuthorizationHeader:      "request contains no authorization header",
	ErrInvalidAuthorizationHeader: "invalid authorization header: '%s'",
}

// EOF
