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
				"this": []string{"1"}, "is": []string{"foo"},
			},
		}},
		{name: "Filters but only 'this' allowed", args: args{raw: "this:1 is:foo a test", allowed: []string{"this"}}, wantQry: Query{
			Orig:   "this:1 is:foo a test",
			Raw:    "is:foo a test",
			Pretty: "this:1 is:foo a test",
			filters: map[string][]string{
				"this": []string{"1"},
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
