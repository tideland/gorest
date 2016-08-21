// Tideland Go REST Server Library - REST - Context
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
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

// jobKey addresses the worksheet inside the context.
const jobKey contextKey = 0

//--------------------
// CONTEXT
//--------------------

// newContext creates a context containing the passed job.
func newContext(job Job) context.Context {
	return context.WithValue(context.Background(), jobKey, job)
}

// FromContext retrieves the job out of a context.
func FromContext(ctx context.Context) (Job, bool) {
	job, ok := ctx.Value(jobKey).(Job)
	return job, ok
}

// EOF
