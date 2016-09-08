// Tideland Go REST Server Library - JSON Web Token - Cache
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
	"sync"
	"time"
)

//--------------------
// CACHE
//--------------------

// Cache provides a caching for tokens so that these
// don't have to be decoded or verified multiple times.
type Cache interface {
	// Get tries to retrieve a token from the cache.
	Get(token string) (JWT, bool)

	// Put adds a token to the cache.
	Put(jwt JWT)
}

// cacheEntry manages a token and its access time.
type cacheEntry struct {
	jwt      JWT
	accessed time.Time
}

// cache implements Cache.
type cache struct {
	mutex   sync.Mutex
	cleanup time.Duration
	leeway  time.Duration
	entries map[string]*cacheEntry
}

// NewCache creates a new JWT caching. It takes two
// durations. The first one is the time a token hasn't
// been used anymore before it is cleaned up. The second
// one is the leeway taken for token time validations.
func NewCache(cleanup, leeway time.Duration) Cache {
	c := &cache{
		cleanup: cleanup,
		leeway:  leeway,
		entries: map[string]*cacheEntry{},
	}
	// TODO Start cleanup goroutine.
	return c
}

// Get implements the Cache interface.
func (c *cache) Get(token string) (JWT, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[token]
	if !ok {
		return nil, false
	}
	if entry.jwt.IsValid(c.leeway) {
		entry.accessed = time.Now()
		return entry.jwt, true
	}
	// Remove invalid token.
	delete(c.entries, token)
	return nil, false
}

// Put implements the Cache interface.
func (c *cache) Put(jwt JWT) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if jwt.IsValid(c.leeway) {
		c.entries[jwt.String()] = &cacheEntry{jwt, time.Now()}
	}
}

// EOF
