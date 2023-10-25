// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package queryparser parses a common "q" http GET variable to look for text
// filters, which can be used for advanced searching. For example:
//    Hello World tags:example,world foo:"something quoted" author:lrstanley
//
// go-queryparser also strips out not very safe/useful characters by default.
// See the DefaultCut function for details (and New() for how you can do this
// yourself.)
package queryparser

import (
	"sort"
	"strings"
)

// Query represents filtered input.
type Query struct {
	// Raw is the raw (trailing) text of items that weren't filters.
	Raw     string
	Filters map[string][]string
}

func (q *Query) Add(key, val string) {
	val = stripDuplicateWS(val)
	var vals []string

	if strings.HasPrefix(val, `"`) {
		vals = []string{strings.Trim(val, `"`)}
	} else if strings.HasPrefix(val, `'`) {
		vals = []string{strings.Trim(val, `'`)}
	} else {
		vals = strings.FieldsFunc(val, func(r rune) bool {
			return r == ','
		})
	}

	key = strings.ToLower(key)
	if _, ok := q.Filters[key]; !ok {
		q.Filters[key] = make([]string, 0, 1)
	}

	q.Filters[key] = append(q.Filters[key], vals...)
}

// Has returns true if there is a filter matching the given name.
func (q *Query) Has(key string) (exists bool) {
	_, exists = q.Filters[strings.ToLower(key)]
	return exists
}

// Get returns the results of the filter if it exists, and if it successfully
// found a result.
func (q *Query) Get(key string) (results []string, ok bool) {
	results, ok = q.Filters[strings.ToLower(key)]
	return results, ok
}

// GetOne returns the last known result for the filter, if it exists. Useful
// if you only want a user to define a filter once. The resulting string
// is empty if no filter of that key was found.
func (q *Query) GetOne(key string) string {
	results, ok := q.Get(key)

	if !ok {
		return ""
	}

	out := results[len(results)-1]

	return out
}

// String will return a string representation of the query, however all filters
// will be sorted alphabetically, and the raw text will be placed at the end
// of the string.
func (q *Query) String() (out string) {
	skeys := make([]string, 0, len(q.Filters))
	for key := range q.Filters {
		skeys = append(skeys, key)
	}
	sort.Strings(skeys)

	for _, key := range skeys {
		vals := q.Filters[key]
		sort.Strings(vals)

		for _, val := range vals {
			out += key + `:"` + val + `" `
		}
	}

	if q.Raw != "" {
		out += q.Raw
	}

	return strings.TrimSpace(out)
}
