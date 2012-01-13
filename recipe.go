// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse builds parse trees for templates.  The grammar is defined
// in the documents for the template package.
package hu

import (
	"bytes"

	"fmt"
	"runtime"
	"strings"
)

type Recipe struct {
	Name        string
	Description string
	Ingredients []string
	Directions  []string
	Attributes  map[string]string
	Photo       string
}

func (r *Recipe) Id() string {
	return strings.ToLower(strings.Replace(r.Name, " ", "_", -1))
}

func (r Recipe) String() string {
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "%v\n\n", r.Name)
	fmt.Fprintf(buffer, "%v\n\n", r.Description)
	for _, i := range r.Ingredients {
		fmt.Fprintf(buffer, "%v\n", i)
	}
	fmt.Fprintf(buffer, "\n")
	for _, i := range r.Directions {
		fmt.Fprintf(buffer, "%v\n\n", i)
	}
	fmt.Fprintf(buffer, "\n")
	return string(buffer.Bytes())
}

// Tree is the representation of a parsed template.
type Tree struct {
	Name   string // Name is the name of the template.
	Recipe *Recipe
	//Root *ListNode // Root is the top-level root of the parse tree.
	// Parsing only; cleared after parse.
	funcs     []map[string]interface{}
	lex       *lexer
	token     [2]item // two-token lookahead for parser.
	peekCount int
	vars      []string // variables defined at the moment.
}

// next returns the next token.
func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// backup2 backs the input stream up two tokens
func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

// peek returns but does not consume the next token.
func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()

	return t.token[0]
}

// Parsing.

// New allocates a new template with the given name.
func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}

// errorf formats the error and terminates processing.
func (t *Tree) errorf(format string, args ...interface{}) {
	t.Recipe = nil
	//format = fmt.Sprintf("template: %s:%d: %s", t.Name, t.lex.lineNumber(), format)
	format = fmt.Sprintf("template: %s: %s", t.Name, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type.
func (t *Tree) expect(expected itemType, context string) item {
	token := t.next()
	if token.typ != expected {
		t.errorf("expected %s in %s; got %s", expected, context, token)
	}
	return token
}

// unexpected complains about the token and terminates processing.
func (t *Tree) unexpected(token item, context string) {
	t.errorf("unexpected %s in %s", token, context)
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.stopParse()
		}
		*errp = e.(error)
	}
	return
}

// startParse starts the template parsing from the lexer.
func (t *Tree) startParse(funcs []map[string]interface{}, lex *lexer) {
	t.Recipe = nil
	t.lex = lex
	t.vars = []string{"$"}
	t.funcs = funcs
}

// stopParse terminates parsing.
func (t *Tree) stopParse() {
	t.lex = nil
	t.vars = nil
	t.funcs = nil
}

// Parse parses the template definition string to construct an internal
// representation of the template for execution.
func (t *Tree) Parse(s string, funcs ...map[string]interface{}) (tree *Tree, err error) {
	defer t.recover(&err)
	t.startParse(funcs, lex(t.Name, strings.NewReader(s)))
	t.parse(true)
	t.stopParse()
	return t, nil
}

// parse is the helper for Parse.
// It triggers an error if we expect EOF but don't reach it.
func (t *Tree) parse(toEOF bool) { //(next Node) {
	r := &Recipe{}
	r.Attributes = make(map[string]string)
	t.Recipe = r
	r.Name = t.parseName()

	t.expect(itemNewline, "name-description-separator")
	r.Description = t.parseDescription()

	t.expect(itemNewline, "description-ingredients-separator")
	for t.peek().typ != itemNewline {
		r.Ingredients = append(r.Ingredients, t.parseIngredient())
	}

	t.expect(itemNewline, "ingredients-directions-separator")
	for {
		r.Directions = append(r.Directions, t.parseDirection())
		next := t.next().typ
		if next != itemNewline || next == itemEOF {
			break
		}
	}
	t.backup()
	t.parseAttributes()
	// for {
	// 	next := t.next()
	// 	if next.typ == itemEOF {
	// 		break
	// 	} else {
	// 		//t.errorf("unexpected %s at line %v", next, t.lex.lineNumber())			
	// 	}
	// }
}

func (t *Tree) parseName() (result string) {
	for {
		switch token := t.next(); token.typ {
		case itemNewline:
			return
		case itemError:
			t.errorf("lex: %s in title", token.val)
		default:
			result += token.val
		}
	}
	return
}

func (t *Tree) parseDescription() (result string) {
	for {
		switch token := t.next(); token.typ {
		case itemNewline:
			return
		case itemError:
			fmt.Println("??", token.typ)
			t.errorf("lex: %s in description", token.val)
		default:
			result += token.val
		}
	}
	return
}

func (t *Tree) parseIngredient() (result string) {
	for {
		switch token := t.next(); token.typ {
		case itemNewline:
			return
		case itemError:
			t.errorf("lex: %s in ingredient", token.val)
		default:
			result += token.val
		}
	}
	return
}

func (t *Tree) parseDirection() (result string) {
	for {
		switch token := t.next(); token.typ {
		case itemNewline, itemEOF:
			return
		case itemError:
			t.errorf("lex: %s in direction", token.val)
		default:
			result += token.val
		}
	}
	return
}

func (t *Tree) parseAttributes() {
	key := ""
	value := ""

	for {
		switch token := t.next(); token.typ {
		case itemWord:
			if key == "" {
				key = token.val
				if key != "" {
					t.expect(itemPunctuation, "attribute-key") // item for colon specifically?
				}
			} else {
				value += token.val
			}
		case itemNewline:
			switch key {
			case "Photo":
				t.Recipe.Photo = value
			default:
				if key != "" {
					t.Recipe.Attributes[key] = value
				}
			}
			key = ""
			value = ""
			//t.Recipe.
		case itemEOF:
			return
		case itemError:
			t.errorf("lex: %s in attribute", token.val)
		default:
			value += token.val
		}
	}
	return
}

type Line []Term

func (line Line) String() string {
	return fmt.Sprintf("Line(%v)", []Term(line))
}

type Section []Term

func (section Section) String() string {
	return fmt.Sprintf("Section(%v)", []Term(section))
}

type Page []Term

func (page Page) String() string {
	return fmt.Sprintf("Page(%v)", []Term(page))
}

type Document []Term

func (document Document) String() string {
	return fmt.Sprintf("Document(%v)", []Term(document))
}


func read_line(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemNewline:
		lexer.nextItem()
		return Line(terms)
	case itemPageBreak, itemEOF:
		return Line(terms)
	default:
		term := read(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}

func read_section(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemNewline:
		lexer.nextItem()
		return Section(terms)
	case itemPageBreak, itemEOF:
		return Section(terms)
	default:
		term := read_line(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}


func read_page(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemPageBreak, itemEOF:
		return Page(terms)
	default:
		term := read_section(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}

func read_pages(lexer *lexer) Term {
	var pages []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemPageBreak:
		lexer.nextItem()
		goto next
	case itemEOF:
		return Document(pages)
	default:
		page := read_page(lexer)
		pages = append(pages, page)
		goto next
	}
	panic("")
}
