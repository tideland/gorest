// Tideland Go REST Server Library - Request
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The request package provides a simple way to handle cross-server
// requests in the Tideland REST ecosystem.
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

// Version returns the version of the REST package.
func Version() version.Version {
	return version.New(2, 6, 0, "alpha", "2016-10-30")
}

// EOF
