// Tideland Go REST Server Library - Handlers - File Serve
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tideland/golib/logger"
	
	"github.com/tideland/gorest/rest"
)

//--------------------
// FILE SERVER HANDLER
//--------------------

// FileServeHandler serves files identified by the resource ID part out
// of the configured local directory.
type FileServeHandler struct {
	id  string
	dir string
}

// NewFileServeHandler creates a new handler with a directory.
func NewFileServeHandler(id, dir string) rest.ResourceHandler {
	pdir := filepath.FromSlash(dir)
	if !strings.HasSuffix(pdir, string(filepath.Separator)) {
		pdir += string(filepath.Separator)
	}
	return &FileServeHandler{id, pdir}
}

// ID is specified on the ResourceHandler interface.
func (h *FileServeHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *FileServeHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *FileServeHandler) Get(job rest.Job) (bool, error) {
	filename := h.dir + job.ResourceID()
	logger.Infof("serving file %q", filename)
	http.ServeFile(job.ResponseWriter(), job.Request(), filename)
	return true, nil
}

// EOF
