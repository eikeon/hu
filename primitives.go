package hu

import (
	"os"
	"bufio"
	"fmt"
	"big"
)

func (interpreter *Interpreter) is_null_proc(arguments Object, environment *Environment) (result Object) {
	if car(arguments) == nil {
		result = TRUE
	} else {
		result = FALSE
	}
	return
}

func is_type_proc_tor(predicate func(Object) bool) func(*Interpreter, Object, *Environment) Object {
	return func(interpreter *Interpreter, arguments Object, environment *Environment) (result Object) {
		if predicate(car(arguments)) {
			result = TRUE
		} else {
			result = FALSE
		}
		return
	}
}

func (interpreter *Interpreter) add_proc(arguments Object, environment *Environment) Object {
	var result = big.NewInt(0)

	for arguments != nil {
		result.Add(result, car(arguments).(*NumberObject).value)
		arguments = cdr(arguments)
	}
	return &NumberObject{result}
}

func (interpreter *Interpreter) subtract_proc(arguments Object, environment *Environment) Object {
	// TODO: implement uniary negation
	result := big.NewInt(0).Set(car(arguments).(*NumberObject).value)
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		result.Sub(result, car(arguments).(*NumberObject).value)
	}
	return &NumberObject{result}
}

func (interpreter *Interpreter) multiply_proc(arguments Object, environment *Environment) Object {
	var result = big.NewInt(1)

	for arguments != nil {
		result.Mul(result, car(arguments).(*NumberObject).value)
		arguments = cdr(arguments)
	}
	return &NumberObject{result}
}

func (interpreter *Interpreter) quotient_proc(arguments Object, environment *Environment) Object {
	result := big.NewInt(0).Quo(car(arguments).(*NumberObject).value, car(cdr(arguments)).(*NumberObject).value)
	return &NumberObject{result}
}

func (interpreter *Interpreter) remainder_proc(arguments Object, environment *Environment) Object {
	result := big.NewInt(0).Rem(car(arguments).(*NumberObject).value, car(cdr(arguments)).(*NumberObject).value)
	return &NumberObject{result}
}

func (interpreter *Interpreter) is_number_equal_proc(arguments Object, environment *Environment) Object {
	value := car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		if value.Cmp(car(arguments).(*NumberObject).value) != 0 {
			return FALSE
		}
	}
	return TRUE
}

func (interpreter *Interpreter) is_less_than_proc(arguments Object, environment *Environment) Object {
	previous := car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		next := car(arguments).(*NumberObject).value
		if previous.Cmp(next) == -1 {
			previous = next
		} else {
			return FALSE
		}
	}
	return TRUE
}

func (interpreter *Interpreter) is_greater_than_proc(arguments Object, environment *Environment) Object {
	previous := car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		next := car(arguments).(*NumberObject).value
		if previous.Cmp(next) == 1 {
			previous = next
		} else {
			return FALSE
		}
	}
	return TRUE
}

func (interpreter *Interpreter) cons_proc(arguments Object, environment *Environment) Object {
	return cons(car(arguments), car(cdr(arguments)))
}

func (interpreter *Interpreter) car_proc(arguments Object, environment *Environment) Object {
	return car(car(arguments))
}

func (interpreter *Interpreter) cdr_proc(arguments Object, environment *Environment) Object {
	return cdr(car(arguments))
}

func (interpreter *Interpreter) set_car_proc(arguments Object, environment *Environment) Object {
	set_car(car(arguments), car(cdr(arguments)))
	return nil
}

func (interpreter *Interpreter) set_cdr_proc(arguments Object, environment *Environment) Object {
	set_cdr(car(arguments), car(cdr(arguments)))
	return nil
}

func (interpreter *Interpreter) list_proc(arguments Object, environment *Environment) Object {
	return arguments
}

func (interpreter *Interpreter) is_eq_proc(arguments Object, environment *Environment) (result Object) {
	obj1 := car(arguments)
	obj2 := car(cdr(arguments))

	// TODO: try switch t1, t2 := ...
	switch t1 := obj1.(type) {
	case *NumberObject:
		t2, ok := obj2.(*NumberObject)
		if ok && t1.value.Cmp(t2.value) == 0 {
			result = TRUE
		}
		break
	case RuneObject:
		t2, ok := obj2.(RuneObject)
		if ok && t1 == t2 {
			result = TRUE
		}
		break
	case *StringObject:
		t2, ok := obj2.(*StringObject)
		if ok && t1.value == t2.value {
			result = TRUE
		}
		break
	default:
		if obj1 == obj2 {
			result = TRUE
		} else {
			result = FALSE
		}
	}
	return
}

func (interpreter *Interpreter) write_char_proc(arguments Object, environment *Environment) Object {
	out := os.Stdout
	character := car(arguments)
	arguments = cdr(arguments)
	if arguments != nil {
		panic("TODO") // out :=
	}
	fmt.Fprintf(out, "%s", character)
	return nil
}

func (interpreter *Interpreter) write_proc(arguments Object, environment *Environment) Object {
	out := os.Stdout
	exp := car(arguments)
	arguments = cdr(arguments)
	if arguments != nil {
		panic("TODO")
	}
	fmt.Fprintf(out, "%s", exp)
	return nil
}

func (interpreter *Interpreter) read_proc(arguments Object, environment *Environment) (result Object) {
	in := bufio.NewReader(os.Stdin)
	if arguments != nil {
		panic("TODO")
	}
	result = interpreter.Read(in)
	if result == nil {
		result = eof_object
	}
	return
}

func (interpreter *Interpreter) error_proc(arguments Object, environment *Environment) Object {
	out := os.Stderr
	for ; arguments != nil; arguments = cdr(arguments) {
		fmt.Fprintf(out, "%s", car(arguments))
		fmt.Fprintf(out, " ")
	}
	panic("exiting")
}
