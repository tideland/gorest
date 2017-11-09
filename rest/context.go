// Tideland GoREST - REST - Context
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
)

//--------------------
// CONST
//--------------------

// contextKey is used to address data inside a context.
type contextKey int

const (
	// envKey addresses the environment inside the context.
	envKey contextKey = 0

	// jobKey addresses the job inside the context.
	jobKey contextKey = 1
)

//--------------------
// CONTEXT
//--------------------

// newEnvironmentContext creates a context based on the passed one
// and containing the passed environment.
func newEnvironmentContext(ctx context.Context, env Environment) context.Context {
	return context.WithValue(ctx, envKey, env)
}

// newJobContext creates a context based on the passed one
// and containing the passed job.
func newJobContext(ctx context.Context, job Job) context.Context {
	return context.WithValue(ctx, jobKey, job)
}

// EnvironmentFromContext retrieves the environment out of a context.
func EnvironmentFromContext(ctx context.Context) (Environment, bool) {
	env, ok := ctx.Value(envKey).(Environment)
	return env, ok
}

// JobFromContext retrieves the job out of a context.
func JobFromContext(ctx context.Context) (Job, bool) {
	job, ok := ctx.Value(jobKey).(Job)
	return job, ok
}

// EOF
