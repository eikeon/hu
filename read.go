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
			pair := read_pair(lexer)
			next := lexer.nextItem()
			if next.typ != itemCloseParenthesis {
				panic("expected )")
			}
			return pair
		case itemCloseParenthesis:
		case itemOpenCurlyBrace:
			expression := read_pair(lexer)
			next := lexer.nextItem()
			if next.typ != itemCloseCurlyBrace {
				panic("expected }")
			}
			return Application{expression}
		case itemCloseCurlyBrace:
		case itemQuote:
			return &Pair{Symbol("quote"), read(lexer)}
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

func read_pair(lexer *lexer) Term {
	var car_term, cdr_term Term
	for {
		switch token := lexer.peekItem(); token.typ {
		case itemCloseParenthesis, itemCloseCurlyBrace:
			return nil
		case itemSpace:
			lexer.nextItem()
			break
		default:
			car_term = read(lexer)
			goto done_car
		}
	}
done_car:
	for {
		switch token := lexer.peekItem(); token.typ {
		case itemSpace:
			lexer.nextItem()
		case itemPeriod:
			lexer.nextItem()
			cdr_term = read(lexer)
			goto done_cdr
		default:
			cdr_term = read_pair(lexer)
			goto done_cdr
		}
	}
done_cdr:
	return &Pair{car_term, cdr_term}
}
