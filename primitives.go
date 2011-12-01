package hu

import (
	"os"
	"bufio"
	"fmt"
)

func is_null_proc(arguments Object) (result Object) {
	if car(arguments) == nil {
		result = TRUE
	} else {
		result = FALSE
	}
	return
}

func is_type_proc_tor(predicate func(Object) bool) func(Object) Object {
	return func(arguments Object) (result Object) {
		if predicate(car(arguments)) {
			result = TRUE
		} else {
			result = FALSE
		}
		return
	}
}

func add_proc(arguments Object) Object {
	var result int64 = 0

	for arguments != nil {
		result += car(arguments).(*NumberObject).value
		arguments = cdr(arguments)
	}
	return &NumberObject{result}
}

func subtract_proc(arguments Object) Object {
	// TODO: implement uniary negation
	result := car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		result -= car(arguments).(*NumberObject).value

	}
	return &NumberObject{result}
}

func multiply_proc(arguments Object) Object {
	var result int64 = 1

	for arguments != nil {
		result *= car(arguments).(*NumberObject).value
		arguments = cdr(arguments)
	}
	return &NumberObject{result}
}

func is_number_equal_proc(arguments Object) Object {
	value := car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		if value != car(arguments).(*NumberObject).value {
			return FALSE
		}
	}
	return TRUE
}

func is_less_than_proc(arguments Object) Object {
	var previous, next int64

	previous = car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		next = car(arguments).(*NumberObject).value
		if previous < next {
			previous = next
		} else {
			return FALSE
		}
	}
	return TRUE
}

func is_greater_than_proc(arguments Object) Object {
	var previous, next int64

	previous = car(arguments).(*NumberObject).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		next = car(arguments).(*NumberObject).value
		if previous > next {
			previous = next
		} else {
			return FALSE
		}
	}
	return TRUE
}

func cons_proc(arguments Object) Object {
	return cons(car(arguments), car(cdr(arguments)))
}

func car_proc(arguments Object) Object {
	return car(car(arguments))
}

func cdr_proc(arguments Object) Object {
	return cdr(car(arguments))
}

func set_car_proc(arguments Object) Object {
	set_car(car(arguments), car(cdr(arguments)))
	return nil
}

func set_cdr_proc(arguments Object) Object {
	set_cdr(car(arguments), car(cdr(arguments)))
	return nil
}

func list_proc(arguments Object) Object {
	return arguments
}

func is_eq_proc(arguments Object) (result Object) {
	obj1 := car(arguments)
	obj2 := car(cdr(arguments))

	// TODO: try switch t1, t2 := ...
	switch t1 := obj1.(type) {
	case *NumberObject:
		t2, ok := obj2.(*NumberObject)
		if ok && t1.value == t2.value {
			result = TRUE
		}
		break
	case *RuneObject:
		t2, ok := obj2.(*RuneObject)
		if ok && t1.rune == t2.rune {
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

func write_char_proc(arguments Object) Object {
	out := os.Stdout
	character := car(arguments)
	arguments = cdr(arguments)
	if arguments != nil {
		panic("TODO") // out :=
	}
	fmt.Fprintf(out, "%s", character)
	return nil
}

func write_proc(arguments Object) Object {
	out := os.Stdout
	exp := car(arguments)
	arguments = cdr(arguments)
	if arguments != nil {
		panic("TODO")
	}
	fmt.Fprintf(out, "%s", exp)
	return nil
}

func (interpreter *Interpreter) read_proc(arguments Object) (result Object) {
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

func error_proc(arguments Object) Object {
	out := os.Stderr
	for ; arguments != nil; arguments = cdr(arguments) {
		fmt.Fprintf(out, "%s", car(arguments))
		fmt.Fprintf(out, " ")
	}
	panic("exiting")
}
