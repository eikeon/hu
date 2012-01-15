package hu

import (
	"fmt"
	"io"
	"math/big"
	"strings"
)

func Read(in io.RuneScanner) (result Term) {
	defer func() {
		switch x := recover().(type) {
		case Term:
			result = x
		case interface{}:
			result = Error(fmt.Sprintf("%v", x))
		}
	}()
	lexer := lex("", in)
	result = read(lexer)
	return
}

func read(lexer *lexer) Term {
	var terms []Term
	for {
		switch token := lexer.nextItem(); token.typ {
		case itemWord:
			return Symbol(token.val)
		case itemString:
			return String(strings.Trim(token.val, string("\"")))
		case itemNumber:
			num := big.NewInt(0)
			num.SetString(token.val, 10)
			return &Number{num}
		case itemOpenParenthesis:
			tuple := read_tuple(lexer)
			next := lexer.nextItem()
			if next.typ != itemCloseParenthesis {
				panic("expected )")
			}
			return tuple
		case itemCloseParenthesis:
		case itemOpenCurlyBrace:
			expression := read_tuple(lexer)
			next := lexer.nextItem()
			if next.typ != itemCloseCurlyBrace {
				panic("expected }")
			}
			return Application{expression}
		case itemCloseCurlyBrace:
		// case itemQuote:
		// 	return &Tuple{Symbol("quote"), read(lexer)}
		case itemPageBreak:
			lexer.backupItem()
			return nil
		case itemEOF:
			lexer.backupItem()
			if terms == nil {
				return nil
			} else {
				return Tuple(terms)
			}
		case itemPunctuation:
			return Symbol(token.val)
		case itemNewline:
		case itemSpace:
			return Symbol(token.val)
		case itemError:
			return Symbol("?")
			//panic(token.val)
		case itemSection:
			terms = append(terms, read_section(lexer))
		default:
			return Symbol(token.val)
		}
	}
	panic("unexpectedly reached this point")
}

func read_tuple(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemCloseParenthesis, itemCloseCurlyBrace:
		return Tuple(terms)
	case itemSpace:
		lexer.nextItem()
		goto next
	default:
		term := read(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}

func read_part(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemNewline:
		lexer.nextItem()
		return Part(terms)
	case itemSection, itemEOF:
		return Part(terms)
	default:
		term := read_line(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}

func read_line(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemNewline:
		lexer.nextItem()
		return Line(terms)
	case itemSection, itemEOF:
		return Line(terms)
	default:
		term := read(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}

/*
 Use § for sections (instead of  for pages)
 Use blank lines to separate items in sections
 Use \t for nesting (for example subsections)
 Use :\n\t as continuations for named items
 */

func read_section(lexer *lexer) Term {
	var terms []Term
next:
	switch token := lexer.peekItem(); token.typ {
	case itemSection, itemEOF:
		return Tuple(terms)
	default:
		term := read_part(lexer)
		terms = append(terms, term)
		goto next
	}
	panic("")
}
