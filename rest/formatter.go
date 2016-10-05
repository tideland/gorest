// Tideland Go REST Server Library - REST - Formatter
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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tideland/golib/errors"
)

//--------------------
// CONST
//--------------------

const (
	// Standard REST status codes.
	StatusOK           = http.StatusOK
	StatusCreated      = http.StatusCreated
	StatusNoContent    = http.StatusNoContent
	StatusBadRequest   = http.StatusBadRequest
	StatusUnauthorized = http.StatusUnauthorized
	StatusNotFound     = http.StatusNotFound
	StatusConflict     = http.StatusConflict

	// Standard REST content types.
	ContentTypePlain = "text/plain"
	ContentTypeHTML  = "text/html"
	ContentTypeXML   = "application/xml"
	ContentTypeJSON  = "application/json"
	ContentTypeGOB   = "application/vnd.tideland.gob"
)

//--------------------
// ENVELOPE
//--------------------

// Envelope is a helper to give a qualified feedback in RESTful requests.
// It contains wether the request has been successful, in case of an
// error an additional message and the payload.
type Envelope struct {
	Success bool
	Message string
	Payload interface{}
}

//--------------------
// FORMATTER
//--------------------

type Formatter interface {
	// Write encodes the passed data to implementers format and writes
	// it with the passed status code to the response writer.
	Write(status int, data interface{}) error

	// Read checks if the request content type matches the implementers
	// format, reads its body and decodes it to the value pointed to by
	// data.
	Read(data interface{}) error
}

// PositiveFeedback writes a positive feedback envelope to the formatter.
func PositiveFeedback(f Formatter, payload interface{}, msg string, args ...interface{}) error {
	fmsg := fmt.Sprintf(msg, args...)
	return f.Write(StatusOK, &Envelope{true, fmsg, payload})
}

// NegativeFeedback writes a negative feedback envelope to the formatter.
func NegativeFeedback(f Formatter, status int, msg string, args ...interface{}) error {
	fmsg := fmt.Sprintf(msg, args...)
	return f.Write(status, &Envelope{false, fmsg, nil})
}

//--------------------
// GOB FORMATTER
//--------------------

// gobFormatter implements Formatter for the GOB encoding.
type gobFormatter struct {
	job Job
}

// Write is specified on the Formatter interface.
func (gf *gobFormatter) Write(status int, data interface{}) error {
	enc := gob.NewEncoder(gf.job.ResponseWriter())
	gf.job.ResponseWriter().WriteHeader(status)
	gf.job.ResponseWriter().Header().Set("Content-Type", ContentTypeGOB)
	err := enc.Encode(data)
	if err != nil {
		http.Error(gf.job.ResponseWriter(), err.Error(), http.StatusInternalServerError)
	}
	return err
}

// Read is specified on the Formatter interface.
func (gf *gobFormatter) Read(data interface{}) error {
	if !gf.job.HasContentType(ContentTypeGOB) {
		return errors.New(ErrInvalidContentType, errorMessages, ContentTypeGOB)
	}
	dec := gob.NewDecoder(gf.job.Request().Body)
	err := dec.Decode(data)
	gf.job.Request().Body.Close()
	return err
}

//--------------------
// JSON FORMATTER
//--------------------

// jsonFormatter implements Formatter for the JSON encoding. Writing
// also can be done with HTML escaping.
type jsonFormatter struct {
	job  Job
	html bool
}

// Write is specified on the Formatter interface.
func (jf *jsonFormatter) Write(status int, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(jf.job.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return err
	}
	if jf.html {
		var buf bytes.Buffer
		json.HTMLEscape(&buf, body)
		body = buf.Bytes()
	}
	jf.job.ResponseWriter().WriteHeader(status)
	jf.job.ResponseWriter().Header().Set("Content-Type", ContentTypeJSON)
	_, err = jf.job.ResponseWriter().Write(body)
	return err
}

// Read is specified on the Formatter interface.
func (jf *jsonFormatter) Read(data interface{}) error {
	if !jf.job.HasContentType(ContentTypeJSON) {
		return errors.New(ErrInvalidContentType, errorMessages, ContentTypeJSON)
	}
	body, err := ioutil.ReadAll(jf.job.Request().Body)
	jf.job.Request().Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &data)
}

//--------------------
// XML FORMATTER
//--------------------

// xmlFormatter implements Formatter for the XML encoding.
type xmlFormatter struct {
	job Job
}

// Write is specified on the Formatter interface.
func (xf *xmlFormatter) Write(status int, data interface{}) error {
	body, err := xml.Marshal(data)
	if err != nil {
		http.Error(xf.job.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return err
	}
	xf.job.ResponseWriter().WriteHeader(status)
	xf.job.ResponseWriter().Header().Set("Content-Type", ContentTypeXML)
	_, err = xf.job.ResponseWriter().Write(body)
	return err
}

// Read is specified on the Formatter interface.
func (xf *xmlFormatter) Read(data interface{}) error {
	if !xf.job.HasContentType(ContentTypeXML) {
		return errors.New(ErrInvalidContentType, errorMessages, ContentTypeXML)
	}
	body, err := ioutil.ReadAll(xf.job.Request().Body)
	xf.job.Request().Body.Close()
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, &data)
}

// EOF
