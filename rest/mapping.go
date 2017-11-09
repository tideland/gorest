// Tideland GoREST - REST - Mapping
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
	"strings"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/logger"
)

//--------------------
// HANDLER LIST
//--------------------

// handlerListEntry is one entry in a list of resource handlers.
type handlerListEntry struct {
	handler ResourceHandler
	next    *handlerListEntry
}

// handlerList maintains a list of handlers responsible
// for one domain and resource.
type handlerList struct {
	head *handlerListEntry
}

// register adds a new resource handler.
func (hl *handlerList) register(handler ResourceHandler) error {
	if hl.head == nil {
		hl.head = &handlerListEntry{handler, nil}
		return nil
	}
	current := hl.head
	for {
		if current.handler == handler {
			return errors.New(ErrDuplicateHandler, errorMessages, handler.ID())
		}
		if current.next == nil {
			break
		}
		current = current.next
	}
	current.next = &handlerListEntry{handler, nil}
	return nil
}

// deregister removes a resource handler.
func (hl *handlerList) deregister(ids ...string) {
	// Check if all shall be deregistered.
	if len(ids) == 0 {
		hl.head = nil
		return
	}
	// No, so iterate over ids.
	for _, id := range ids {
		var head, tail *handlerListEntry
		current := hl.head
		for current != nil {
			if current.handler.ID() != id {
				if head == nil {
					head = current
					tail = current
				} else {
					tail.next = current
					tail = tail.next
				}
			}
			current = current.next
		}
		hl.head = head
	}
}

// ids returns the handler ids of this handler list.
func (hl *handlerList) ids() []string {
	ids := []string{}
	current := hl.head
	for current != nil {
		ids = append(ids, current.handler.ID())
		current = current.next
	}
	return ids
}

// handle lets all resource handlers process the request.
func (hl *handlerList) handle(job Job) error {
	current := hl.head
	for current != nil {
		goOn, err := handleJob(current.handler, job)
		if err != nil {
			return err
		}
		if !goOn {
			return nil
		}
		current = current.next
	}
	return nil
}

//--------------------
// MAPPING
//--------------------

// mapping maps domains and resources to lists of
// resource handlers.
type mapping struct {
	ignoreFavicon bool
	handlers      map[string]*handlerList
}

// newMapping returns a new handler mapping.
func newMapping(cfg etc.Etc) *mapping {
	return &mapping{
		ignoreFavicon: cfg.ValueAsBool("ignore-favicon", true),
		handlers:      make(map[string]*handlerList),
	}
}

// register adds a resource handler.
func (m *mapping) register(domain, resource string, handler ResourceHandler) error {
	location := m.location(domain, resource)
	hl, ok := m.handlers[location]
	if !ok {
		hl = &handlerList{}
		m.handlers[location] = hl
	}
	return hl.register(handler)
}

// registeredHandlers returns the IDs of the registered resource handlers.
func (m *mapping) registeredHandlers(domain, resource string) []string {
	location := m.location(domain, resource)
	hl, ok := m.handlers[location]
	if !ok {
		return nil
	}
	return hl.ids()
}

// deregister removes a resource handler.
func (m *mapping) deregister(domain, resource string, ids ...string) {
	location := m.location(domain, resource)
	hl, ok := m.handlers[location]
	if !ok {
		return
	}
	hl.deregister(ids...)
	if hl.head == nil {
		delete(m.handlers, location)
	}
}

// handle handles a request.
func (m *mapping) handle(job Job) error {
	// Check for favicon.ico.
	if m.ignoreFavicon {
		if job.Domain() == "favicon.ico" {
			job.ResponseWriter().WriteHeader(StatusNoContent)
			return nil
		}
	}
	// Find handler list.
	hl, err := m.handlerList(job)
	if err != nil {
		return err
	}
	// Let the handler list handle the job.
	logger.Infof("handling %s", job)
	return hl.handle(job)
}

// handlerList retrieves the handler list for the job.
func (m *mapping) handlerList(job Job) (*handlerList, error) {
	location := m.location(job.Domain(), job.Resource())
	hl, ok := m.handlers[location]
	if ok {
		return hl, nil
	}
	location = m.location(job.Domain(), job.Environment().DefaultResource())
	hl, ok = m.handlers[location]
	if ok {
		return hl, nil
	}
	location = m.location(job.Environment().DefaultDomain(), job.Environment().DefaultResource())
	hl, ok = m.handlers[location]
	if ok {
		return hl, nil
	}
	return nil, errors.New(ErrNoHandler, errorMessages, location)
}

// location builds the map key for domain and resource.
func (m *mapping) location(domain, resource string) string {
	return strings.ToLower(domain + "/" + resource)
}

// EOF
