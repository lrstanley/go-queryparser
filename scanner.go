// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package queryparser

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// stateFn represents the state of the scanner as a function that returns the
// next state.
type stateFn func(*scanner) stateFn

type token int

type tokenRef struct {
	tok token  // Token that represents the literal.
	pos int    // Position in the input.
	lit string // Literal/value of item.
}

const (
	tokenEOF token = iota // EOF.

	tokenDELIM // :
	tokenFIELD // Quoted (with spaces/words) or unquoted (single word).
	tokenIDENT // Raw text, an IDENT can also be a WORD.
	tokenWS    // Whitespaces.

	// eof isn't a token, but rather the literal reference to the EOF token.
	eof = 1
)

func (i tokenRef) String() string {
	if i.tok == tokenEOF {
		return "EOF"
	}

	if len(i.lit) > 10 {
		return fmt.Sprintf("%.10q...", i.lit)
	}

	return fmt.Sprintf("%q", i.lit)
}

// scanner represents a lexical scanner.
type scanner struct {
	items   chan tokenRef // The channel of scanned items.
	input   string        // The string being scanned.
	pos     int           // Current position in the input.
	start   int           // Start position of the acive item.
	width   int           // Width of last rune read from input.
	lastPos int           // Position of most recent item returned by nextItem.
}

// newScanner returns a new instance of Scanner. This starts a goroutine. Make sure
// to call drain() on it to ensure it doesn't leak goroutines.
func newScanner(input string) *scanner {
	s := &scanner{
		input: input,
		items: make(chan tokenRef),
	}

	go s.run()
	return s
}

// emit passes a tokenRef back to the client.
func (s *scanner) emit(t token) {
	s.items <- tokenRef{t, s.start, s.input[s.start:s.pos]}
	s.start = s.pos
}

// nextToken returns the next tokenRef from the input. Called by the parser, not
// in the lexing goroutine.
func (s *scanner) nextToken() tokenRef {
	item := <-s.items
	s.lastPos = item.pos
	return item
}

// run runs the state machine for the lexer.
func (s *scanner) run() {
	for state := scanMain; state != nil; {
		state = state(s)
	}
	close(s.items)
}

// drain drains the output so the lexing goroutine will exit. Called by the
// parser, not in the lexing goroutine.
func (s *scanner) drain() {
	for range s.items {
	}
}

// read reads the next rune from the buffered reader.
// Returns eof if an error occurs (or io.EOF is returned).
func (s *scanner) next() rune {
	if s.pos >= len(s.input) {
		s.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.width = w
	s.pos += s.width

	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (s *scanner) backup() {
	s.pos -= s.width
}

// peek steps forward one rune, reads, and backs up again.
func (s *scanner) peek() rune {
	r := s.next()
	s.backup()
	return r
}

func scanMain(s *scanner) stateFn {
	switch r := s.next(); {
	case r == eof:
		s.emit(tokenEOF)
		return nil
	case isWhitespace(r):
		return scanWhitespace
	case r == ':':
		s.emit(tokenDELIM)
		return scanMain
	case r == '"':
		return scanDoubleQuote
	case r == '\'':
		return scanSingleQuote
	case isWord(r):
		return scanWord
	}
	return nil
}

// scanWhitespace scans a run of space characters. One space has already been
// seen.
func scanWhitespace(s *scanner) stateFn {
	for isWhitespace(s.peek()) {
		s.next()
	}
	s.emit(tokenWS)
	return scanMain
}

// scanWord scans a run of word characters. One word character has already been
// seen.
func scanWord(s *scanner) stateFn {
	for isWord(s.peek()) {
		s.next()
	}
	s.emit(tokenIDENT)
	return scanMain
}

// scanSingleQuote scans a quoted string.
func scanSingleQuote(s *scanner) stateFn {
Loop:
	for {
		switch s.next() {
		case '\\':
			if r := s.next(); r != eof {
				break
			}
			fallthrough
		case eof:
			// Should this be a req? Unterminated quoted string.
			break Loop
		case '\'':
			break Loop
		}
	}
	s.emit(tokenFIELD)
	return scanMain
}

// scanDoubleQuote scans a quoted string.
func scanDoubleQuote(s *scanner) stateFn {
Loop:
	for {
		switch s.next() {
		case '\\':
			if r := s.next(); r != eof {
				break
			}
			fallthrough
		case eof:
			// Should this be a req? Unterminated quoted string.
			break Loop
		case '"':
			break Loop
		}
	}
	s.emit(tokenFIELD)
	return scanMain
}

// isWhitespace returns true if ch is a space or a tab.
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

// isWord returns true if ch is character allowed in raw text.
func isWord(ch rune) bool {
	return ch >= '!' && ch <= '~' && ch != ':'
}

func isIdent(input string) bool {
	for _, r := range input {
		if r != '_' && r != '-' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
