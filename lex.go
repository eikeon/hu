// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hu

import (
	"fmt"
	"strings"
	"unicode"
	"utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemNewline     //
	itemWord        // alphanumeric word
	itemNumber      // simple number, including imaginary
	itemPunctuation //
	itemString      // quoted string (includes quotes)
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:       "error",
	itemEOF:         "EOF",
	itemNewline:     "newline",
	itemWord:        "word",
	itemPunctuation: "punctuation",
	itemNumber:      "number",
	itemString:      "string",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // the name of the input; used only for error reports.
	input string    // the string being scanned.
	state stateFn   // the next lexing function to enter
	pos   int       // current position in the input.
	start int       // start position of this item.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// next returns the next rune in the input.
func (l *lexer) next() (rune int) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() int {
	rune := l.next()
	l.backup()
	return rune
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		state: lexItem,
		items: make(chan item, 2), // Two items of buffering is sufficient for all state functions
	}
	return l
}

// state functions

// lexItem ...
func lexItem(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
		//return l.errorf("unexpected EOF")
	case r == '"':
		return lexQuote
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	//case isAlphaNumeric(r):
	case isPunctuation(r) == false:
		l.backup()
		return lexWord
	default:
		return l.errorf("unrecognized character in Item: %#U", r)
	}
	return nil //lexItem
}

// lexWord scans an alphanumeric
func lexWord(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(itemEOF)
			return nil
		//case isAlphaNumeric(r):
		//	// absorb.
		case isPunctuation(r):
			l.backup()
			//word := l.input[l.start:l.pos]
			switch {
			default:
				l.emit(itemWord)
			}
			break Loop
		default:
			//
		}

	}
	return lexPunctuation
}

func lexPunctuation(l *lexer) stateFn {
	switch {
	case l.accept("\n"):
		l.emit(itemNewline)
		return lexPunctuation
	case l.accept(" ,:;.!"):
		l.emit(itemPunctuation)
		return lexPunctuation
	default:
		//l.backup()
		return lexItem
	}
	return lexPunctuation
	// switch r := l.next(); {
	// case r == eof:
	// 	l.emit(itemEOF)
	// 	return nil
	// default:
	// 	l.backup()
	// }

	// if !l.scanPunctuation() {
	// 	return l.errorf("bad punctuation: %q", l.input[l.start:l.pos])
	// }
	// l.emit(itemPunctuation)
	// return lexItem
}

func (l *lexer) scanPunctuation() bool {
	l.acceptRun(" ,;.!\n")
	return true
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary.  This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexPunctuation
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("/") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// lexQuote scans a quoted string.
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	return lexPunctuation
}

// isPunctuation reports whether r is a punctuation character.
func isPunctuation(r int) bool {
	switch r {
	case ' ', '\t', '\n', '\r', ',', ':':
		return true
	}
	return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r int) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
