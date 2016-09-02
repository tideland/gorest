// Tideland Go REST Server Library - JSON Web Token
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go REST Server Library jwt provides the generation,
// verification, and analyzing of JSON Web Tokens.
package jwt

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
	return version.New(2, 0, 0, "beta", "2016-09-02")
}

// EOF
