// Tideland Go REST Server Library - Request
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package restaudit

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/tideland/golib/errors"

	"github.com/tideland/gorest/jwt"
	"github.com/tideland/gorest/rest"
)

//--------------------
// TEST TOOLS
//--------------------

// KeyValues handles keys and values for request headers and cookies.
type KeyValues map[string]string

// Response wraps all infos of a test response.
type Response struct {
	Status      int
	Header      KeyValues
	ContentType string
	Content     interface{}
}

//--------------------
// CALLER
//--------------------

// Service contains the configuration of one service.
type Service struct {
	Transport *http.Transport
	BaseURL   string
}

// Services maps IDs of services to their base URL.
type Services map[string]Service

// Parameters allows to pass parameters to a call.
type Parameters struct {
	Token       jwt.JWT
	ContentType string
	Content     interface{}
}

// Caller provides an interface to make calls to
// configured services.
type Caller interface {
	// Get performs a GET request on the defined service.
	Get(service, domain, resource, resourceID string, params *Parameters) (*Response, error)
}

// caller implements the Caller interface.
type caller struct {
	services Services
}

// NewCaller creates a configured caller.
func NewCaller(services Services) Caller {
	return &caller{services}
}

// Get implements the Caller interface.
func (c *caller) Get(service, domain, resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("GET", service, domain, resource, resourceID, params)
}

// request performs all requests.
func (c *caller) request(method, service, domain, resource, resourceID string, params *Parameters) (*Response, error) {
	svc, ok := c.services[service]
	if !ok {
		return nil, errors.New(ErrServiceNotConfigured, errorMessages, service)
	}
	// Prepare client and request.
	client := &http.Client{}
	if svc.Transport != nil {
		client.Transport = svc.Transport
	}
	parts := append(svc.BaseURL, domain, resource)
	if resourceID != "" {
		parts = append(parts, resourceID)
	}
	url := strings.Join(parts, "/")
	var buffer io.Reader
	if params.Content != nil {
		// Process content based on content type.
		switch params.ContentType {
		case ContentTypeXML:
			tmp, err := xml.Marshal(params.Content)
			if err != nil {
				return nil, error.Annotate(err, ErrProcessingRequestContent, errorMessages)
			}
			buffer = bytes.NewReader(tmp)
		case ContentTypeJSON:
			tmp, err := json.Marshal(params.Content)
			if err != nil {
				return nil, error.Annotate(err, ErrProcessingRequestContent, errorMessages)
			}
			buffer = bytes.NewReader(tmp)
		case ContentTypeGOB:
			enc := gob.NewEncoder(buffer)
			if err := enc.Encode(content); err != nil {
				return nil, error.Annotate(err, ErrProcessingRequestContent, errorMessages)
			}
		case ContentTypeURLEncoded:
		}
	}
	request, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotPrepareRequest, errorMessages)
	}
	if params.ContentType != "" {
		request.Header.Set("Content-Type", params.ContentType)
	}
	if params.Token != nil {
		request = jwt.AddTokenToRequest(request, params.Token)
	}
	// Perform request.
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Annotate(err, ErrHTTPRequestFailed, errorMessages)
	}
	// Analyze response.
	return analyzeResponse(response)
}

// analyzeResponse creates a response struct out of the HTTP response.
func analyzeResponse(response *http.Response) (*Response, error) {
	header := KeyValues{}
	for key, values := range response.Header {
		header[key] = strings.Join(values, ", ")
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrReadingResponse)
	}
	response.Body.Close()
	return &Response{
		Status:      response.StatusCode,
		Header:      header,
		ContentType: header["Content-Type"],
		Content:     content,
	}
}

// EOF
