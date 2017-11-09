// Tideland GoREST - REST - Tools
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
	"fmt"
	"net/url"
	"strings"
)

//--------------------
// LANGUAGE
//--------------------

// Language is the valued language a request accepts as response.
type Language struct {
	Locale string
	Value  float64
}

// Languages is the ordered set of accepted languages.
type Languages []Language

// Len returns the number of languages to fulfill the sort interface.
func (ls Languages) Len() int {
	return len(ls)
}

// Less returns if the language with the index i has a smaller
// value than the one with index j to fulfill the sort interface.
func (ls Languages) Less(i, j int) bool {
	return ls[i].Value < ls[j].Value
}

// Swap swaps the languages with the indexes i and j.
func (ls Languages) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

//--------------------
// KEY / VALUE
//--------------------

// KeyValue assigns a value to a key.
type KeyValue struct {
	Key   string
	Value interface{}
}

// String prints the encoded form key=value for URLs.
func (kv KeyValue) String() string {
	return fmt.Sprintf("%v=%v", url.QueryEscape(kv.Key), url.QueryEscape(fmt.Sprintf("%v", kv.Value)))
}

// KeyValues is a number of key/value pairs.
type KeyValues []KeyValue

// String prints the encoded form key=value joind by & for URLs.
func (kvs KeyValues) String() string {
	kvss := make([]string, len(kvs))
	for i, kv := range kvs {
		kvss[i] = kv.String()
	}
	return strings.Join(kvss, "&")
}

// EOF
