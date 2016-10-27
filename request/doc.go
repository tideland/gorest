// Tideland Go REST Server Library - REST Request
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go REST Server Library request provides simpler
// requests to handlers of the Tideland Go REST Server Library world.
package request

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
	return version.New(2, 5, 2)
}

// EOF
