// Tideland Go REST Server Library - REST - Templates
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
	"io/ioutil"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/tideland/golib/errors"
)

//--------------------
// TEMPLATE CACHE ITEM
//--------------------

// templatesItem stores the parsed template and the
// content type.
type templatesItem struct {
	id             string
	timestamp      time.Time
	parsedTemplate *template.Template
	contentType    string
}

// isValid checks if the the entry is younger than the
// passed validity period.
func (ti *templatesItem) isValid(validityPeriod time.Duration) bool {
	return ti.timestamp.Add(validityPeriod).After(time.Now())
}

// render the cached entry.
func (ti *templatesItem) render(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", ti.contentType)
	if err := ti.parsedTemplate.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

//--------------------
// TEMPLATE CACHE
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
}

// templates implements the TemplatesCache interface.
type templates struct {
	mutex sync.RWMutex
	items map[string]*templatesItem
}

// NewTemplates creates a new template cache.
func NewTemplatesCache() TemplatesCache {
	return &templates{
		items: make(map[string]*templatesItem),
	}
}

// Parse is specified on the Templates interface.
func (t *templates) Parse(id, rawTemplate, contentType string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	parsedTemplate, err := template.New(id).Parse(rawTemplate)
	if err != nil {
		return err
	}
	t.items[id] = &templatesItem{id, time.Now(), parsedTemplate, contentType}
	return nil
}

// LoadAndParse is specified on the Templates interface.
func (t *templates) LoadAndParse(id, filename, contentType string) error {
	rawTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return t.Parse(id, string(rawTemplate), contentType)
}

// Render is specified on the Templates interface.
func (t *templates) Render(rw http.ResponseWriter, id string, data interface{}) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	entry, ok := t.items[id]
	if !ok {
		return errors.New(ErrNoCachedTemplate, errorMessages, id)
	}
	return entry.render(rw, data)
}

// EOF
