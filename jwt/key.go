// Tideland GoREST - JSON Web Token - Keys
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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"

	"github.com/tideland/golib/errors"
)

//--------------------
// KEY
//--------------------

// Key is the used key to sign a token. The real implementation
// controls signing and verification.
type Key interface{}

// ReadECPrivateKey reads a PEM formated ECDSA private key
// from the passed reader.
func ReadECPrivateKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New(ErrCannotReadPEM, errorMessages)
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, errors.New(ErrCannotDecodePEM, errorMessages)
	}
	var parsed *ecdsa.PrivateKey
	if parsed, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages)
	}
	return parsed, nil
}

// ReadECPublicKey reads a PEM encoded ECDSA public key
// from the passed reader.
func ReadECPublicKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New(ErrCannotReadPEM, errorMessages)
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, errors.New(ErrCannotDecodePEM, errorMessages)
	}
	var parsed interface{}
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotParseECDSA, errorMessages)
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New(ErrNoECDSAKey, errorMessages)
	}
	return publicKey, nil
}

// ReadRSAPrivateKey reads a PEM encoded PKCS1 or PKCS8 private key
// from the passed reader.
func ReadRSAPrivateKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New(ErrCannotReadPEM, errorMessages)
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, errors.New(ErrCannotDecodePEM, errorMessages)
	}
	var parsed interface{}
	parsed, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotParseRSA, errorMessages)
		}
	}
	privateKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New(ErrNoRSAKey, errorMessages)
	}
	return privateKey, nil
}

// ReadRSAPublicKey reads a PEM encoded PKCS1 or PKCS8 public key
// from the passed reader.
func ReadRSAPublicKey(r io.Reader) (Key, error) {
	pemkey, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New(ErrCannotReadPEM, errorMessages)
	}
	var block *pem.Block
	if block, _ = pem.Decode(pemkey); block == nil {
		return nil, errors.New(ErrCannotDecodePEM, errorMessages)
	}
	var parsed interface{}
	parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotParseRSA, errorMessages)
		}
		parsed = certificate.PublicKey
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New(ErrNoRSAKey, errorMessages)
	}
	return publicKey, nil
}

// EOF
