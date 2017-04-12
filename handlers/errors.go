// Tideland Go REST Server Library - Handlers - Errors
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the handlers package.
const (
	ErrUploadingFile = iota + 1
	ErrDownloadingFile
)

var errorMessages = errors.Messages{
	ErrUploadingFile:   "uploaded file cannot be handled by '%s'",
	ErrDownloadingFile: "file '%s' cannot be downloaded",
}

// EOF
