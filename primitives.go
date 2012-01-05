package hu

import (
	"big"
)


func lambda(environment *Environment, term Term) Term {
	parameters := &Pair{nil, &Pair{car(term), nil}}
	term = car(cdr(term))
	return Abstraction{parameters, term}
}

func operator(environment *Environment, term Term) Term {
	parameters := car(term)
	term = car(cdr(term))
	return Abstraction{parameters, term}
}

func add_numbers(environment *Environment, arguments Term) Term {
	var result = big.NewInt(0)
	for arguments != nil {
		num := environment.evaluate(car(arguments)).(*Number)
		result.Add(result, num.value)
		arguments = cdr(arguments)
	}
	return &Number{result}
}

func add_numbersP(environment *Environment) Term {
	var result = big.NewInt(0)
	numbers := environment.Get(Symbol("numbers"))
	for numbers != nil {
		num := environment.evaluate(car(numbers)).(*Number)
		result.Add(result, num.value)
		numbers = cdr(numbers)
	}
	return &Number{result}
}

func add_lists(environment *Environment, arguments Term) Term {
	var pairs []*Pair
	for arguments != nil {
		pairs = append(pairs, environment.evaluate(car(arguments)).(*Pair))
		arguments = cdr(arguments)
	}
	return concat(pairs...)
}

func subtract_proc(environment *Environment, arguments Term) Term {
	// TODO: implement uniary negation
	num := environment.evaluate(car(arguments)).(*Number)
	result := big.NewInt(0).Set(num.value)
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		num = environment.evaluate(car(arguments)).(*Number)
		result.Sub(result, num.value)
	}
	return &Number{result}
}

func multiply_proc(environment *Environment, arguments Term) Term {
	var result = big.NewInt(1)

	for arguments != nil {
		result.Mul(result, car(arguments).(*Number).value)
		arguments = cdr(arguments)
	}
	return &Number{result}
}

func quotient_proc(environment *Environment, arguments Term) Term {
	a := environment.evaluate(car(arguments)).(*Number)
	b := environment.evaluate(car(cdr(arguments))).(*Number)
	result := big.NewInt(0).Quo(a.value, b.value)
	return &Number{result}
}

func remainder_proc(environment *Environment, arguments Term) Term {
	a := environment.evaluate(car(arguments)).(*Number)
	b := environment.evaluate(car(cdr(arguments))).(*Number)
	result := big.NewInt(0).Rem(a.value, b.value)
	return &Number{result}
}

func is_number_equal_proc(environment *Environment, arguments Term) Term {
	value := car(arguments).(*Number).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		num := environment.evaluate(car(arguments)).(*Number)
		if value.Cmp(num.value) != 0 {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func is_less_than_proc(environment *Environment, arguments Term) Term {
	num := environment.evaluate(car(arguments))
	previous := num.(*Number).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		num = environment.evaluate(car(arguments))
		next := num.(*Number).value
		if previous.Cmp(next) == -1 {
			previous = next
		} else {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func is_greater_than_proc(environment *Environment, arguments Term) Term {
	num := environment.evaluate(car(arguments))
	previous := num.(*Number).value
	for arguments = cdr(arguments); arguments != nil; arguments = cdr(arguments) {
		num = environment.evaluate(car(arguments))
		next := num.(*Number).value
		if previous.Cmp(next) == 1 {
			previous = next
		} else {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func define(environment *Environment, term Term) Term {
	var variable Symbol
	var value Term

	switch v := car(term).(type) {
	case Symbol:
		variable = v
		value = car(cdr(term))
	default:
		variable = car(car(term)).(Symbol)
		parameters := cdr(car(term))
		body := cdr(term)
		value = lambda(environment, &Pair{parameters, body})
	}
	environment.Define(variable, value)
	return nil
}

func set(environment *Environment, term Term) Term {
	variable := car(term)
	value := car(cdr(term))
	environment.Set(variable.(Symbol), value)
	return nil
}

func begin(environment *Environment, term Term) Term {
	var result Term
	for expressions := term; expressions != nil; expressions = cdr(expressions) {
		expression := car(expressions)
		result = environment.evaluate(expression)
	}
	return result
}

func and(environment *Environment, term Term) Term {
	tests := term
	for exp := tests; exp != nil; exp = cdr(exp) {
		first_exp := car(exp)
		result := environment.evaluate(first_exp).(Boolean)
		if !result {
			return result
		}
	}
	return Boolean(true)
}

func or(environment *Environment, term Term) Term {
	tests := term
	for exp := tests; exp != nil; exp = cdr(exp) {
		first_exp := car(exp)
		result := environment.evaluate(first_exp).(Boolean)
		if result {
			return result
		}
	}
	return Boolean(false)
}

func ifPrimitive(environment *Environment, term Term) Term {
	if_predicate := car(term)
	if environment.evaluate(if_predicate).(Boolean) {
		if_consequent := car(cdr(term))
 		term = if_consequent
 	} else {
 		var if_alternative Term
 		if cdr(cdr(term)) == nil {
 			if_alternative = Boolean(false)
 		} else {
 			if_alternative = car(cdr(cdr(term)))
 		}
 		term = if_alternative
 	}
 	return environment.evaluate(term)
}

func apply(environment *Environment, term Term) Term {
	return Application{term}
}

func evalPrimitive(environment *Environment, term Term) Term {
	return environment.evaluate(car(term))
}

func let(environment *Environment, term Term) Term {
	bindings := car(term)
	body := cdr(term)

	binding_parameter := func(binding Term) Term { return car(binding) }
	parameters := list_from(bindings, binding_parameter)

	binding_arguments := func(binding Term) Term { return car(cdr(binding)) }
	arguments := list_from(bindings, binding_arguments)

	operator := lambda(environment, &Pair{parameters, body})
	operands := arguments

	return Application{&Pair{operator, operands}}
}
