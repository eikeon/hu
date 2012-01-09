package hu

import (
	"io"
	"big"
	"strings"
	"fmt"
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
		case itemEOF:
			return nil
		case itemPunctuation:
		case itemNewline:
		case itemSpace:
		case itemError:
			panic(token.val)
		default:
			panic(token.typ)
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
