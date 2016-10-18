// Tideland Go REST Server Library - Handlers
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go REST Server Library handlers package defines
// some initial resource handlers to integrate into own solutions.
package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/version"
)

//--------------------
// VERSION
//--------------------

// Version returns the version of the handlers package.
func Version() version.Version {
	return version.New(2, 5, 0)
}

// EOF
