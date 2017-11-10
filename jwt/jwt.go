// Tideland GoREST - JSON Web Token
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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// CONTEXT
//--------------------

// key for the storage of values in a context.
type key int

const (
	jwtKey key = iota
)

// NewContext returns a new context that carries a token.
func NewContext(ctx context.Context, token JWT) context.Context {
	return context.WithValue(ctx, jwtKey, token)
}

// FromContext returns the token stored in ctx, if any.
func FromContext(ctx context.Context) (JWT, bool) {
	token, ok := ctx.Value(jwtKey).(JWT)
	return token, ok
}

//--------------------
// JSON Web Token
//--------------------

// JWT describes the interface to access the parts of a
// JSON Web Token.
type JWT interface {
	// Stringer provides the String() method.
	fmt.Stringer

	// Claims returns the claims payload of the token.
	Claims() Claims

	// Key return the key of the token only when
	// it is a result of encoding or verification.
	Key() (Key, error)

	// Algorithm returns the algorithm of the token
	// after encoding, decoding, or verification.
	Algorithm() Algorithm

	// IsValid is a convenience method checking the
	// registered claims if the token is valid.
	IsValid(leeway time.Duration) bool
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type jwt struct {
	claims    Claims
	key       Key
	algorithm Algorithm
	token     string
}

// Encode creates a JSON Web Token for the given claims
// based on key and algorithm.
func Encode(claims Claims, key Key, algorithm Algorithm) (JWT, error) {
	jwt := &jwt{
		claims:    claims,
		key:       key,
		algorithm: algorithm,
	}
	headerPart, err := marshallAndEncode(jwtHeader{string(algorithm), "JWT"})
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotEncode, errorMessages, "header")
	}
	claimsPart, err := marshallAndEncode(claims)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotEncode, errorMessages, "claims")
	}
	dataParts := headerPart + "." + claimsPart
	signaturePart, err := signAndEncode([]byte(dataParts), key, algorithm)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotEncode, errorMessages, "signature")
	}
	jwt.token = dataParts + "." + signaturePart
	return jwt, nil
}

// Decode creates a token out of a string without verification.
func Decode(token string) (JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New(ErrCannotDecode, errorMessages, "parts")
	}
	var header jwtHeader
	err := decodeAndUnmarshall(parts[0], &header)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotDecode, errorMessages, "header")
	}
	var claims Claims
	err = decodeAndUnmarshall(parts[1], &claims)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotDecode, errorMessages, "claims")
	}
	return &jwt{
		claims:    claims,
		algorithm: Algorithm(header.Algorithm),
		token:     token,
	}, nil
}

// Verify creates a token out of a string and varifies it against
// the passed key.
func Verify(token string, key Key) (JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New(ErrCannotVerify, errorMessages, "parts")
	}
	var header jwtHeader
	err := decodeAndUnmarshall(parts[0], &header)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotVerify, errorMessages, "header")
	}
	err = decodeAndVerify(parts, key, Algorithm(header.Algorithm))
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotVerify, errorMessages, "signature")
	}
	var claims Claims
	err = decodeAndUnmarshall(parts[1], &claims)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotVerify, errorMessages, "claims")
	}
	return &jwt{
		claims:    claims,
		key:       key,
		algorithm: Algorithm(header.Algorithm),
		token:     token,
	}, nil
}

// Claims implements the JWT interface.
func (jwt *jwt) Claims() Claims {
	return jwt.claims
}

// Key implements the JWT interface.
func (jwt *jwt) Key() (Key, error) {
	if jwt.key == nil {
		return nil, errors.New(ErrNoKey, errorMessages)
	}
	return jwt.key, nil
}

// Algorithm implements the JWT interface.
func (jwt *jwt) Algorithm() Algorithm {
	return jwt.algorithm
}

// IsValid implements the JWT interface.
func (jwt *jwt) IsValid(leeway time.Duration) bool {
	return jwt.claims.IsValid(leeway)
}

// String implements the Stringer interface.
func (jwt *jwt) String() string {
	return jwt.token
}

//--------------------
// PRIVATE HELPERS
//--------------------

// marshallAndEncode marshals the passed value to JSON and
// creates a BASE64 string out of it.
func marshallAndEncode(value interface{}) (string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", errors.Annotate(err, ErrJSONMarshalling, errorMessages)
	}
	encoded := base64.RawURLEncoding.EncodeToString(jsonValue)
	return encoded, nil
}

// decodeAndUnmarshall decodes a BASE64 encoded JSON string and
// unmarshals it into the passed value.
func decodeAndUnmarshall(part string, value interface{}) error {
	decoded, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return errors.Annotate(err, ErrInvalidTokenPart, errorMessages)
	}
	err = json.Unmarshal(decoded, value)
	if err != nil {
		return errors.Annotate(err, ErrJSONUnmarshalling, errorMessages)
	}
	return nil
}

// signAndEncode creates the signature for the data part (header and
// payload) of the token using the passed key and algorithm. The result
// is then encoded to BASE64.
func signAndEncode(data []byte, key Key, algorithm Algorithm) (string, error) {
	sig, err := algorithm.Sign(data, key)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(sig)
	return encoded, nil
}

// decodeAndVerify decodes a BASE64 encoded signature and verifies
// the correct signing of the data part (header and payload) using the
// passed key and algorithm.
func decodeAndVerify(parts []string, key Key, algorithm Algorithm) error {
	data := []byte(parts[0] + "." + parts[1])
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return errors.Annotate(err, ErrInvalidTokenPart, errorMessages)
	}
	return algorithm.Verify(data, sig, key)
}

// EOF
