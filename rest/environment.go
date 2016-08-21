// Tideland Go REST Server Library - REST - Environment
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rest

//--------------------
// IMPORTS
//--------------------

import ()

//--------------------
// ENVIRONMENT
//--------------------

type Environment interface {
	// BasePath returns the configured base path.
	BasePath() string

	// DefaultDomain returns the configured default domain.
	DefaultDomain() string

	// DefaultResource returns the configured default resource.
	DefaultResource() string

	// Templates returns the template cache.
	Templates() TemplatesCache
}

// environment implements the Environment interface.
type environment struct {
	basePath        string
	defaultDomain   string
	defaultResource string
	templates       TemplatesCache
}

// Option defines a function for setting an option.
type Option func(env Environment)

// newEnvironment crerates the default environment and
// checks the passed options.
func newEnvironment(options ...Option) Environment {
	env := &environment{
		basePath:        "/",
		defaultDomain:   "default",
		defaultResource: "default",
		templates:       NewTemplatesCache(),
	}
	for _, option := range options {
		option(env)
	}
	return env
}

// BasePath is specified on the Environment interface.
func (env *environment) BasePath() string {
	return env.basePath
}

// DefaultDomain is specified on the Environment interface.
func (env *environment) DefaultDomain() string {
	return env.defaultDomain
}

// DefaultResource is specified on the Environment interface.
func (env *environment) DefaultResource() string {
	return env.defaultResource
}

// Templates is specified on the Environment interface.
func (env *environment) Templates() TemplatesCache {
	return env.templates
}

//--------------------
// OPTIONS
//--------------------

// BasePath sets the path thats used as prefix before
// domain and resource.
func BasePath(basePath string) Option {
	return func(env Environment) {
		if basePath == "" {
			basePath = "/"
		}
		if basePath[len(basePath)-1] != '/' {
			basePath += "/"
		}
		envImpl, ok := env.(*environment)
		if ok {
			envImpl.basePath = basePath
		}
	}
}

// DefaultDomain sets the default domain.
func DefaultDomain(defaultDomain string) Option {
	return func(env Environment) {
		if defaultDomain == "" {
			defaultDomain = "default"
		}
		envImpl, ok := env.(*environment)
		if ok {
			envImpl.defaultDomain = defaultDomain
		}
	}
}

// DefaultResource sets the default resource.
func DefaultResource(defaultResource string) Option {
	return func(env Environment) {
		if defaultResource == "" {
			defaultResource = "default"
		}
		envImpl, ok := env.(*environment)
		if ok {
			envImpl.defaultResource = defaultResource
		}
	}
}

// Templates sets the templates cache.
func Templates(templates TemplatesCache) Option {
	return func(env Environment) {
		if templates == nil {
			templates = NewTemplatesCache()
		}
		envImpl, ok := env.(*environment)
		if ok {
			envImpl.templates = templates
		}
	}
}

// EOF
