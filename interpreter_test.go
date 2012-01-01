package hu

import (
	"testing"
	"strings"
	"big"
)

type testCase struct {
	input       string
	is_expected func(Term) bool
}

func is_eq_number(number int64) func(Term) bool {
	return func(result Term) bool {
		num, ok := result.(*Number)
		return ok && num.value.Cmp(big.NewInt(number)) == 0
	}
}

func is_eq(expected Term) func(Term) bool {
	return func(result Term) bool {
		return result==expected
	}
}

func is_pair() func(Term) bool {
	return func(result Term) bool {
		_, ok := result.(*Pair)
		return ok
	}
}

func is_unbound() func(Term) bool {
	return func(result Term) bool {
		_, ok := result.(UnboundVariableError)
		return ok
	}
}

func is_nil() func(Term) bool {
	return func(result Term) bool {
		return result == nil
	}
}

var tests = []testCase{
	{"true", is_eq(Boolean(true))},
	{"{+ 1 2}", is_eq_number(3)},
	{"{+ 1 2}", is_eq_number(3)},
	{"{+ 1 {+ 2 3}}", is_eq_number(6)},
	{"{concat (1 2) (3 4)}", is_pair()},
	{"{begin {define foo 1} foo}", is_eq_number(1)},
	{"{{lambda (x y) {+ x y}} 4 5}", is_eq_number(9)},
	{"{begin {define (double x) {+ x x}} {double 5}}}", is_eq_number(10)},
	{"{begin {define n 1} {define c 5} {{lambda (n) {+ n c}} {+ n 1}}}", is_eq_number(7)},
	{"{begin {define fib {lambda (n) {if {< n 2} n {+ {fib {- n 1}} {fib {- n 2}}}}}} {fib 15}}",	is_eq_number(610)},
	{"{begin {define plus {operator ((lhs) (rhs)) {+ lhs rhs}}} {1 plus 2}}", is_eq_number(3)},
	{"{begin {define plus {operator (lhs rhs) {concat lhs rhs}}} {1 2 plus 3 4}}", is_pair()},
	{"{apply + 1 2}", is_eq_number(3)},
	{"{let ((x 2)) {+ x x}}", is_eq_number(4)},
	{"{quotient 10 3}", is_eq_number(3)},
	{"{remainder 5 3}", is_eq_number(2)},
	{"foo", is_unbound()},
	{"{+ 1 foo}", is_unbound()},
}

func TestInterpreter(t *testing.T) {
	for _, test := range tests {
		environment := NewEnvironment()
		AddDefaultBindings(environment)
		reader := strings.NewReader(test.input)
		expression := Read(reader)
		result := environment.Evaluate(expression)
		if test.is_expected(result) {
			t.Logf("  PASS: %v resulted in %v as expected", test.input, result)
		} else {
			t.Errorf("  FAIL: %v unexpectedly resulted in %v", test.input, result)
		}
	}
}
