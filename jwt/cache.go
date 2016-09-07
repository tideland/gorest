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
	mutex   sync.RWMutex
	entries map[string]*cacheEntry
}

// Get implements the Cache interface.
func (c *cache) Get(token string) (JWT, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.entries[token]
	if !ok {
		return nil, false
	}
	// TODO Check claims and their validity.
	return entry.jwt, true
}

// Put implements the Cache interface.
func (c *cache) Put(jwt JWT) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.entries, jwt.Token())
}

// EOF
