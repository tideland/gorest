// Tideland Go REST Server Library - JSON Web Token - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt_test

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
)

//--------------------
// TESTS
//--------------------

// TestClaimsBasic tests the low level operations
// on claims.
func TestClaimsBasic(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	// First with uninitialised claims.
	var claims jwt.Claims
	ok := claims.Contains("foo")
	assert.False(ok)
	nothing, ok := claims.Get("foo")
	assert.Nil(nothing)
	assert.False(ok)
	old := claims.Set("foo", "bar")
	assert.Nil(old)
	old = claims.Delete("foo")
	assert.Nil(old)
	// Now initialise it.
	claims = jwt.NewClaims()
	ok = claims.Contains("foo")
	assert.False(ok)
	mothing, ok = claims.Get("foo")
	assert.Nil(nothing)
	assert.False(ok)
	old = claims.Set("foo", "bar")
	assert.Nil(old)
	ok = claims.Contains("foo")
	assert.True(ok)
	foo, ok := claims.Get("foo")
	assert.Equal(foo, "bar")
	assert.True(ok)
	old = claims.Set("foo", "yadda")
	assert.Equal(old, "bar")
	// Finally delete it.
	old = claims.Delete("foo")
	assert.Equal(old, "yadda")
	old = claims.Delete("foo")
	assert.Nil(old)
	ok = claims.Contains("foo")
	assert.True(false)
}

// EOF
