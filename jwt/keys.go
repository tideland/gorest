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
	"crypto/x509"
	// "crypto/rsa"
	"encoding/pem"

	"github.com/tideland/golib/errors"
)

//--------------------
// ECDSA
//--------------------

// PEMToECPrivateKey converts a PEM formatted ECDSA private
// key into its Go representation.
func PEMToECPrivateKey(key []byte) (*ecdsa.PrivateKey, error) {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New(ErrCannotDecodePEM, errorMessages)
	}
	var parsed *ecdsa.PrivateKey
	var err error
	if parsed, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages)
	}
	return parsed, nil
}

// PEMToECPublicKey converts a PEM formatted ECDSA public
// key into its Go representation.
func PEMToECPublicKey(key []byte) (*ecdsa.PublicKey, error) {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New(ErrCannotDecodePEM, errorMessages)
	}
	var parsed interface{}
	var err error
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages, "certificate")
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New(ErrCannotParseECDSA, errorMessages, "public key")
	}
	return publicKey, nil
}

// EOF
