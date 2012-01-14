// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hu

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
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
	itemOpenParenthesis
	itemCloseParenthesis
	itemOpenCurlyBrace
	itemCloseCurlyBrace
	itemQuote
	itemSpace
	itemPeriod
	itemPageBreak
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:            "error",
	itemEOF:              "EOF",
	itemNewline:          "newline",
	itemWord:             "word",
	itemPunctuation:      "punctuation",
	itemNumber:           "number",
	itemString:           "string",
	itemOpenParenthesis:  "(",
	itemCloseParenthesis: ")",
	itemOpenCurlyBrace:   "{",
	itemCloseCurlyBrace:  "}",
	itemQuote:            "'",
	itemSpace:            "space",
	itemPeriod:           "period",
	itemPageBreak:        "page break",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = rune(-1)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name      string         // the name of the input; used only for error reports.
	input     io.RuneScanner // the string being scanned.
	current   bytes.Buffer
	state     stateFn   // the next lexing function to enter
	width     int       // width of last rune read from input.
	items     chan item // channel of scanned items.
	token     [2]item   // two-token lookahead for parser.
	peekCount int
}

// next returns the next rune in the input.
func (l *lexer) next() (rune rune) {
	rune, size, err := l.input.ReadRune()
	if err == nil {
		l.width, _ = l.current.WriteRune(rune)
		if size != l.width {
			fmt.Println("size: ", size, "width: ", l.width)
		}
		return rune
	} else {
		l.width = 0
		return eof
	}
	panic("")
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	if l.width > 0 {
		err := l.input.UnreadRune()
		if err != nil {
			fmt.Println("err: ", err)
		}
		l.current.Truncate(l.current.Len() - l.width)
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.current.String()}
	l.current.Reset()
	l.width = 0
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

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItemFromInput() item {
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

// next returns the next item taking into account peek
func (l *lexer) nextItem() item {
	if l.peekCount > 0 {
		l.peekCount--
	} else {
		l.token[0] = l.nextItemFromInput()
	}
	return l.token[l.peekCount]
}

// backup backs the input stream up one token.
func (l *lexer) backupItem() {
	l.peekCount++
}

// peek returns but does not consume the next token.
func (l *lexer) peekItem() item {
	if l.peekCount > 0 {
		return l.token[l.peekCount-1]
	}
	l.peekCount = 1
	l.token[0] = l.nextItemFromInput()

	return l.token[0]
}

// lex creates a new scanner for the input string.
func lex(name string, input io.RuneScanner) *lexer {
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
	case r == '"':
		return lexQuote
	case r == '+' || r == '-':
		rr := l.peek()
		if isPunctuation(rr) {
			l.emit(itemWord)
			return lexItem
		} else if '0' <= rr && rr <= '9' {
			return lexNumber
		} else {
			l.emit(itemPunctuation)
			return lexItem
		}
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case isPunctuation(r):
		l.backup()
		return lexPunctuation
	default:
		l.backup()
		return lexWord
	}
	return nil
}

// lexWord scans an alphanumeric
func lexWord(l *lexer) stateFn {
top:
	switch r := l.next(); {
	case isPunctuation(r):
		l.backup()
		l.emit(itemWord)
	default:
		goto top
	}
	return lexItem
}

func lexPunctuation(l *lexer) stateFn {
	switch l.next() {
	case '\n':
		l.emit(itemNewline)
	case '(':
		l.emit(itemOpenParenthesis)
	case ')':
		l.emit(itemCloseParenthesis)
	case '{':
		l.emit(itemOpenCurlyBrace)
	case '}':
		l.emit(itemCloseCurlyBrace)
	case '\'':
		l.emit(itemQuote)
	case ' ':
		l.emit(itemSpace)
	case '.':
		l.emit(itemPeriod)
	case '\f':
		l.emit(itemPageBreak)
	case ',', ':', ';', '!', '-':
		l.emit(itemPunctuation)
	default:
		l.emit(itemPunctuation)
		//panic("???")
	}
	return lexItem
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary.  This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.current.String())
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
func isPunctuation(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\f', '(', ')', '{', '}', '\'', '-', eof:
		return true
	case '.', '!', ',', ':':
		return true
	}
	return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
