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

	"github.com/tideland/golib/loop"
)

//--------------------
// CACHE
//--------------------

// cleanupInterval defines the timespan between
// two cleanup runs.
var cleanupInterval = 5 * time.Minute

// Cache provides a caching for tokens so that these
// don't have to be decoded or verified multiple times.
type Cache interface {
	// Get tries to retrieve a token from the cache.
	Get(token string) (JWT, bool)

	// Put adds a token to the cache.
	Put(jwt JWT)

	// Stop tells the cache to end working.
	Stop() error
}

// cacheEntry manages a token and its access time.
type cacheEntry struct {
	jwt      JWT
	accessed time.Time
}

// cache implements Cache.
type cache struct {
	mutex   sync.Mutex
	ttl     time.Duration
	leeway  time.Duration
	entries map[string]*cacheEntry
	loop    loop.Loop
}

// NewCache creates a new JWT caching. It takes two
// durations. The first one is the time a token hasn't
// been used anymore before it is cleaned up. The second
// one is the leeway taken for token time validations.
func NewCache(ttl, leeway time.Duration) Cache {
	c := &cache{
		ttl:     ttl,
		leeway:  leeway,
		entries: map[string]*cacheEntry{},
	}
	c.loop = loop.Go(c.backendLoop, "jwt", "cache")
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

// Stop implements the Cache interface.
func (c *cache) Stop() error {
	return c.loop.Stop()
}

// backendLoop runs a cleaning session every five minutes.
func (c *cache) backendLoop(l loop.Loop) error {
	defer func() {
		// Some cleanup after stop or error.
		c.ttl = 0
		c.leeway = 0
		c.entries = nil
	}()
	ticker := time.NewTicker(cleanupInterval)
	for {
		select {
		case <-l.ShallStop():
			return nil
		case <-ticker.C:
			c.cleanup()
		}
	}
}

// cleanup checks for invalid or unused tokens.
func (c *cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	valids := map[string]*cacheEntry{}
	now := time.Now()
	for token, entry := range c.entries {
		if entry.jwt.IsValid(c.leeway) {
			if entry.accessed.Add(c.ttl).After(now) {
				// Everything fine.
				valids[token] = entry
			}
		}
	}
	c.entries = valids
}

// EOF
