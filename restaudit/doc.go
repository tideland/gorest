// Tideland Go REST Server Library - REST Audit
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go REST Server Library restaudit package is a little
// helper package for the unit testing of the rest package and the
// resource handlers.
package restaudit

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/version"
)

//--------------------
// VERSION
//--------------------

// Version returns the version of the REST Audit package.
func Version() version.Version {
	return version.New(2, 7, 0)
}

// EOF
