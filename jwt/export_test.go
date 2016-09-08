// Tideland Go REST Server Library - JSON Web Token - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt

//--------------------
// IMPORTS
//--------------------

import (
	"time"
)

//--------------------
// HELPERS
//--------------------

// SetCleanupInterval allows to configure a shorter cleanup
// interval for tests. The returned function resets the
// current interval when called, e.g. with defer.
func SetCleanupInterval(interval time.Duration) func() {
	currentInterval := cleanupInterval
	reset := func() {
		cleanupInterval = currentInterval
	}
	cleanupInterval = interval
	return reset
}

// EOF
