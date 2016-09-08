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
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
)

//--------------------
// TESTS
//--------------------

// TestCachePutGet tests the putting and getting of tokens
// to the cache.
func TestCachePutGet(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing cache put and get")
	cache := jwt.NewCache(time.Minute, time.Minute)
	key := []byte("secret")
	claims := init.Claims()
	jwtIn, err := jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	cache.Put(jwt)
	token := jwtIn.String()
	jwtOut, ok := cache.Get(token)
	assert.True(ok)
	assert.Equal(jwtIn, jwtOut)
	jwtOut, ok = cache.Get("is.not.there")
	assert.False(ok)
	assert.Nil(jwtOut)
	err := cache.Stop()
	assert.Nil(err)
}

// EOF