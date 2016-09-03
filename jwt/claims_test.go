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
	"encoding/json"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
)

//--------------------
// TESTS
//--------------------

// TestClaimsMarshalling tests the marshalling of Claims
// to JSON and back.
func TestClaimsMarshalling(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims marshalling")
	// First with uninitialised or empty claims.
	var claims jwt.Claims
	jsonValue, err := json.Marshal(claims)
	assert.Equal(string(jsonValue), "{}")
	assert.Nil(err)
	claims = jwt.NewClaims()
	jsonValue, err = json.Marshal(claims)
	assert.Equal(string(jsonValue), "{}")
	assert.Nil(err)
	// Now fill it.
	claims.Set("foo", "yadda")
	claims.Set("bar", 12345)
	assert.Length(claims, 2)
	jsonValue, err = json.Marshal(claims)
	assert.NotNil(jsonValue)
	assert.Nil(err)
	var unmarshalled jwt.Claims
	err = json.Unmarshal(jsonValue, &unmarshalled)
	assert.Nil(err)
	assert.Length(unmarshalled, 2)
	foo, ok := claims.Get("foo")
	assert.Equal(foo, "yadda")
	assert.True(ok)
	bar, ok := claims.GetInt("bar")
	assert.Equal(bar, 12345)
	assert.True(ok)
}

// TestClaimsBasic tests the low level operations
// on claims.
func TestClaimsBasic(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims basic functions handling")
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
	nothing, ok = claims.Get("foo")
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
	assert.False(ok)
}

// TestClaimsString tests the string operations
// on claims.
func TestClaimsString(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims string handling")
	claims := jwt.NewClaims()
	nothing := claims.Set("foo", "bar")
	assert.Nil(nothing)
	var foo string
	foo, ok := claims.GetString("foo")
	assert.Equal(foo, "bar")
	assert.True(ok)
	claims.Set("foo", 4711)
	foo, ok = claims.GetString("foo")
	assert.Equal(foo, "4711")
	assert.True(ok)
}

// TestClaimsInt tests the int operations
// on claims.
func TestClaimsInt(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims int handling")
	claims := jwt.NewClaims()
	claims.Set("foo", 4711)
	claims.Set("bar", "4712")
	claims.Set("baz", 4713.0)
	claims.Set("yadda", "nope")
	foo, ok := claims.GetInt("foo")
	assert.Equal(foo, 4711)
	assert.True(ok)
	bar, ok := claims.GetInt("bar")
	assert.Equal(bar, 4712)
	assert.True(ok)
	baz, ok := claims.GetInt("baz")
	assert.Equal(baz, 4713)
	assert.True(ok)
	yadda, ok := claims.GetInt("yadda")
	assert.Equal(yadda, 0)
	assert.False(ok)
}

// TestClaimsFloat64 tests the float64 operations
// on claims.
func TestClaimsFloat64(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims float64 handling")
	claims := jwt.NewClaims()
	claims.Set("foo", 4711)
	claims.Set("bar", "4712")
	claims.Set("baz", 4713.0)
	claims.Set("yadda", "nope")
	foo, ok := claims.GetFloat64("foo")
	assert.Equal(foo, 4711.0)
	assert.True(ok)
	bar, ok := claims.GetFloat64("bar")
	assert.Equal(bar, 4712.0)
	assert.True(ok)
	baz, ok := claims.GetFloat64("baz")
	assert.Equal(baz, 4713.0)
	assert.True(ok)
	yadda, ok := claims.GetFloat64("yadda")
	assert.Equal(yadda, 0.0)
	assert.False(ok)
}

// TestClaimsTime tests the time operations
// on claims.
func TestClaimsTime(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims time handling")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	claims := jwt.NewClaims()
	claims.SetTime("foo", goLaunch)
	claims.Set("bar", goLaunch.Unix())
	claims.Set("baz", goLaunch.Format(time.RFC3339))
	claims.Set("yadda", "nope")
	foo, ok := claims.GetTime("foo")
	assert.Equal(foo.Unix(), goLaunch.Unix())
	assert.True(ok)
	bar, ok := claims.GetTime("bar")
	assert.Equal(bar.Unix(), goLaunch.Unix())
	assert.True(ok)
	baz, ok := claims.GetTime("baz")
	assert.Equal(baz.Unix(), goLaunch.Unix())
	assert.True(ok)
	yadda, ok := claims.GetTime("yadda")
	assert.Equal(yadda, time.Time{})
	assert.False(ok)
}

// EOF
