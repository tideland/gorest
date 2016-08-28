// Tideland Go REST Server Library - JSON Web Token - Algorithm
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
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"math/big"

	"github.com/tideland/golib/errors"
)

//--------------------
// INIT
//--------------------

// Assure linking of crypto pachages.
func init() {
	sha256.New()
	sha512.New()
}

//--------------------
// SIGNATURE
//--------------------

// Signature is the resulting signature when signing
// a token.
type Signature []byte

//--------------------
// KEY
//--------------------

// ecPoint is needed to marshal R and S of the ECDSA algorithms.
type ecPoint struct {
	R *big.Int
	S *big.Int
}

// Key is the used key to sign a token.
type Key interface{}

// sign signs the passed data based on the key and the passed hash.
func sign(data []byte, k Key, h crypto.Hash, isRSAPSS bool) (Signature, error) {
	hashSum := func() []byte {
		hasher := h.New()
		hasher.Write(data)
		return hasher.Sum(nil)
	}
	switch key := k.(type) {
	case *ecdsa.PrivateKey:
		// ECDSA algorithms.
		r, s, err := ecdsa.Sign(rand.Reader, key, hashSum())
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
		}
		sig, err := asn1.Marshal(ecPoint{r, s})
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
		}
		return Signature(sig), nil
	case []byte:
		// HMAC algorithms.
		hasher := hmac.New(h.New, key)
		hasher.Write(data)
		sig := hasher.Sum(nil)
		return Signature(sig), nil
	case *rsa.PrivateKey:
		// RSA and RSAPSS algorithms.
		if isRSAPSS {
			// RSAPSS.
			options := &rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthAuto,
				Hash:       h,
			}
			sig, err := rsa.SignPSS(rand.Reader, key, h, hashSum(), options)
			if err != nil {
				return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
			}
			return Signature(sig), nil
		} else {
			// RSA.
			sig, err := rsa.SignPKCS1v15(rand.Reader, key, h, hashSum())
			if err != nil {
				return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
			}
			return Signature(sig), nil
		}
	default:
		// No valid key type.
		return nil, errors.New(ErrInvalidKeyType, errorMessages, k)
	}
}

// verify checks if the signature is correct for the passed data
// based on the key and the passed hash.
func verify(data []byte, sig Signature, k Key, h crypto.Hash, isRSAPSS bool) error {
	hashSum := func() []byte {
		hasher := h.New()
		hasher.Write(data)
		return hasher.Sum(nil)
	}
	switch key := k.(type) {
	case *ecdsa.PublicKey:
		// ECDSA algorithms.
		var ecp ecPoint
		if _, err := asn1.Unmarshal(sig, &ecp); err != nil {
			return errors.Annotate(err, ErrCannotVerify, errorMessages)
		}
		if !ecdsa.Verify(key, hashSum(), ecp.R, ecp.S) {
			return errors.New(ErrInvalidSignature, errorMessages)
		}
		return nil
	case []byte:
		// HMAC algorithms.
		expectedSig, err := sign(data, k, h, isRSAPSS)
		if err != nil {
			return errors.Annotate(err, ErrCannotVerify, errorMessages)
		}
		if !hmac.Equal(sig, expectedSig) {
			return errors.New(ErrInvalidSignature, errorMessages)
		}
		return nil
	case *rsa.PublicKey:
		// RSA and RSAPSS algorithms.
		if isRSAPSS {
			// RSAPSS.
			options := &rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthAuto,
				Hash:       h,
			}
			if err := rsa.VerifyPSS(key, h, hashSum(), sig, options); err != nil {
				return errors.Annotate(err, ErrInvalidSignature, errorMessages)
			}
		} else {
			// RSA.
			if err := rsa.VerifyPKCS1v15(key, h, hashSum(), sig); err != nil {
				return errors.Annotate(err, ErrInvalidSignature, errorMessages)
			}
		}
		return nil
	default:
		// No valid key type.
		return errors.New(ErrInvalidKeyType, errorMessages, k)
	}
}

//--------------------
// ALGORITHM
//--------------------

// Algorithm describes the algorithm used to sign a token.
type Algorithm string

const (
	ES256 Algorithm = "ES256"
	ES384 Algorithm = "ES384"
	ES512 Algorithm = "ES512"
	HS256 Algorithm = "HS256"
	HS384 Algorithm = "HS384"
	HS512 Algorithm = "HS512"
	PS256 Algorithm = "PS256"
	PS384 Algorithm = "PS384"
	PS512 Algorithm = "PS512"
	RS256 Algorithm = "RS256"
	RS384 Algorithm = "RS384"
	RS512 Algorithm = "RS512"
)

// Sign creates the signature for the data based on the
// algorithm and the key.
func (a Algorithm) Sign(data []byte, key Key) (Signature, error) {
	switch a {
	case ES256, HS256, PS256, RS256:
		return sign(data, key, crypto.SHA256, a.isRSAPSS())
	case ES384, HS384, PS384, RS384:
		return sign(data, key, crypto.SHA384, a.isRSAPSS())
	case ES512, HS512, PS512, RS512:
		return sign(data, key, crypto.SHA512, a.isRSAPSS())
	default:
		return nil, errors.New(ErrInvalidAlgorithm, errorMessages, a)
	}
}

// Verify checks if the signature is correct for the data when using
// the passed key.
func (a Algorithm) Verify(data []byte, sig Signature, key Key) error {
	switch a {
	case ES256, HS256, PS256, RS256:
		return verify(data, sig, key, crypto.SHA256, a.isRSAPSS())
	case ES384, HS384, PS384, RS384:
		return verify(data, sig, key, crypto.SHA384, a.isRSAPSS())
	case ES512, HS512, PS512, RS512:
		return verify(data, sig, key, crypto.SHA512, a.isRSAPSS())
	default:
		return errors.New(ErrInvalidAlgorithm, errorMessages, a)
	}
}

// isRSAPSS returns true when the algorithm is one of
// the RSAPSS algorithms.
func (a Algorithm) isRSAPSS() bool {
	return a[0] == 'P'
}

// EOF
