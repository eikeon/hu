package hu

import (
	"fmt"
	"io"
	"math/big"
	"strings"
)

func read(lexer *lexer) (term Term) {
	switch token := lexer.nextItem(); token.typ {
	case itemString:
		term = String(strings.Trim(token.val, string("\"")))
	case itemNumber:
		num := big.NewRat(0, 1)
		num.SetString(token.val)
		term = &Number{num}
	case itemOpenParenthesis:
		reader := &partReader{ignore: 1<<itemSpace | 1<<itemCloseParenthesis, end: 1 << itemCloseParenthesis}
		term = Tuple(reader.read(lexer))
	case itemOpenCurlyBrace:
		reader := &partReader{ignore: 1<<itemSpace | 1<<itemCloseCurlyBrace, end: 1 << itemCloseCurlyBrace}
		term = Application{Tuple(reader.read(lexer))}
	case itemEOF:
		lexer.backupItem()
		term = nil
	case itemError:
		term = Error(token.val)
	default:
		term = Symbol(token.val)
	}
	return
}

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

type partReader struct {
	ignore, start, end uint64
	sub                *partReader
}

func (reader *partReader) read(lexer *lexer) (result Part) {
	first := true
	for {
		token := lexer.peekItem()
		if token.typ == itemEOF {
			break
		}
		if first == false {
			if 1<<token.typ&reader.start != 0 {
				break
			}
		} else {
			first = false
		}
		if 1<<token.typ&reader.ignore != 0 {
			lexer.nextItem()
		} else {
			var t Term
			if reader.sub == nil || 1<<token.typ&reader.end != 0 {
				t = read(lexer)
			} else {
				t = reader.sub.read(lexer)
			}
			result = append(result, t)
		}
		if 1<<token.typ&reader.end != 0 {
			break
		}
	}
	return
}

func ReadDocument(in io.RuneScanner) Part {
	line := &partReader{end: 1 << itemNewline}
	part := &partReader{end: 1 << itemNewline, sub: line}
	section := &partReader{start: 1 << itemSection, sub: part}
	document := &partReader{sub: section}
	return document.read(lex("", in))
}
