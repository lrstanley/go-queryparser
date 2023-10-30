// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package queryparser

import (
	"reflect"
	"testing"
)

type caseArgs struct {
	name    string
	input   string
	tokens  []tokenRef
	allowed []string
	query   *Query
}

var cases = []caseArgs{
	{
		name:  "quoted fields",
		input: `foo:"bar" bar:"baz baz1"`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `"bar"`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "bar"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `"baz baz1"`},
		},
		allowed: []string{"foo", "bar"},
		query: &Query{
			Raw: "",
			Filters: map[string][]string{
				"foo": {"bar"},
				"bar": {"baz baz1"},
			},
		},
	},
	{
		name:  "unquoted fields",
		input: `foo:1,2 bar:3`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenIDENT, lit: "1,2"},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "bar"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenIDENT, lit: "3"},
		},
		query: &Query{
			Raw: "",
			Filters: map[string][]string{
				"foo": {"1", "2"},
				"bar": {"3"},
			},
		},
	},
	{
		name:  "trailing",
		input: `foo:"bar" test`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `"bar"`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "test"},
		},
		query: &Query{
			Raw: "test",
			Filters: map[string][]string{
				"foo": {"bar"},
			},
		},
	},
	{
		name:  "trailing with single quote",
		input: `foo:'bar' test`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `'bar'`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "test"},
		},
		query: &Query{
			Raw: "test",
			Filters: map[string][]string{
				"foo": {"bar"},
			},
		},
	},
	{
		name:  "trailing with single and inner double quote",
		input: `foo:'bar' test`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `'ba"r'`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "test"},
		},
		query: &Query{
			Raw: "test",
			Filters: map[string][]string{
				"foo": {"bar"},
			},
		},
	},
	{
		name:  "trailing with random double quotes",
		input: `foo:"bar" test " :" a:"`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `"bar"`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "test"},
			{tok: tokenWS, lit: " "},
			{tok: tokenFIELD, lit: `" :"`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "a"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `"`},
		},
		query: &Query{
			Raw: "test : ",
			Filters: map[string][]string{
				"foo": {"bar"},
				"a":   {},
			},
		},
	},
	{
		name:  "strip DefaultCut",
		input: `foo:"bar" te$st#!`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "foo"},
			{tok: tokenDELIM, lit: ":"},
			{tok: tokenFIELD, lit: `"bar"`},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "test"},
		},
		query: &Query{
			Raw: "test",
			Filters: map[string][]string{
				"foo": {"bar"},
			},
		},
	},
	{
		name:  "text only",
		input: `test test1`,
		tokens: []tokenRef{
			{tok: tokenIDENT, lit: "test"},
			{tok: tokenWS, lit: " "},
			{tok: tokenIDENT, lit: "test1"},
		},
		query: &Query{
			Raw:     "test test1",
			Filters: map[string][]string{},
		},
	},
	{
		name:   "empty",
		input:  ``,
		tokens: []tokenRef{},
		query:  &Query{Filters: map[string][]string{}},
	},
}

func TestScanner(t *testing.T) {
	for _, tt := range cases {
		t.Run("scanner_"+tt.name, func(t *testing.T) {
			s := newScanner(tt.input)

			for _, valid := range tt.tokens {
				tr := s.nextToken()

				// Make sure both have had their fair share of cutsets.
				tr.lit = cutsetFunc(tr.lit, DefaultCut)
				valid.lit = cutsetFunc(valid.lit, DefaultCut)

				if tr.lit != valid.lit || tr.tok != valid.tok {
					t.Fatalf("expected %#v but got %#v", valid, tr)
				}
			}

			if tr := s.nextToken(); tr.tok != tokenEOF {
				t.Fatalf("expected EOF, got %#v", tr)
			}

			s.drain()
		})
	}
}

func TestParser(t *testing.T) {
	for _, tt := range cases {
		t.Run("parser_"+tt.name, func(t *testing.T) {
			p := New(tt.input, Options{Allowed: tt.allowed, CutFn: DefaultCut})
			qp := p.Parse()

			if !reflect.DeepEqual(tt.query, qp) {
				t.Fatalf("expected query %#v, but got %#v", tt.query, qp)
			}
		})
	}
}
