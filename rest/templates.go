// Tideland GoREST - REST - Templates
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
	"io/ioutil"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// templatesCache CACHE ITEM
//--------------------

// templatesCacheItem stores the parsed template and the
// content type.
type templatesCacheItem struct {
	id             string
	timestamp      time.Time
	parsedTemplate *template.Template
	contentType    string
}

// isValid checks if the the entry is younger than the
// passed validity period.
func (ti *templatesCacheItem) isValid(validityPeriod time.Duration) bool {
	return ti.timestamp.Add(validityPeriod).After(time.Now())
}

// render the cached entry.
func (ti *templatesCacheItem) render(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", ti.contentType)
	if err := ti.parsedTemplate.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

//--------------------
// TEMPLATES CACHE
//--------------------

// TemplatesCache caches and renders templates.
type TemplatesCache interface {
	// Parse parses a raw template an stores it.
	Parse(id, rawTmpl, contentType string) error

	// LoadAndParse loads a template out of the filesystem,
	// parses and stores it.
	LoadAndParse(id, filename, contentType string) error

	// Render executes the pre-parsed template with the data.
	// It also sets the content type header.
	Render(rw http.ResponseWriter, id string, data interface{}) error

	// LoadAndRender checks if the template with the given id
	// has already been parsed. In this case it will use it,
	// otherwise the template will be loaded, parsed, added
	// to the cache, and used then.
	LoadAndRender(rw http.ResponseWriter, id, filename, contentType string, data interface{}) error
}

// templatesCache implements the TemplatesCache interface.
type templatesCache struct {
	mutex sync.RWMutex
	items map[string]*templatesCacheItem
}

// newTemplatesCache creates a new template cache.
func newTemplatesCache() *templatesCache {
	return &templatesCache{
		items: make(map[string]*templatesCacheItem),
	}
}

// Parse impements the TemplatesCache interface.
func (t *templatesCache) Parse(id, rawTemplate, contentType string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	parsedTemplate, err := template.New(id).Parse(rawTemplate)
	if err != nil {
		return err
	}
	t.items[id] = &templatesCacheItem{id, time.Now(), parsedTemplate, contentType}
	return nil
}

// LoadAndParse implements the TemplatesCache interface.
func (t *templatesCache) LoadAndParse(id, filename, contentType string) error {
	rawTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return t.Parse(id, string(rawTemplate), contentType)
}

// Render implements the TemplatesCache interface.
func (t *templatesCache) Render(rw http.ResponseWriter, id string, data interface{}) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	entry, ok := t.items[id]
	if !ok {
		return errors.New(ErrNoCachedTemplate, errorMessages, id)
	}
	return entry.render(rw, data)
}

// LoadAndRender implements the TemplatesCache interface.
func (t *templatesCache) LoadAndRender(rw http.ResponseWriter, id, filename, contentType string, data interface{}) error {
	t.mutex.RLock()
	_, ok := t.items[id]
	t.mutex.RUnlock()
	if !ok {
		if err := t.LoadAndParse(id, filename, contentType); err != nil {
			return err
		}
	}
	return t.Render(rw, id, data)
}

//--------------------
// RENDERER
//--------------------

// Renderer renders templates. It is returned by a Job and knows
// where to render it.
type Renderer interface {
	// Render executes the pre-parsed template with the data.
	// It also sets the content type header.
	Render(id string, data interface{}) error

	// LoadAndRender checks if the template with the given id
	// has already been parsed. In this case it will use it,
	// otherwise the template will be loaded, parsed, added
	// to the cache, and used then.
	LoadAndRender(id, filename, contentType string, data interface{}) error
}

// renderer implements the Renderer interface.
type renderer struct {
	rw http.ResponseWriter
	tc TemplatesCache
}

// Render implements the Renderer interface.
func (r *renderer) Render(id string, data interface{}) error {
	return r.tc.Render(r.rw, id, data)
}

// LoadAndRender implements the Renderer interface.
func (r *renderer) LoadAndRender(id, filename, contentType string, data interface{}) error {
	return r.tc.LoadAndRender(r.rw, id, filename, contentType, data)
}

// EOF
