// Tideland GoREST - Handlers - File Serve
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
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gorest/rest"
)

//--------------------
// FILE SERVE HANDLER
//--------------------

// fileServeHandler implements the file server.
type fileServeHandler struct {
	id  string
	dir string
}

// NewFileServeHandler creates a new handler serving the files names
// by the resource ID part out of the passed directory.
func NewFileServeHandler(id, dir string) rest.ResourceHandler {
	pdir := filepath.FromSlash(dir)
	if !strings.HasSuffix(pdir, string(filepath.Separator)) {
		pdir += string(filepath.Separator)
	}
	return &fileServeHandler{id, pdir}
}

// ID is specified on the ResourceHandler interface.
func (h *fileServeHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *fileServeHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *fileServeHandler) Get(job rest.Job) (bool, error) {
	filename, err := filepath.Abs(filepath.Clean(h.dir + job.ResourceID()))
	if err != nil {
		logger.Errorf("file '%s' does not exist", filename)
		return false, errors.Annotate(err, ErrDownloadingFile, errorMessages, filename)
	}
	logger.Infof("serving file '%s'", filename)
	http.ServeFile(job.ResponseWriter(), job.Request(), filename)
	return true, nil
}

// EOF
