// Tideland GoREST - JSON Web Token - Unit Tests
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
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

// TestClaimsBool tests the bool operations
// on claims.
func TestClaimsBool(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims bool handling")
	claims := jwt.NewClaims()
	claims.Set("foo", true)
	claims.Set("bar", false)
	claims.Set("baz", "T")
	claims.Set("bingo", "0")
	claims.Set("yadda", "nope")
	foo, ok := claims.GetBool("foo")
	assert.True(foo)
	assert.True(ok)
	bar, ok := claims.GetBool("bar")
	assert.False(bar)
	assert.True(ok)
	baz, ok := claims.GetBool("baz")
	assert.True(baz)
	assert.True(ok)
	bingo, ok := claims.GetBool("bingo")
	assert.False(bingo)
	assert.True(ok)
	yadda, ok := claims.GetBool("yadda")
	assert.False(yadda)
	assert.False(ok)
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

// nestedValue is used as a structured value of a claim.
type nestedValue struct {
	Name  string
	Value int
}

// TestClaimsMarshalledValue tests the marshalling and
// unmarshalling of structures as values.
func TestClaimsMarshalledValue(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims deep value unmarshalling")
	baz := []*nestedValue{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}
	claims := jwt.NewClaims()
	claims.Set("foo", "bar")
	claims.Set("baz", baz)
	// Now marshal and unmarshal the claims.
	jsonValue, err := json.Marshal(claims)
	assert.NotNil(jsonValue)
	assert.Nil(err)
	var unmarshalled jwt.Claims
	err = json.Unmarshal(jsonValue, &unmarshalled)
	assert.Nil(err)
	assert.Length(unmarshalled, 2)
	foo, ok := claims.Get("foo")
	assert.Equal(foo, "bar")
	assert.True(ok)
	var unmarshalledBaz []*nestedValue
	ok, err = claims.GetMarshalled("baz", &unmarshalledBaz)
	assert.True(ok)
	assert.Nil(err)
	assert.Length(unmarshalledBaz, 3)
	assert.Equal(unmarshalledBaz[0].Name, "one")
	assert.Equal(unmarshalledBaz[2].Value, 3)
}

// TestClaimsAudience checks the setting, getting, and
// deleting of the audience claim.
func TestClaimsAudience(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"aud\"")
	audience := []string{"foo", "bar", "baz"}
	claims := jwt.NewClaims()
	aud, ok := claims.Audience()
	assert.False(ok)
	none := claims.SetAudience(audience...)
	assert.Length(none, 0)
	aud, ok = claims.Audience()
	assert.Equal(aud, audience)
	assert.True(ok)
	old := claims.DeleteAudience()
	assert.Equal(old, aud)
	_, ok = claims.Audience()
	assert.False(ok)
}

// TestClaimsExpiration checks the setting, getting, and
// deleting of the expiration claim.
func TestClaimsExpiration(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"exp\"")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	claims := jwt.NewClaims()
	exp, ok := claims.Expiration()
	assert.False(ok)
	none := claims.SetExpiration(goLaunch)
	assert.Equal(none, time.Time{})
	exp, ok = claims.Expiration()
	assert.Equal(exp.Unix(), goLaunch.Unix())
	assert.True(ok)
	old := claims.DeleteExpiration()
	assert.Equal(old.Unix(), exp.Unix())
	exp, ok = claims.Expiration()
	assert.False(ok)
}

// TestClaimsIdentifier checks the setting, getting, and
// deleting of the identifier claim.
func TestClaimsIdentifier(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"jti\"")
	identifier := "foo"
	claims := jwt.NewClaims()
	jti, ok := claims.Identifier()
	assert.False(ok)
	none := claims.SetIdentifier(identifier)
	assert.Equal(none, "")
	jti, ok = claims.Identifier()
	assert.Equal(jti, identifier)
	assert.True(ok)
	old := claims.DeleteIdentifier()
	assert.Equal(old, jti)
	_, ok = claims.Identifier()
	assert.False(ok)
}

// TestClaimsIssuedAt checks the setting, getting, and
// deleting of the issued at claim.
func TestClaimsIssuedAt(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"iat\"")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	claims := jwt.NewClaims()
	iat, ok := claims.IssuedAt()
	assert.False(ok)
	none := claims.SetIssuedAt(goLaunch)
	assert.Equal(none, time.Time{})
	iat, ok = claims.IssuedAt()
	assert.Equal(iat.Unix(), goLaunch.Unix())
	assert.True(ok)
	old := claims.DeleteIssuedAt()
	assert.Equal(old.Unix(), iat.Unix())
	iat, ok = claims.IssuedAt()
	assert.False(ok)
}

// TestClaimsIssuer checks the setting, getting, and
// deleting of the issuer claim.
func TestClaimsIssuer(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"iss\"")
	issuer := "foo"
	claims := jwt.NewClaims()
	iss, ok := claims.Issuer()
	assert.False(ok)
	none := claims.SetIssuer(issuer)
	assert.Equal(none, "")
	iss, ok = claims.Issuer()
	assert.Equal(iss, issuer)
	assert.True(ok)
	old := claims.DeleteIssuer()
	assert.Equal(old, iss)
	_, ok = claims.Issuer()
	assert.False(ok)
}

// TestClaimsNotBefore checks the setting, getting, and
// deleting of the not before claim.
func TestClaimsNotBefore(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"nbf\"")
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	claims := jwt.NewClaims()
	nbf, ok := claims.NotBefore()
	assert.False(ok)
	none := claims.SetNotBefore(goLaunch)
	assert.Equal(none, time.Time{})
	nbf, ok = claims.NotBefore()
	assert.Equal(nbf.Unix(), goLaunch.Unix())
	assert.True(ok)
	old := claims.DeleteNotBefore()
	assert.Equal(old.Unix(), nbf.Unix())
	nbf, ok = claims.NotBefore()
	assert.False(ok)
}

// TestClaimsSubject checks the setting, getting, and
// deleting of the subject claim.
func TestClaimsSubject(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claim \"sub\"")
	subject := "foo"
	claims := jwt.NewClaims()
	sub, ok := claims.Subject()
	assert.False(ok)
	none := claims.SetSubject(subject)
	assert.Equal(none, "")
	sub, ok = claims.Subject()
	assert.Equal(sub, subject)
	assert.True(ok)
	old := claims.DeleteSubject()
	assert.Equal(old, sub)
	_, ok = claims.Subject()
	assert.False(ok)
}

// TestClaimsValidity checks the validation of the not before
// and the expiring time.
func TestClaimsValidity(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing claims validity")
	// Fresh claims.
	now := time.Now()
	leeway := time.Minute
	claims := jwt.NewClaims()
	valid := claims.IsAlreadyValid(leeway)
	assert.True(valid)
	valid = claims.IsStillValid(leeway)
	assert.True(valid)
	valid = claims.IsValid(leeway)
	assert.True(valid)
	// Set times.
	nbf := now.Add(-time.Hour)
	exp := now.Add(time.Hour)
	claims.SetNotBefore(nbf)
	valid = claims.IsAlreadyValid(leeway)
	assert.True(valid)
	claims.SetExpiration(exp)
	valid = claims.IsStillValid(leeway)
	assert.True(valid)
	valid = claims.IsValid(leeway)
	assert.True(valid)
	// Invalid claims.
	nbf = now.Add(time.Hour)
	exp = now.Add(-time.Hour)
	claims.SetNotBefore(nbf)
	claims.DeleteExpiration()
	valid = claims.IsAlreadyValid(leeway)
	assert.False(valid)
	valid = claims.IsValid(leeway)
	assert.False(valid)
	claims.DeleteNotBefore()
	claims.SetExpiration(exp)
	valid = claims.IsStillValid(leeway)
	assert.False(valid)
	valid = claims.IsValid(leeway)
	assert.False(valid)
	claims.SetNotBefore(nbf)
	valid = claims.IsValid(leeway)
	assert.False(valid)
}

//--------------------
// HELPERS
//--------------------

// initClaims creates test claims.
func initClaims() jwt.Claims {
	claims := jwt.NewClaims()
	claims.SetSubject("1234567890")
	claims.Set("name", "John Doe")
	claims.Set("admin", true)
	return claims
}

// testClaims checks the passed claims.
func testClaims(assert audit.Assertion, claims jwt.Claims) {
	sub, ok := claims.Subject()
	assert.True(ok)
	assert.Equal(sub, "1234567890")
	name, ok := claims.GetString("name")
	assert.True(ok)
	assert.Equal(name, "John Doe")
	admin, ok := claims.GetBool("admin")
	assert.True(ok)
	assert.True(admin)
}

// EOF
