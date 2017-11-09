// Tideland GoREST - REST - Job
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
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/tideland/golib/logger"
	"github.com/tideland/golib/version"
)

//--------------------
// JOB
//--------------------

// Job encapsulates all the needed information for handling
// a request.
type Job interface {
	// Return the Job as string.
	fmt.Stringer

	// Environment returns the server environment.
	Environment() Environment

	// Request returns the used Go HTTP request.
	Request() *http.Request

	// ResponseWriter returns the used Go HTTP response writer.
	ResponseWriter() http.ResponseWriter

	// Domain returns the requests domain. It's deprecated,
	// use Job.Path().Domain() instead.
	Domain() string

	// Resource returns the requests resource. It's deprecated,
	// use Job.Path().Resource() instead.
	Resource() string

	// ResourceID return the requests resource ID. It's deprecated,
	// use Job.Path().JoinedResourceID() instead.
	ResourceID() string

	// Path returns access to the request path inside the URL.
	Path() Path

	// Context returns a job context also containing the
	// job itself.
	Context() context.Context

	// EnhanceContext allows to enhance the job context
	// values, a deadline, a timeout, or a cancel. So
	// e.g. a first handler in a handler queue can
	// store authentication information in the context
	// and a following handler can use it (see the
	// JWTAuthorizationHandler).
	EnhanceContext(func(ctx context.Context) context.Context)

	// Version returns the requested API version for this job. If none
	// is set the version 1.0.0 will be returned as default. It will
	// be retrieved aut of the header Version.
	Version() version.Version

	// SetVersion allows to set an API version for the response. If
	// none is set the version 1.0.0 will be set as default. It will
	// be set in the header Version.
	SetVersion(v version.Version)

	// AcceptsContentType checks if the requestor accepts a given content type.
	AcceptsContentType(contentType string) bool

	// HasContentType checks if the sent content has the given content type.
	HasContentType(contentType string) bool

	// Languages returns the accepted language with the quality values.
	Languages() Languages

	// InternalPath builds an internal path out of the passed parts.
	InternalPath(domain, resource, resourceID string, query ...KeyValue) string

	// Redirect to a domain, resource and resource ID (optional).
	Redirect(domain, resource, resourceID string)

	// Renderer returns a template renderer.
	Renderer() Renderer

	// GOB returns a GOB formatter.
	GOB() Formatter

	// JSON returns a JSON formatter.
	JSON(html bool) Formatter

	// XML returns a XML formatter.
	XML() Formatter

	// Query returns a convenient access to query values.
	Query() Values

	// Form returns a convenient access to form values.
	Form() Values
}

// job implements the Job interface.
type job struct {
	environment    *environment
	ctx            context.Context
	request        *http.Request
	responseWriter http.ResponseWriter
	version        version.Version
	path           Path
}

// newJob parses the URL and returns the prepared job.
func newJob(env *environment, r *http.Request, rw http.ResponseWriter) Job {
	// Init the job.
	j := &job{
		environment:    env,
		request:        r,
		responseWriter: rw,
		path:           newPath(env, r),
	}
	// Retrieve the requested version of the API.
	vsnstr := j.request.Header.Get("Version")
	if vsnstr == "" {
		j.version = version.New(1, 0, 0)
	} else {
		vsn, err := version.Parse(vsnstr)
		if err != nil {
			logger.Errorf("invalid request version: %v", err)
			j.version = version.New(1, 0, 0)
		} else {
			j.version = vsn
		}
	}
	return j
}

// String is defined on the Stringer interface.
func (j *job) String() string {
	path := j.createPath(j.Domain(), j.Resource(), j.ResourceID())
	return fmt.Sprintf("%s %s", j.request.Method, path)
}

// Environment implements the Job interface.
func (j *job) Environment() Environment {
	return j.environment
}

// Request implements the Job interface.
func (j *job) Request() *http.Request {
	return j.request
}

// ResponseWriter implements the Job interface.
func (j *job) ResponseWriter() http.ResponseWriter {
	return j.responseWriter
}

// Domain implements the Job interface.
func (j *job) Domain() string {
	return j.path.Domain()
}

// Resource implements the Job interface.
func (j *job) Resource() string {
	return j.path.Resource()
}

// ResourceID implements the Job interface.
func (j *job) ResourceID() string {
	return j.path.JoinedResourceID()
}

// Path implements the Job interface.
func (j *job) Path() Path {
	return j.path
}

// Context implements the Job interface.
func (j *job) Context() context.Context {
	// Lazy init.
	if j.ctx == nil {
		j.ctx = newJobContext(j.environment.ctx, j)
	}
	return j.ctx
}

// EnhanceContext implements the Job interface.
func (j *job) EnhanceContext(f func(ctx context.Context) context.Context) {
	ctx := j.Context()
	j.ctx = f(ctx)
}

// Version implements the Job interface.
func (j *job) Version() version.Version {
	return j.version
}

// SerVersion implements the Job interface.
func (j *job) SetVersion(vsn version.Version) {
	if vsn != nil {
		j.version = vsn
	}
}

// AcceptsContentType implements the Job interface.
func (j *job) AcceptsContentType(contentType string) bool {
	return strings.Contains(j.request.Header.Get("Accept"), contentType)
}

// HasContentType implements the Job interface.
func (j *job) HasContentType(contentType string) bool {
	return strings.Contains(j.request.Header.Get("Content-Type"), contentType)
}

// Languages implements the Job interface.
func (j *job) Languages() Languages {
	accept := j.request.Header.Get("Accept-Language")
	languages := Languages{}
	for _, part := range strings.Split(accept, ",") {
		lv := strings.Split(part, ";")
		if len(lv) == 1 {
			languages = append(languages, Language{lv[0], 1.0})
		} else {
			value, err := strconv.ParseFloat(lv[1], 64)
			if err != nil {
				value = 0.0
			}
			languages = append(languages, Language{lv[0], value})
		}
	}
	languages = sort.Reverse(languages).(Languages)
	return languages
}

// createPath creates a path out of the major URL parts.
func (j *job) createPath(domain, resource, resourceID string) string {
	parts := append(j.environment.baseparts, domain, resource)
	if resourceID != "" {
		parts = append(parts, resourceID)
	}
	path := strings.Join(parts, "/")
	return "/" + path
}

// InternalPath implements the Job interface.
func (j *job) InternalPath(domain, resource, resourceID string, query ...KeyValue) string {
	path := j.createPath(domain, resource, resourceID)
	if len(query) > 0 {
		path += "?" + KeyValues(query).String()
	}
	return path
}

// Redirect implements the Job interface.
func (j *job) Redirect(domain, resource, resourceID string) {
	path := j.createPath(domain, resource, resourceID)
	http.Redirect(j.responseWriter, j.request, path, http.StatusTemporaryRedirect)
}

// Renderer implements the Job interface.
func (j *job) Renderer() Renderer {
	return &renderer{j.responseWriter, j.environment.templatesCache}
}

// GOB implements the Job interface.
func (j *job) GOB() Formatter {
	return &gobFormatter{j}
}

// JSON implements the Job interface.
func (j *job) JSON(html bool) Formatter {
	return &jsonFormatter{j, html}
}

// XML implements the Job interface.
func (j *job) XML() Formatter {
	return &xmlFormatter{j}
}

// Query implements the Job interface.
func (j *job) Query() Values {
	return &values{j.request.URL.Query()}
}

// Form implements the Job interface.
func (j *job) Form() Values {
	return &values{j.request.PostForm}
}

// EOF
