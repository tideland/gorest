// Tideland GoREST - JSON Web Token - Algorithm
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
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"
	"math/big"

	// Import hashing packages just to register them via init().
	_ "crypto/sha256"
	_ "crypto/sha512"

	"github.com/tideland/golib/errors"
)

//--------------------
// SIGNATURE
//--------------------

// Signature is the resulting signature when signing
// a token.
type Signature []byte

//--------------------
// ALGORITHM
//--------------------

// Algorithm describes the algorithm used to sign a token.
type Algorithm string

// Definition of the supported algorithms.
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
	NONE  Algorithm = "none"
)

// ecPoint is needed to marshal R and S of the ECDSA algorithms.
type ecPoint struct {
	R *big.Int
	S *big.Int
}

// Sign creates the signature for the data based on the
// algorithm and the key.
func (a Algorithm) Sign(data []byte, key Key) (Signature, error) {
	switch a {
	case ES256, HS256, PS256, RS256:
		return a.sign(data, key, crypto.SHA256)
	case ES384, HS384, PS384, RS384:
		return a.sign(data, key, crypto.SHA384)
	case ES512, HS512, PS512, RS512:
		return a.sign(data, key, crypto.SHA512)
	case NONE:
		return a.sign(data, key, 0)
	default:
		return nil, errors.New(ErrInvalidAlgorithm, errorMessages, a)
	}
}

// Verify checks if the signature is correct for the data when using
// the passed key.
func (a Algorithm) Verify(data []byte, sig Signature, key Key) error {
	switch a {
	case ES256, HS256, PS256, RS256:
		return a.verify(data, sig, key, crypto.SHA256)
	case ES384, HS384, PS384, RS384:
		return a.verify(data, sig, key, crypto.SHA384)
	case ES512, HS512, PS512, RS512:
		return a.verify(data, sig, key, crypto.SHA512)
	case NONE:
		return a.verify(data, sig, key, 0)
	default:
		return errors.New(ErrInvalidAlgorithm, errorMessages, a)
	}
}

// isRSAPSS returns true when the algorithm is one of
// the RSAPSS algorithms.
func (a Algorithm) isRSAPSS() bool {
	return a[0] == 'P'
}

// sign signs the passed data based on the key and the passed hash.
func (a Algorithm) sign(data []byte, k Key, h crypto.Hash) (Signature, error) {
	switch key := k.(type) {
	case *ecdsa.PrivateKey:
		// ECDSA algorithms.
		return a.signECDSA(data, key, h)
	case []byte:
		// HMAC algorithms.
		return a.signHMAC(data, key, h)
	case *rsa.PrivateKey:
		// RSA and RSAPSS algorithms.
		return a.signRSA(data, key, h)
	case string:
		// None algorithm.
		if a != "none" {
			return nil, errors.New(ErrInvalidCombination, errorMessages, a, "none")
		}
		return Signature(""), nil
	default:
		// No valid key type.
		return nil, errors.New(ErrInvalidKeyType, errorMessages, k)
	}
}

// signECDSA signs the data using the ECDSA algorithm.
func (a Algorithm) signECDSA(data []byte, key *ecdsa.PrivateKey, h crypto.Hash) (Signature, error) {
	if a[0] != 'E' {
		return nil, errors.New(ErrInvalidCombination, errorMessages, a, "ECDSA")
	}
	r, s, err := ecdsa.Sign(rand.Reader, key, hashSum(data, h))
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
	}
	sig, err := asn1.Marshal(ecPoint{r, s})
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
	}
	return Signature(sig), nil
}

// signHMAC signs the data using the HMAC algorithm.
func (a Algorithm) signHMAC(data, key []byte, h crypto.Hash) (Signature, error) {
	if a[0] != 'H' {
		return nil, errors.New(ErrInvalidCombination, errorMessages, a, "HMAC")
	}
	hasher := hmac.New(h.New, key)
	hasher.Write(data)
	sig := hasher.Sum(nil)
	return Signature(sig), nil
}

// signRSA signs the data using the RSAPSS or RSA algorithm.
func (a Algorithm) signRSA(data []byte, key *rsa.PrivateKey, h crypto.Hash) (Signature, error) {
	if a[0] != 'P' && a[0] != 'R' {
		return nil, errors.New(ErrInvalidCombination, errorMessages, a, "RSA(PSS)")
	}
	if a.isRSAPSS() {
		// RSAPSS.
		options := &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       h,
		}
		sig, err := rsa.SignPSS(rand.Reader, key, h, hashSum(data, h), options)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
		}
		return Signature(sig), nil
	}
	// RSA.
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, h, hashSum(data, h))
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotSign, errorMessages)
	}
	return Signature(sig), nil
}

// verify checks if the signature is correct for the passed data
// based on the key and the passed hash.
func (a Algorithm) verify(data []byte, sig Signature, k Key, h crypto.Hash) error {
	switch key := k.(type) {
	case *ecdsa.PublicKey:
		// ECDSA algorithms.
		return a.verifyECDSA(data, sig, key, h)
	case []byte:
		// HMAC algorithms.
		return a.verifyHMAC(data, sig, key, h)
	case *rsa.PublicKey:
		// RSA and RSAPSS algorithms.
		return a.verifyRSA(data, sig, key, h)
	case string:
		// None algorithm.
		if a != "none" {
			return errors.New(ErrInvalidCombination, errorMessages, a, "none")
		}
		if len(sig) > 0 {
			return errors.New(ErrInvalidSignature, errorMessages)
		}
		return nil
	default:
		// No valid key type.
		return errors.New(ErrInvalidKeyType, errorMessages, k)
	}
}

// verifyECDSA verifies the data using the ECDSA algorithm.
func (a Algorithm) verifyECDSA(data []byte, sig Signature, key *ecdsa.PublicKey, h crypto.Hash) error {
	if a[0] != 'E' {
		return errors.New(ErrInvalidCombination, errorMessages, a, "ECDSA")
	}
	var ecp ecPoint
	if _, err := asn1.Unmarshal(sig, &ecp); err != nil {
		return errors.Annotate(err, ErrCannotVerify, errorMessages)
	}
	if !ecdsa.Verify(key, hashSum(data, h), ecp.R, ecp.S) {
		return errors.New(ErrInvalidSignature, errorMessages)
	}
	return nil
}

// verifyHMAC verifies the data using the HMAC algorithm.
func (a Algorithm) verifyHMAC(data []byte, sig Signature, key []byte, h crypto.Hash) error {
	if a[0] != 'H' {
		return errors.New(ErrInvalidCombination, errorMessages, a, "HMAC")
	}
	expectedSig, err := a.sign(data, key, h)
	if err != nil {
		return errors.Annotate(err, ErrCannotVerify, errorMessages)
	}
	if !hmac.Equal(sig, expectedSig) {
		return errors.New(ErrInvalidSignature, errorMessages)
	}
	return nil
}

// verifyRSA verifies the data using the RSAPSS or RSS algorithm.
func (a Algorithm) verifyRSA(data []byte, sig Signature, key *rsa.PublicKey, h crypto.Hash) error {
	if a[0] != 'P' && a[0] != 'R' {
		return errors.New(ErrInvalidCombination, errorMessages, a, "RSA(PSS)")
	}
	if a.isRSAPSS() {
		// RSAPSS.
		options := &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       h,
		}
		if err := rsa.VerifyPSS(key, h, hashSum(data, h), sig, options); err != nil {
			return errors.Annotate(err, ErrInvalidSignature, errorMessages)
		}
	} else {
		// RSA.
		if err := rsa.VerifyPKCS1v15(key, h, hashSum(data, h), sig); err != nil {
			return errors.Annotate(err, ErrInvalidSignature, errorMessages)
		}
	}
	return nil
}

//--------------------
// HELPERS
//--------------------

// hashSum determines the hash sum of the passed data.
func hashSum(data []byte, h crypto.Hash) []byte {
	hasher := h.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// EOF
