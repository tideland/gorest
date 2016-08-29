// Tideland Go REST Server Library - JSON Web Token - Keys
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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/X509"
	"encoding/pem"

	"github.com/tideland/golib/errors"
)

//--------------------
// ECDSA
//--------------------

// PEMToECPrivateKey converts a PEM formatted ECDSA private
// key into its Go representation.
func PEMToECPrivateKey(pemKey string) (*ecdsa.PrivateKey, error) {
	if pemBlock, err := pem.Decode(pemKey); err != nil {
		return nil, errors.Annotate(err, ErrCannotDecodePEM, errorMessages)
	}
	if parsedKey, err := x509.ParseECPrivateKey(pemBlock.Bytes); err != nil {
		return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages)
	}
	return parsedKey, nil
}

// PEMToECPublicKey converts a PEM formatted ECDSA public
// key into its Go representation.
func PEMToECPrivateKey(pemKey string) (*ecdsa.PublicKey, error) {
	if pemBlock, err := pem.Decode(pemKey); err != nil {
		return nil, errors.Annotate(err, ErrCannotDecodePEM, errorMessages)
	}
	if parsedKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes); err != nil {
		return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages, "DER public key")
	}
	if cert, err := x509.ParseCertificate(pemBlock.Bytes); err != nil {
		return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages, "certificate")
	}
	return parsedKey, nil
}

// EOF