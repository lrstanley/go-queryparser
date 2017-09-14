// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package queryparser

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		raw     string
		allowed []string
	}
	tests := []struct {
		name    string
		args    args
		wantQry Query
	}{
		{name: "Empty raw", args: args{raw: "", allowed: []string{}}, wantQry: Query{filters: map[string][]string{}}},
		{name: "some invalid", args: args{raw: "this #$% test  ", allowed: []string{}}, wantQry: Query{
			Orig:    "this #$% test  ",
			Raw:     "this test",
			Pretty:  "this test",
			filters: map[string][]string{},
		}},
		{name: "No filters", args: args{raw: "this is a test", allowed: []string{}}, wantQry: Query{
			Orig:    "this is a test",
			Raw:     "this is a test",
			Pretty:  "this is a test",
			filters: map[string][]string{},
		}},
		{name: "Filters but all allowed", args: args{raw: "this:1 is:foo a test", allowed: []string{}}, wantQry: Query{
			Orig:   "this:1 is:foo a test",
			Raw:    "a test",
			Pretty: "this:1 is:foo a test",
			filters: map[string][]string{
				"this": {"1"}, "is": {"foo"},
			},
		}},
		{name: "Filters but only 'this' allowed", args: args{raw: "this:1 is:foo a test", allowed: []string{"this"}}, wantQry: Query{
			Orig:   "this:1 is:foo a test",
			Raw:    "is:foo a test",
			Pretty: "this:1 is:foo a test",
			filters: map[string][]string{
				"this": {"1"},
			},
		}},
		{name: "Just a filter", args: args{raw: "this:test", allowed: []string{}}, wantQry: Query{
			Orig:   "this:test",
			Raw:    "",
			Pretty: "this:test",
			filters: map[string][]string{
				"this": {"test"},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotQry := New(tt.args.raw, tt.args.allowed); !reflect.DeepEqual(gotQry, tt.wantQry) {
				t.Errorf("New() = %#v, want %#v", gotQry, tt.wantQry)
			}
		})
	}
}
