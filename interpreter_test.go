package hu

import (
	"testing"
	"strings"
	"fmt"
)

type testCase struct {
	name        string
	input       string
	is_expected func(Object) bool
}

func foo(result Object) bool { return false }

func TestInterpreter(t *testing.T) {
	var tests = []testCase{
		{"primitive function test", "{add (1 2)}", func(result Object) bool {
			if is_number(result) && result.(*NumberObject).value == 3 {
				return true
			}
			return false
		}},
		{"apply test", "{apply (add 1 2)}", func(result Object) bool {
			if is_number(result) && result.(*NumberObject).value == 3 {
				return true
			}
			return false
		}},
		{"begin test", "{begin ((define foo 1) foo)}", func(result Object) bool {
			return is_number(result)
		}},
		{"lambda call test", "{(lambda (x) (1 (add 3 4) x)) (4)}", func(result Object) bool {
			_, ok := result.(*PairObject)
			if !ok {
				return false
			}
			return true
		}},
		{"primitive procedure test", "{+ (1 2)}", func(result Object) bool {
			if is_number(result) && result.(*NumberObject).value == 3 {
				return true
			}
			return false
		}},
		{"closure test", "{begin ((define (double x) (+ x x)) (double 5))}", func(result Object) bool {
			if is_number(result) && result.(*NumberObject).value == 10 {
				return true
			}
			return false
		}},
		{"closure test setup", "{begin ((define (fib n) (if (< n 2) n (+ (fib (- n 1)) (fib (- n 2))))) (fib 15))}", func(result Object) bool {
			if is_number(result) && result.(*NumberObject).value == 610 {
				return true
			}
			return false
		}},
		//         {"closure test", "(fib 15)", func(result Object) bool {
		// 		if is_number(result) && result.(*NumberObject).value==610 {
		// 			return true
		// 		}
		// 		return false
		// }},
	}
	for _, test := range tests {
		fmt.Println("Running: ", test.name)
		interpreter := NewInterpreter()
		interpreter.AddDefaultBindings()
		// interpreter.AddPrimitive("add", (*Interpreter).add_proc)
		// interpreter.AddPrimitive("apply", (*Interpreter).apply)
		// interpreter.AddPrimitive("begin", (*Interpreter).begin)
		// interpreter.AddPrimitive("lambda", (*Interpreter).lambda)

		// interpreter.AddPrimitiveProcedure("+", (*Interpreter).add_proc)

		reader := strings.NewReader(test.input)
		expression := interpreter.Read(reader)
		result := interpreter.Evaluate(expression)
		if !test.is_expected(result) {
			t.Error("unexpected result", result, " for ", test.name)
		}
	}
}
