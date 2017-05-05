// Tideland Go REST Server Library - REST - Path
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
	"github.com/tideland/golib/stringex"
)

//--------------------
// CONSTANTS
//--------------------

// Path indexes for the different parts.
const (
	PathDomain = 0
	PathResource = 1
	PathResourceID = 2
)

//--------------------
// PATH
//--------------------

// Path provides access to the parts of a
// request path interesting for handling a
// job. 
type Path interface {
	// Length returns the number of parts of the path.
	Length() int
	
	// Part returns the parts of the URL path based on the
	// index or an empty string.
	Part(index int) string
	
	// Domain returns the requests domain.
	Domain() string

	// Resource returns the requests resource.
	Resource() string

	// ResourceID return the requests resource ID.
	ResourceID() string
}

// path implements Path.
type path struct {
	path []string
}

func newPath(url ) *path {
}

// EOF