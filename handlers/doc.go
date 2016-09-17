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

// PackageVersion returns the version of the version package.
func PackageVersion() version.Version {
	return version.New(2, 0, 0, "beta", "2016-09-17")
}

// EOF
