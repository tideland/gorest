// Tideland Go REST Server Library - REST - Job
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
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

	// Domain returns the requests domain.
	Domain() string

	// Resource returns the requests resource.
	Resource() string

	// ResourceID return the requests resource ID.
	ResourceID() string

	// Context returns a context containing the job.
	Context() context.Context

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

	// RenderTemplate renders a template with the passed data.
	RenderTemplate(templateID string, data interface{}) error

	// GOB returns a GOB formatter.
	GOB() Formatter

	// JSON returns a JSON formatter.
	JSON(html bool) Formatter

	// XML returns a XML formatter.
	XML() Formatter
}

// job implements the Job interface.
type job struct {
	environment    Environment
	ctx            context.Context
	request        *http.Request
	responseWriter http.ResponseWriter
	domain         string
	resource       string
	resourceID     string
}

// newJob parses the URL and returns the prepared job.
func newJob(env Environment, r *http.Request, rw http.ResponseWriter) Job {
	// Init the job.
	j := &job{
		environment:    env,
		request:        r,
		responseWriter: rw,
	}
	// Split path for REST identifiers.
	parts := strings.Split(r.URL.Path[len(env.BasePath()):], "/")
	switch len(parts) {
	case 3:
		j.resourceID = parts[2]
		j.resource = parts[1]
		j.domain = parts[0]
	case 2:
		j.resource = parts[1]
		j.domain = parts[0]
	case 1:
		j.resource = j.environment.DefaultResource()
		j.domain = parts[0]
	case 0:
		j.resource = j.environment.DefaultResource()
		j.domain = j.environment.DefaultDomain()
	default:
		j.resourceID = strings.Join(parts[2:], "/")
		j.resource = parts[1]
		j.domain = parts[0]
	}
	return j
}

// String is defined on the Stringer interface.
func (j *job) String() string {
	if j.resourceID == "" {
		return fmt.Sprintf("%s /%s/%s", j.request.Method, j.domain, j.resource)
	}
	return fmt.Sprintf("%s /%s/%s/%s", j.request.Method, j.domain, j.resource, j.resourceID)
}

// Environment is specified on the Job interface.
func (j *job) Environment() Environment {
	return j.environment
}

// Request is specified on the Job interface.
func (j *job) Request() *http.Request {
	return j.request
}

// ResponseWriter is specified on the Job interface.
func (j *job) ResponseWriter() http.ResponseWriter {
	return j.responseWriter
}

// Domain is specified on the Job interface.
func (j *job) Domain() string {
	return j.domain
}

// Resource is specified on the Job interface.
func (j *job) Resource() string {
	return j.resource
}

// ResourceID is specified on the Job interface.
func (j *job) ResourceID() string {
	return j.resourceID
}

// Context is specified on the Job interface.
func (j *job) Context() context.Context {
	// Lazy init.
	if j.ctx == nil {
		j.ctx = newContext(j)
	}
	return j.ctx
}

// AcceptsContentType is specified on the Job interface.
func (j *job) AcceptsContentType(contentType string) bool {
	return strings.Contains(j.request.Header.Get("Accept"), contentType)
}

// HasContentType is specified on the Job interface.
func (j *job) HasContentType(contentType string) bool {
	return strings.Contains(j.request.Header.Get("Content-Type"), contentType)
}

// Languages is specified on the Job interface.
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
	sort.Reverse(languages)
	return languages
}

// createPath creates a path out of the major URL parts.
func (j *job) createPath(domain, resource, resourceID string) string {
	path := j.environment.BasePath() + domain + "/" + resource
	if resourceID != "" {
		path = path + "/" + resourceID
	}
	return path
}

// InternalPath is specified on the Job interface.
func (j *job) InternalPath(domain, resource, resourceID string, query ...KeyValue) string {
	path := j.createPath(domain, resource, resourceID)
	if len(query) > 0 {
		path += "?" + KeyValues(query).String()
	}
	return path
}

// Redirect is specified on the Job interface.
func (j *job) Redirect(domain, resource, resourceID string) {
	path := j.createPath(domain, resource, resourceID)
	http.Redirect(j.responseWriter, j.request, path, http.StatusTemporaryRedirect)
}

// RenderTemplate is specified on the Job interface.
func (j *job) RenderTemplate(templateID string, data interface{}) error {
	return j.environment.Templates().Render(j.responseWriter, templateID, data)
}

// GOB is specified on the Job interface.
func (j *job) GOB() Formatter {
	return &gobFormatter{j}
}

// JSON is specified on the Job interface.
func (j *job) JSON(html bool) Formatter {
	return &jsonFormatter{j, html}
}

// XML is specified on the Job interface.
func (j *job) XML() Formatter {
	return &xmlFormatter{j}
}

// EOF
