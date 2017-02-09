// Tideland Go REST Server Library - Request
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/version"

	"github.com/tideland/gorest/jwt"
	"github.com/tideland/gorest/rest"
)

//--------------------
// SERVERS
//--------------------

// key is to address the servers inside a context.
type key int

var serversKey key = 0

// server contains the configuration of one server.
type server struct {
	URL       string
	Transport *http.Transport
}

// Servers maps IDs of domains to their server configurations.
// Multiple ones can be added per domain for spreading the
// load or provide higher availability.
type Servers interface {
	// Add adds a domain server configuration.
	Add(domain string, url string, transport *http.Transport)

	// Caller retrieves a caller for a domain.
	Caller(domain string) (Caller, error)
}

// servers implements servers.
type servers struct {
	mutex   sync.RWMutex
	servers map[string][]*server
}

// NewServers creates a new servers manager.
func NewServers() Servers {
	rand.Seed(time.Now().Unix())
	return &servers{
		servers: make(map[string][]*server),
	}
}

// Add implements the Servers interface.
func (s *servers) Add(domain, url string, transport *http.Transport) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	srvs, ok := s.servers[domain]
	if ok {
		s.servers[domain] = append(srvs, &server{url, transport})
		return
	}
	s.servers[domain] = []*server{&server{url, transport}}
}

// Caller implements the Servers interface.
func (s *servers) Caller(domain string) (Caller, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	srvs, ok := s.servers[domain]
	if !ok {
		return nil, errors.New(ErrNoServerDefined, errorMessages, domain)
	}
	return newCaller(domain, srvs), nil
}

// NewContext returns a new context that carries configured servers.
func NewContext(ctx context.Context, servers Servers) context.Context {
	return context.WithValue(ctx, serversKey, servers)
}

// FromContext returns the servers configuration stored in ctx, if any.
func FromContext(ctx context.Context) (Servers, bool) {
	servers, ok := ctx.Value(serversKey).(Servers)
	return servers, ok
}

//--------------------
// RESPONSE
//--------------------

// KeyValues handles keys and values for request headers and cookies.
type KeyValues map[string]string

// Response wraps all infos of a test response.
type Response interface {
	// StatusCode returns the HTTP status code of the response.
	StatusCode() int

	// Header returns the HTTP header of the response.
	Header() http.Header

	// HasContentType checks the content type regardless of charsets.
	HasContentType(contentType string) bool

	// Read decodes the content into the passed data depending
	// on the content type.
	Read(data interface{}) error
}

// response implements Response.
type response struct {
	httpResp    *http.Response
	contentType string
	content     []byte
}

// StatusCode implements the Response interface.
func (r *response) StatusCode() int {
	return r.httpResp.StatusCode
}

// Header implements the Response interface.
func (r *response) Header() http.Header {
	return r.httpResp.Header
}

// HasContentType implements the Response interface.
func (r *response) HasContentType(contentType string) bool {
	return strings.Contains(r.contentType, contentType)
}

// Read implements the Response interface.
func (r *response) Read(data interface{}) error {
	switch {
	case r.HasContentType(rest.ContentTypeGOB):
		dec := gob.NewDecoder(bytes.NewBuffer(r.content))
		if err := dec.Decode(data); err != nil {
			return errors.Annotate(err, ErrDecodingResponse, errorMessages)
		}
		return nil
	case r.HasContentType(rest.ContentTypeJSON):
		if err := json.Unmarshal(r.content, &data); err != nil {
			return errors.Annotate(err, ErrDecodingResponse, errorMessages)
		}
		return nil
	case r.HasContentType(rest.ContentTypeXML):
		if err := xml.Unmarshal(r.content, &data); err != nil {
			return errors.Annotate(err, ErrDecodingResponse, errorMessages)
		}
		return nil
	case r.HasContentType(rest.ContentTypeURLEncoded):
		values, err := url.ParseQuery(string(r.content))
		if err != nil {
			return errors.Annotate(err, ErrDecodingResponse, errorMessages)
		}
		// Check for data type url.Values.
		duv, ok := data.(url.Values)
		if ok {
			for key, value := range values {
				duv[key] = value
			}
			return nil
		}
		// Check for data type KeyValues.
		kvv, ok := data.(KeyValues)
		if !ok {
			return errors.New(ErrDecodingResponse, errorMessages)
		}
		for key, value := range values {
			kvv[key] = strings.Join(value, " / ")
		}
		return nil
	}
	return errors.New(ErrInvalidContentType, errorMessages, r.contentType)
}

//--------------------
// CALL PARAMETERS
//--------------------

// Parameters allows to pass parameters to a call.
type Parameters struct {
	Version     version.Version
	Token       jwt.JWT
	ContentType string
	Content     interface{}
	Accept      string
}

// body returns the content as body data depending on
// the content type.
func (p *Parameters) body() (io.Reader, error) {
	buffer := bytes.NewBuffer(nil)
	if p.Content == nil {
		return buffer, nil
	}
	// Process content based on content type.
	switch p.ContentType {
	case rest.ContentTypeXML:
		tmp, err := xml.Marshal(p.Content)
		if err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, errorMessages)
		}
		buffer.Write(tmp)
	case rest.ContentTypeJSON:
		tmp, err := json.Marshal(p.Content)
		if err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, errorMessages)
		}
		buffer.Write(tmp)
	case rest.ContentTypeGOB:
		enc := gob.NewEncoder(buffer)
		if err := enc.Encode(p.Content); err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, errorMessages)
		}
	case rest.ContentTypeURLEncoded:
		values, err := p.values()
		if err != nil {
			return nil, err
		}
		_, err = buffer.WriteString(values.Encode())
		if err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, errorMessages)
		}
	}
	return buffer, nil
}

// values returns the content as URL encoded values.
func (p *Parameters) values() (url.Values, error) {
	if p.Content == nil {
		return url.Values{}, nil
	}
	// Check if type is already ok.
	urlvs, ok := p.Content.(url.Values)
	if ok {
		return urlvs, nil
	}
	// Check for simple key/values.
	kvs, ok := p.Content.(KeyValues)
	if !ok {
		return nil, errors.New(ErrInvalidContent, errorMessages)
	}
	values := url.Values{}
	for key, value := range kvs {
		values.Set(key, value)
	}
	return values, nil
}

//--------------------
// CALLER
//--------------------

// Caller provides an interface to make calls to
// configured services.
type Caller interface {
	// Get performs a GET request on the defined resource.
	Get(resource, resourceID string, params *Parameters) (Response, error)

	// Head performs a HEAD request on the defined resource.
	Head(resource, resourceID string, params *Parameters) (Response, error)

	// Put performs a PUT request on the defined resource.
	Put(resource, resourceID string, params *Parameters) (Response, error)

	// Post performs a POST request on the defined resource.
	Post(resource, resourceID string, params *Parameters) (Response, error)

	// Patch performs a PATCH request on the defined resource.
	Patch(resource, resourceID string, params *Parameters) (Response, error)

	// Delete performs a DELETE request on the defined resource.
	Delete(resource, resourceID string, params *Parameters) (Response, error)

	// Options performs a OPTIONS request on the defined resource.
	Options(resource, resourceID string, params *Parameters) (Response, error)
}

// caller implements the Caller interface.
type caller struct {
	domain string
	srvs   []*server
}

// newCaller creates a configured caller.
func newCaller(domain string, srvs []*server) Caller {
	return &caller{domain, srvs}
}

// Get implements the Caller interface.
func (c *caller) Get(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("GET", resource, resourceID, params)
}

// Head implements the Caller interface.
func (c *caller) Head(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("HEAD", resource, resourceID, params)
}

// Put implements the Caller interface.
func (c *caller) Put(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("PUT", resource, resourceID, params)
}

// Post implements the Caller interface.
func (c *caller) Post(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("POST", resource, resourceID, params)
}

// Patch implements the Caller interface.
func (c *caller) Patch(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("PATCH", resource, resourceID, params)
}

// Delete implements the Caller interface.
func (c *caller) Delete(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("DELETE", resource, resourceID, params)
}

// Options implements the Caller interface.
func (c *caller) Options(resource, resourceID string, params *Parameters) (Response, error) {
	return c.request("OPTIONS", resource, resourceID, params)
}

// request performs all requests.
func (c *caller) request(method, resource, resourceID string, params *Parameters) (Response, error) {
	if params == nil {
		params = &Parameters{}
	}
	// Prepare client.
	// TODO Mue 2016-10-28 Add more algorithms than just random selection.
	srv := c.srvs[rand.Intn(len(c.srvs))]
	client := &http.Client{}
	if srv.Transport != nil {
		client.Transport = srv.Transport
	}
	u, err := url.Parse(srv.URL)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotPrepareRequest, errorMessages)
	}
	upath := strings.Trim(u.Path, "/")
	path := []string{upath, c.domain, resource}
	if resourceID != "" {
		path = append(path, resourceID)
	}
	u.Path = strings.Join(path, "/")
	// Prepare request, check the parameters first.
	var request *http.Request
	if method == "GET" || method == "HEAD" {
		// These allow only URL encoded.
		request, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotPrepareRequest, errorMessages)
		}
		values, err := params.values()
		if err != nil {
			return nil, err
		}
		request.URL.RawQuery = values.Encode()
		request.Header.Set("Content-Type", rest.ContentTypeURLEncoded)
	} else {
		// Here use the body for content.
		body, err := params.body()
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(method, u.String(), body)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotPrepareRequest, errorMessages)
		}
		request.Header.Set("Content-Type", params.ContentType)
	}
	if params.Version != nil {
		request.Header.Set("Version", params.Version.String())
	}
	if params.Token != nil {
		request = jwt.AddTokenToRequest(request, params.Token)
	}
	if params.Accept == "" {
		params.Accept = params.ContentType
	}
	if params.Accept != "" {
		request.Header.Set("Accept", params.Accept)
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
func analyzeResponse(resp *http.Response) (Response, error) {
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrAnalyzingResponse, errorMessages)
	}
	resp.Body.Close()
	return &response{
		httpResp:    resp,
		contentType: resp.Header.Get("Content-Type"),
		content:     content,
	}, nil
}

// EOF
