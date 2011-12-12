package hu

import (
	"io"
	"big"
	"strings"
)

func (interpreter *Interpreter) Read(in io.RuneScanner) Object {
	lexer := lex("", in)
	return interpreter.read(lexer)
}

func (interpreter *Interpreter) read(lexer *lexer) Object {
	for {
		switch token := lexer.nextItem(); token.typ {
		case itemWord:
			return Symbol(token.val)
		case itemString:
			return &StringObject{strings.Trim(token.val, string("\""))}
		case itemNumber:
			num := big.NewInt(0)
			num.SetString(token.val, 10)
			return &NumberObject{num}
		case itemOpenParenthesis:
			pair := interpreter.read_pair(lexer)
			next := lexer.nextItem()
			if next.typ != itemCloseParenthesis {
				panic("expected )")
			}
			return pair
		case itemCloseParenthesis:
		case itemOpenCurlyBrace:
			expression := interpreter.read_pair(lexer)
			next := lexer.nextItem()
			if next.typ != itemCloseCurlyBrace {
				panic("expected }")
			}
			return &ExpressionObject{car(expression), cdr(expression)}
		case itemCloseCurlyBrace:
		case itemQuote:
			return cons(quote_symbol, interpreter.read(lexer))
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

func (interpreter *Interpreter) read_pair(lexer *lexer) Object {
	var car_object, cdr_object Object
	for {
		switch token := lexer.peekItem(); token.typ {
		case itemCloseParenthesis, itemCloseCurlyBrace:
			return nil
		case itemSpace:
			lexer.nextItem()
			break
		default:
			car_object = interpreter.read(lexer)
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
			cdr_object = interpreter.read(lexer)
			goto done_cdr
		default:
			cdr_object = interpreter.read_pair(lexer)
			goto done_cdr
		}
	}
done_cdr:
	return cons(car_object, cdr_object)
}
