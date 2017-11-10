// Tideland GoREST - REST - Environment
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest

//--------------------
// IMPORTS
//--------------------

import (
	"context"

	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/stringex"
)

//--------------------
// ENVIRONMENT
//--------------------

// Environment describes the environment of a RESTful application.
type Environment interface {
	// Context returns the context of the environment.
	Context() context.Context

	// Basepath returns the configured basepath.
	Basepath() string

	// DefaultDomain returns the configured default domain.
	DefaultDomain() string

	// DefaultResource returns the configured default resource.
	DefaultResource() string

	// TemplatesCache returns the template cache.
	TemplatesCache() TemplatesCache
}

// environment implements the Environment interface.
type environment struct {
	ctx             context.Context
	basepath        string
	baseparts       []string
	basepartsLen    int
	defaultDomain   string
	defaultResource string
	templatesCache  TemplatesCache
}

// newEnvironment crerates an environment using the
// passed context and configuration.
func newEnvironment(ctx context.Context, cfg etc.Etc) *environment {
	env := &environment{
		basepath:        "/",
		baseparts:       []string{},
		defaultDomain:   "default",
		defaultResource: "default",
		templatesCache:  newTemplatesCache(),
	}
	// Check configuration.
	if cfg != nil {
		env.basepath = cfg.ValueAsString("basepath", env.basepath)
		env.defaultDomain = cfg.ValueAsString("default-domain", env.defaultDomain)
		env.defaultResource = cfg.ValueAsString("default-resource", env.defaultResource)
	}
	// Check basepath and remove empty parts.
	env.baseparts = stringex.SplitMap(env.basepath, "/", func(p string) (string, bool) {
		if p == "" {
			return "", false
		}
		return p, true
	})
	env.basepartsLen = len(env.baseparts)
	// Set context.
	if ctx == nil {
		ctx = context.Background()
	}
	env.ctx = newEnvironmentContext(ctx, env)
	return env
}

// Context implements the Environment interface.
func (env *environment) Context() context.Context {
	return env.ctx
}

// Basepath implements the Environment interface.
func (env *environment) Basepath() string {
	return env.basepath
}

// DefaultDomain implements the Environment interface.
func (env *environment) DefaultDomain() string {
	return env.defaultDomain
}

// DefaultResource implements the Environment interface.
func (env *environment) DefaultResource() string {
	return env.defaultResource
}

// TemplatesCache implements the Environment interface.
func (env *environment) TemplatesCache() TemplatesCache {
	return env.templatesCache
}

// EOF
