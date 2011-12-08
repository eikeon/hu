package hu

import (
	"testing"
	"strings"
	"fmt"
	"big"
)

type testCase struct {
	name        string
	input       string
	is_expected func(Object) bool
}

func is_eq_number(result Object, number int64) bool {
	return is_number(result) && result.(*NumberObject).value.Cmp(big.NewInt(number)) == 0
}

func TestInterpreter(t *testing.T) {
	var tests = []testCase{
		{"primitive function test", "{add 1 2}", func(result Object) bool {
			return is_eq_number(result, 3)
		}},
		{"apply test", "{apply add 1 2}", func(result Object) bool {
			return is_eq_number(result, 3)
		}},
		{"begin test", "{begin {define foo 1} foo}", func(result Object) bool {
			return is_number(result)
		}},
		{"lambda call test", "{{lambda (x y) (1 (add 3 4) x)} 4 5}", func(result Object) bool {
			_, ok := result.(*PairObject)
			if !ok {
				return false
			}
			return true
		}},
		{"primitive procedure test", "{+ 1 2}", func(result Object) bool {
			return is_eq_number(result, 3)
		}},
		{"closure test", "{begin {define (double x) {+ x x}} {double 5}}}", func(result Object) bool {
			return is_eq_number(result, 10)
		}},
		{"closure test setup", "{begin {define (fib n) {if {< n 2} n {+ {fib {- n 1}} {fib {- n 2}}}}} {fib 15}}", func(result Object) bool {
			return is_eq_number(result, 610)
		}},
		{"let test", "{let ((x 2)) {+ x x}}", func(result Object) bool {
			return is_eq_number(result, 4)
		}},
		{"quotient test", "{quotient 10 3}", func(result Object) bool {
			return is_eq_number(result, 3)
		}},
		{"remainder test", "{remainder 5 3}", func(result Object) bool {
			return is_eq_number(result, 2)
		}},

	}
	for _, test := range tests {
		fmt.Println("Running: ", test.name)
		interpreter := NewInterpreter()
		interpreter.AddDefaultBindings()
		reader := strings.NewReader(test.input)
		expression := interpreter.Read(reader)
		result := interpreter.Evaluate(expression)
		if !test.is_expected(result) {
			t.Error("unexpected result", result, " for ", test.name)
		}
	}
}
