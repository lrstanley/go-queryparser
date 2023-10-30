// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package queryparser

import (
	"strings"
	"unicode"
)

// Options allow the adjustment of allowed filters and characters.
type Options struct {
	// CutFn allows excluding specific characters from being allowed within
	// the filter fields. When the function returns true on a rune, it will
	// be excluded from the filter field.
	CutFn func(rune) bool
	// Allowed is a slice of allowed filter names. If no allowed filter names
	// are provided, all are considered allowed.
	Allowed []string
}

// Parser represents a parser.
type Parser struct {
	s   *scanner
	opt *Options

	buf []tokenRef
}

// New returns a new instance of Parser. Make sure Parser.Parser() is called
// or this will leak goroutines.
func New(query string, opt Options) *Parser {
	return &Parser{s: newScanner(query), opt: &opt}
}

// Parse is a higher level helper method to return a query from a query string.
func Parse(query string) *Query {
	return New(query, Options{CutFn: DefaultCut}).Parse()
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tr tokenRef) {
	// If we have a token on the buffer, then return it.
	if len(p.buf) > 0 {
		// Pop first item off the buffer.
		tr = p.buf[0]
		copy(p.buf, p.buf[1:])
		p.buf = p.buf[:len(p.buf)-1]
		return tr
	}

	// Otherwise read the next token from the scanner.
	tr = p.s.nextToken()
	return tr
}

// unscan pushes provided token/literal back onto the buffer.
func (p *Parser) unscan(tr tokenRef) {
	p.buf = append(p.buf, tr)
}

// accept scans if the provided token matches, otherwise unscans.
func (p *Parser) accept(tok token) bool {
	tr := p.scan()
	if tr.tok == tok {
		return true
	}

	p.unscan(tr)
	return false
}

// Parse parses the input query and returns a new instance of Query if there
// were no errors.
func (p *Parser) Parse() *Query {
	defer p.s.drain()

	qp := &Query{Filters: make(map[string][]string)}

	for {
		tr := p.scan()

		switch tr.tok {
		case tokenIDENT:
			p.scanField(tr, qp)
		case tokenEOF:
			if p.opt.CutFn != nil {
				qp.Raw = cutsetFunc(qp.Raw, p.opt.CutFn)
				qp.Raw = stripDuplicateWS(qp.Raw)
			}
			return qp
		default:
			qp.Raw += tr.lit
		}
	}
}

func (p *Parser) scanField(ident tokenRef, qp *Query) {
	if !isIdent(ident.lit) {
		qp.Raw += ident.lit
		return
	}

	// Return early if it's not allowed.
	if p.opt.Allowed != nil && len(p.opt.Allowed) > 0 {
		var in bool
		for i := 0; i < len(p.opt.Allowed); i++ {
			if strings.EqualFold(p.opt.Allowed[i], ident.lit) {
				in = true
				break
			}
		}
		if !in {
			qp.Raw += ident.lit
			return
		}
	}

	delim := p.scan()
	if delim.tok != tokenDELIM {
		qp.Raw += ident.lit
		p.unscan(delim)
		return
	}

	// Chomp all trailing fields.
	var fields []tokenRef
	var count int
	for {
		field := p.scan()
		count++

		if field.tok == tokenFIELD || field.tok == tokenIDENT {
			fields = append(fields, field)
			continue
		}

		if count == 1 {
			qp.Raw += ident.lit
			p.unscan(delim)
			p.unscan(field)
			return
		}
		p.unscan(field)
		break
	}

	// Chomp trailing whitespaces if there are any.
	_ = p.accept(tokenWS)

	var fieldText string
	for i := 0; i < len(fields); i++ {
		fieldText += fields[i].lit
	}

	if p.opt.CutFn != nil {
		qp.Add(ident.lit, cutsetFunc(fieldText, p.opt.CutFn))
		return
	}

	qp.Add(ident.lit, fieldText)
}

func cutsetFunc(input string, cutFn func(rune) bool) (out string) {
	for _, c := range input {
		if !cutFn(c) {
			out += string(c)
		}
	}
	return out
}

// DefaultCut is the default cut function, which allowed stripping out potentially
// unwanted characters from filter fields and raw text. Only allows
// " _,-.:A-Za-z0-9" (or unicode equivalents).
func DefaultCut(r rune) (strip bool) {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != ' ' && r != '\t' &&
		r != '_' && r != ',' && r != '-' && r != '.' && r != ':'
}

func stripDuplicateWS(val string) string {
	for strings.Contains(val, "  ") {
		val = strings.ReplaceAll(val, "  ", " ")
	}

	return val
}
