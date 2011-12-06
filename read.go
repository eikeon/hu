package hu

import (
	"fmt"
	"io"
	"unicode"
)

func is_delimiter(rune int) bool {
	return unicode.IsSpace(rune) || rune == 0 ||
		rune == '(' || rune == ')' ||
		rune == '"' || rune == ';'
}

func is_initial(rune int) bool {
	return unicode.IsLetter(rune) || rune == '*' || rune == '/' || rune == '>' ||
		rune == '<' || rune == '=' || rune == '?' || rune == '!'
}

func getc(in io.RuneReader) int {
	rune, size, err := in.ReadRune()
	if size == 0 {
	}
	if err != nil {
		return 0
	}
	return rune
}

func ungetc(in io.RuneScanner) {
	err := in.UnreadRune()
	if err != nil {
		fmt.Println("err:", err)
	}
}

func peek(in io.RuneScanner) int {
	c := getc(in)
	ungetc(in)
	return c
}

func ignoreWhitespace(in io.RuneScanner) {
	for c := getc(in); c != 0; c = getc(in) {
		if unicode.IsSpace(c) {
		} else if c == '\n' {
		} else if c == ';' {
			for c = getc(in); c != 0 && c != '\n'; c = getc(in) {
			}
			continue
		} else {
			ungetc(in)
			break
		}
	}
}

func ignore_expected_string(in io.RuneScanner, expected string) {
	for _, c := range expected {
		rune := getc(in)
		if c != rune {
			panic(fmt.Sprintf("didn't find all of expected string '%s'", expected))
		}
	}
}

func peek_expected_delimiter(in io.RuneScanner) {
	if !is_delimiter(peek(in)) {
		panic("character not an expected delimiter\n")
	}
}

func (interpreter *Interpreter) Read(in io.RuneScanner) Object {
	ignoreWhitespace(in)

	rune := getc(in)

	if rune == '#' {
		rune = getc(in)
		switch rune {
		case 't':
			return TRUE
		case 'f':
			return FALSE
		case '\\':
			rune = getc(in)

			switch rune {
			case 0:
				panic("incomplete character literal\n")
			case 's':
				if peek(in) == 'p' {
					ignore_expected_string(in, "pace")
					peek_expected_delimiter(in)
					return &RuneObject{' '}
				}
				break
			case 'n':
				if peek(in) == 'e' {
					ignore_expected_string(in, "ewline")
					peek_expected_delimiter(in)
					return &RuneObject{'\n'}
				}
				break
			default:
				peek_expected_delimiter(in)
				return &RuneObject{rune}
			}
		default:
			panic("unknown boolean or character literal\n")
		}
	} else if unicode.IsDigit(rune) ||
		(rune == '-' && (unicode.IsDigit(peek(in)))) {

		var sign int64 = 1
		var num int64 = 0
		if rune == '-' {
			sign = -1
		} else {
			ungetc(in)
		}
		for rune = getc(in); unicode.IsDigit(rune); rune = getc(in) {
			num = (num * 10) + int64(rune-'0')
		}
		num *= sign
		if is_delimiter(rune) {
			ungetc(in)
			return &NumberObject{num}
		} else {
			panic("number not followed by delimiter\n")
		}
	} else if is_initial(rune) ||
		((rune == '+' || rune == '-') &&
			is_delimiter(peek(in))) {

		var buffer []int
		for is_initial(rune) || unicode.IsDigit(rune) || rune == '+' || rune == '-' {
			buffer = append(buffer, rune)
			rune = getc(in)
		}
		if is_delimiter(rune) {
			ungetc(in)
			return Symbol(string(buffer))
		} else {
			panic("symbol not followed by delimiter.")
		}
	} else if rune == '"' {
		var buffer []int
		for rune = getc(in); rune != '"'; rune = getc(in) {
			if rune == '\\' {
				rune = getc(in)
				if rune == 'n' {
					rune = '\n'
				}
			}
			if rune == 0 {
				panic("non-terminated string literal\n")
			}
			buffer = append(buffer, rune)
		}
		return &StringObject{string(buffer)}
	} else if rune == '(' {
		return interpreter.read_pair(in)
	} else if rune == '\'' {
		return cons(quote_symbol, interpreter.Read(in))
	} else if rune == 0 || rune == 10 {
		return nil
	} else {
		panic(fmt.Sprintf("unexpected input: %v", rune))
	}
	panic("unexpectedly reached this point")
}

func (interpreter *Interpreter) read_pair(in io.RuneScanner) Object {

	var car_object, cdr_object Object

	ignoreWhitespace(in)

	rune := getc(in)
	if rune == ')' {
		return nil
	}
	ungetc(in)

	car_object = interpreter.Read(in)

	ignoreWhitespace(in)

	rune = getc(in)
	if rune == '.' {
		rune = peek(in)
		if !is_delimiter(rune) {
			panic("dot not followed by delimiter\n")
		}
		cdr_object = interpreter.Read(in)
		ignoreWhitespace(in)
		rune = getc(in)
		if rune != ')' {
			panic("where was the trailing right paren?\n")
		}
	} else {
		ungetc(in)
		cdr_object = interpreter.read_pair(in)
	}
	return cons(car_object, cdr_object)
}
