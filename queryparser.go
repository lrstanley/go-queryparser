// Package queryparser parses a common "q" http GET variable to strip out
// filters, which can be used for advanced searching. For example:
//    Hello World tags:example,world author:lrstanley
//
// go-queryparser also strips out all not very safe, or useful characters.
// This means only " _,-.:A-Za-z0-9" are allowed.
package queryparser

import "strings"

// Query contains everything necessary to pull specific filters, and
// everything else, that originated from the original query string.
type Query struct {
	// Orig is the original query string.
	Orig string
	// Raw is everything besides filters (with commands stripped out.)
	Raw string
	// Pretty is much like Orig, however it has the invalid characters that
	// may have been contained within Origin, removed.
	Pretty string
	// filters are the "filters" or "tags" within the query string. This
	// allows multiple filters of the same name, and by default does not
	// overwrite previously defined filters.
	filters map[string][]string
}

// Get returns the results of the filter if it exists, and if it successfully
// found a result.
func (q *Query) Get(name string) (results []string, ok bool) {
	results, ok = q.filters[name]

	return results, ok
}

// GetOne returns the last known result for the filter, if it exists. Useful
// if you only want a user to define a filter once. The resulting string
// is empty if no filter of that name was found.
func (q *Query) GetOne(name string) string {
	results, ok := q.Get(name)

	if !ok {
		return ""
	}

	out := results[len(results)-1]

	return out
}

// IsZero checks to see if the query string parsed result is of a zero value,
// indicating that the end user supplied an invalid search query.
func (q *Query) IsZero() bool {
	if q.Pretty == "" {
		return true
	}

	if q.Raw == "" && len(q.filters) == 0 {
		return true
	}

	return false
}

// Strip strips all invalid characters from the search query.
func Strip(raw string) (out string) {
	for i := 0; i < len(raw); i++ {
		if raw[i] == 0x20 || raw[i] == 0x5F || // Space and underscore.
			(raw[i] >= 0x2C && raw[i] <= 0x2E) || // Comma, dash, and period.
			(raw[i] >= 0x30 && raw[i] <= 0x3A) || // 0-9, and colon.
			(raw[i] >= 0x41 && raw[i] <= 0x5A) || // A-Z.
			(raw[i] >= 0x61 && raw[i] <= 0x7A) { // a-z.
			out += string(raw[i])
		}
	}

	return out
}

// New parses the query string, and splits it up as necessary.
func New(raw string, allowed []string) (qry Query) {
	qry.Orig = raw
	qry.Pretty = qry.Orig

	// Remove all un-needed characters.
	qry.Pretty = strings.TrimSpace(Strip(qry.Pretty))

	// Strip all double spaces from the query.
	// E.g. "something     else" -> "something else".
	for {
		if strings.Contains(qry.Pretty, "  ") {
			qry.Pretty = strings.Replace(qry.Pretty, "  ", " ", -1)
			continue
		}

		break
	}

	qry.filters = make(map[string][]string)

	if qry.Pretty == "" {
		return qry
	}

	var j int
	var isin bool
	var name, tmp string
	var values []string

	for i := 0; i < len(qry.Pretty); i++ {
		j = strings.Index(qry.Pretty[i:], " ")
		if j < 0 {
			// Assume it's nearing the end of the line, or there are no special
			// filters.
			tmp += qry.Pretty[i:]
			break
		}

		c := strings.Index(qry.Pretty[i:i+j], ":")

		// Assume it's random, e.g. "something ".
		if c < 0 {
			tmp += qry.Pretty[i : i+j+1]
			i += j
			continue
		}

		// If it finds a "<name>:<value> " match.
		name = strings.ToLower(qry.Pretty[i : i+c])

		// Check to see if we allow the filter.
		if len(allowed) > 0 {
			isin = false
			for _, ftype := range allowed {
				// Assume it's a formatted filter that we don't allow.
				if strings.ToLower(ftype) == name {
					isin = true
					break
				}
			}

			// Just trunk it and move on.
			if !isin {
				tmp += qry.Pretty[i : i+j+1]
				i += j
				continue
			}
		}

		if _, ok := qry.filters[name]; !ok {
			qry.filters[name] = []string{}
		}

		values = strings.Split(qry.Pretty[i+c+1:i+j], ",")

		qry.filters[name] = append(qry.filters[name], values...)
		i += j
	}

	// Strip out commas.
	for i := 0; i < len(tmp); i++ {
		if tmp[i] != 0x2C {
			qry.Raw += string(tmp[i])
		}
	}
	qry.Raw = strings.TrimSpace(qry.Raw)

	return qry
}
