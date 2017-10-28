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
		name string
		args args
		want Query
	}{
		{name: "empty raw", args: args{raw: "", allowed: []string{}}, want: Query{filters: map[string][]string{}}},
		{name: "some invalid", args: args{raw: "this #$% test  ", allowed: []string{}}, want: Query{
			Orig:    "this #$% test  ",
			Raw:     "this test",
			Pretty:  "this test",
			filters: map[string][]string{},
		}},
		{name: "no filters", args: args{raw: "this is a test", allowed: []string{}}, want: Query{
			Orig:    "this is a test",
			Raw:     "this is a test",
			Pretty:  "this is a test",
			filters: map[string][]string{},
		}},
		{name: "filters but all allowed", args: args{raw: "foo:1 bar:foo a test", allowed: []string{}}, want: Query{
			Orig:   "foo:1 bar:foo a test",
			Raw:    "a test",
			Pretty: "foo:1 bar:foo a test",
			filters: map[string][]string{
				"foo": {"1"}, "bar": {"foo"},
			},
		}},
		{name: "filters but only 'foo' allowed", args: args{raw: "foo:1 bar:foo a :test", allowed: []string{"foo"}}, want: Query{
			Orig:   "foo:1 bar:foo a :test",
			Raw:    "bar:foo a :test",
			Pretty: "foo:1 bar:foo a :test",
			filters: map[string][]string{
				"foo": {"1"},
			},
		}},
		{name: "1 filter", args: args{raw: "foo:test", allowed: []string{}}, want: Query{
			Orig:   "foo:test",
			Raw:    "",
			Pretty: "foo:test",
			filters: map[string][]string{
				"foo": {"test"},
			},
		}},
		{name: "2 filters", args: args{raw: "foo:test bar:test1", allowed: []string{"foo"}}, want: Query{
			Orig:   "foo:test bar:test1",
			Raw:    "bar:test1",
			Pretty: "foo:test bar:test1",
			filters: map[string][]string{
				"foo": {"test"},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.raw, tt.args.allowed)

			if tt.args.raw == "" && !got.IsZero() {
				t.Error("Query.IsZero failed with empty input")
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %#v, want %#v", got, tt.want)
			}

			if len(got.filters) > 0 && len(tt.args.allowed) > 0 {
				for key := range got.filters {
					var isin bool
					for i := 0; i < len(tt.args.allowed); i++ {
						if tt.args.allowed[i] == key {
							isin = true
							break
						}
					}

					if !isin {
						t.Errorf("New() = %#v, returned filter that wasn't allowed", got)
					} else {
						if !got.Has(key) {
							t.Errorf("New() = %#v, Query.Has(%q) = false, wanted true", got, key)
						}
					}
				}
			}
		})
	}
}
