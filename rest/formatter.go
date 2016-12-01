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
	"net/url"
	"time"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/stringex"
)

//--------------------
// CONST
//--------------------

const (
	// Standard REST status codes.
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusNoContent           = http.StatusNoContent
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusNotFound            = http.StatusNotFound
	StatusConflict            = http.StatusConflict
	StatusInternalServerError = http.StatusInternalServerError

	// Standard REST content types.
	ContentTypePlain      = "text/plain"
	ContentTypeHTML       = "text/html"
	ContentTypeXML        = "application/xml"
	ContentTypeJSON       = "application/json"
	ContentTypeGOB        = "application/vnd.tideland.gob"
	ContentTypeURLEncoded = "application/x-www-form-urlencoded"
)

//--------------------
// GLOBAL
//--------------------

var (
	defaulter = stringex.NewDefaulter("job", false)
)

//--------------------
// ENVELOPE
//--------------------

// envelope is a helper to give a qualified feedback in RESTful requests.
// It contains wether the request has been successful, in case of an
// error an additional message and the payload.
type envelope struct {
	Status  string      `json:"status" xml:"status"`
	Message string      `json:"message,omitempty" xml:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty" xml:"payload,omitempty"`
}

//--------------------
// FORMATTER
//--------------------

type Formatter interface {
	// Write encodes the passed data to implementers format and writes
	// it with the passed status code and possible header values to the
	// response writer.
	Write(status int, data interface{}, headers ...KeyValue) error

	// Read checks if the request content type matches the implementers
	// format, reads its body and decodes it to the value pointed to by
	// data.
	Read(data interface{}) error
}

// PositiveFeedback writes a positive feedback envelope to the formatter.
func PositiveFeedback(f Formatter, payload interface{}, msg string, args ...interface{}) error {
	fmsg := fmt.Sprintf(msg, args...)
	return f.Write(StatusOK, envelope{"success", fmsg, payload})
}

// NegativeFeedback writes a negative feedback envelope to the formatter.
func NegativeFeedback(f Formatter, status int, msg string, args ...interface{}) error {
	fmsg := fmt.Sprintf(msg, args...)
	return f.Write(status, envelope{"fail", fmsg, nil})
}

//--------------------
// GOB FORMATTER
//--------------------

// gobFormatter implements Formatter for the GOB encoding.
type gobFormatter struct {
	job Job
}

// Write is specified on the Formatter interface.
func (gf *gobFormatter) Write(status int, data interface{}, headers ...KeyValue) error {
	enc := gob.NewEncoder(gf.job.ResponseWriter())
	for _, header := range headers {
		gf.job.ResponseWriter().Header().Add(header.Key, fmt.Sprintf("%v", header.Value))
	}
	gf.job.ResponseWriter().Header().Set("Content-Type", ContentTypeGOB)
	gf.job.ResponseWriter().Header().Set("Version", gf.job.Version().String())
	gf.job.ResponseWriter().WriteHeader(status)
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
func (jf *jsonFormatter) Write(status int, data interface{}, headers ...KeyValue) error {
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
	for _, header := range headers {
		jf.job.ResponseWriter().Header().Add(header.Key, fmt.Sprintf("%v", header.Value))
	}
	jf.job.ResponseWriter().Header().Set("Content-Type", ContentTypeJSON)
	jf.job.ResponseWriter().Header().Set("Version", jf.job.Version().String())
	jf.job.ResponseWriter().WriteHeader(status)
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
func (xf *xmlFormatter) Write(status int, data interface{}, headers ...KeyValue) error {
	body, err := xml.Marshal(data)
	if err != nil {
		http.Error(xf.job.ResponseWriter(), err.Error(), http.StatusInternalServerError)
		return err
	}
	for _, header := range headers {
		xf.job.ResponseWriter().Header().Add(header.Key, fmt.Sprintf("%v", header.Value))
	}
	xf.job.ResponseWriter().Header().Set("Content-Type", ContentTypeXML)
	xf.job.ResponseWriter().Header().Set("Version", xf.job.Version().String())
	xf.job.ResponseWriter().WriteHeader(status)
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

//--------------------
// QUERY
//--------------------

// Query allows typed access with default values to a jobs
// request values passed as query.
type Query interface {
	// ValueAsString retrieves the string value of a given key. If it
	// doesn't exist the default value dv is returned.
	ValueAsString(key, dv string) string

	// ValueAsBool retrieves the bool value of a given key. If it
	// doesn't exist the default value dv is returned.
	ValueAsBool(key string, dv bool) bool

	// ValueAsInt retrieves the int value of a given key. If it
	// doesn't exist the default value dv is returned.
	ValueAsInt(key string, dv int) int

	// ValueAsFloat64 retrieves the float64 value of a given key. If it
	// doesn't exist the default value dv is returned.
	ValueAsFloat64(key string, dv float64) float64

	// ValueAsTime retrieves the string value of a given key and
	// interprets it as time with the passed format. If it
	// doesn't exist the default value dv is returned.
	ValueAsTime(key, layout string, dv time.Time) time.Time

	// ValueAsDuration retrieves the duration value of a given key.
	// If it doesn't exist the default value dv is returned.
	ValueAsDuration(key string, dv time.Duration) time.Duration
}

// query implements Query.
type query struct {
	values url.Values
}

// ValueAsString implements the Query interface.
func (q *query) ValueAsString(key, dv string) string {
	value := queryValuer(q.values.Get(key))
	return defaulter.AsString(value, dv)
}

// ValueAsBool implements the Query interface.
func (q *query) ValueAsBool(key string, dv bool) bool {
	value := queryValuer(q.values.Get(key))
	return defaulter.AsBool(value, dv)
}

// ValueAsInt implements the Query interface.
func (q *query) ValueAsInt(key string, dv int) int {
	value := queryValuer(q.values.Get(key))
	return defaulter.AsInt(value, dv)
}

// ValueAsFloat64 implements the Query interface.
func (q *query) ValueAsFloat64(key string, dv float64) float64 {
	value := queryValuer(q.values.Get(key))
	return defaulter.AsFloat64(value, dv)
}

// ValueAsTime implements the Query interface.
func (q *query) ValueAsTime(key, format string, dv time.Time) time.Time {
	value := queryValuer(q.values.Get(key))
	return defaulter.AsTime(value, format, dv)
}

// ValueAsDuration implements the Query interface.
func (q *query) ValueAsDuration(key string, dv time.Duration) time.Duration {
	value := queryValuer(q.values.Get(key))
	return defaulter.AsDuration(value, dv)
}

// queryValues implements the stringex.Valuer interface for
// the usage inside of query.
type queryValuer string

// Value implements the Valuer interface.
func (qv queryValuer) Value() (string, error) {
	v := string(qv)
	if len(v) == 0 {
		return "", errors.New(ErrQueryValueNotFound, errorMessages)
	}
	return v, nil
}

// EOF
