// Tideland Go REST Server Library - REST
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go REST Server Library provides the package rest for the
// implementation of servers with a RESTful API. The business has to
// be implemented in types fullfilling the ResourceHandler interface.
// This basic interface only allows the initialization of the handler.
// More interesting are the other interfaces like GetResourceHandler
// which defines the Get() method for the HTTP request method GET.
// Others are for PUT, POST, HEAD, PATCH, DELETE, and OPTIONS. Their
// according methods get a Job as argument. It provides convenient
// helpers for the processing of the job.
//
// The processing methods return two values: a boolean and an error.
// The latter is pretty clear, it signals a job processing error. The
// boolean is more interesting. Registering a handler is based on a
// domain and a resource. The URL
//
// /<DOMAIN>/<RESOURCE>
//
// leads to a handler, or even better, to a list of handlers. All
// are used as long as the returned boolean value is true. E.g. the
// first handler can check the authentication, the second one the
// authorization, and the third one does the business. Additionally
// the URL
//
// /<DOMAIN>/<RESOURCE>/<ID>
//
// provides the resource identifier via Job.ResourceID().
package rest

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
	return version.New(2, 0, 0, "alpha", "2016-08-21")
}

// EOF
