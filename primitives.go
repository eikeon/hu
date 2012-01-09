package hu

import (
	"big"
)


func lambda(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	parameters := Tuple([]Term{nil, terms[0]})
	//parameters := Tuple([]Term{nil, Tuple([]Term{terms[0]})})
	term = terms[1]
	return Abstraction{parameters, term}
}

func operator(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	parameters := terms[0]
	term = terms[1]
	return Abstraction{parameters, term}
}

func add_numbers(environment *Environment, term Term) Term {
	var result = big.NewInt(0)
	for _, argument := range(term.(Tuple)) {
		num := environment.evaluate(argument).(*Number)
		result.Add(result, num.value)
	}
	return &Number{result}
}

func add_numbersP(environment *Environment) Term {
	var result = big.NewInt(0)
	numbers := environment.Get(Symbol("numbers"))
	for _, number := range(numbers.(Tuple)) {
		num := environment.evaluate(number).(*Number)
		result.Add(result, num.value)
	}
	return &Number{result}
}

func add_lists(environment *Environment, arguments Term) Term {
	var terms []Term
	for _, argument := range(arguments.(Tuple)) {
		for _, term := range(environment.evaluate(argument).(Tuple)) {
			terms = append(terms, term)
		}
	}
	return Tuple(terms)
}

func subtract_proc(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	// TODO: implement uniary negation
	num := environment.evaluate(terms[0]).(*Number)
	result := big.NewInt(0).Set(num.value)
	for _, argument := range(terms[1:]) {
		num = environment.evaluate(argument).(*Number)
		result.Sub(result, num.value)
	}
	return &Number{result}
}

func multiply_proc(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	var result = big.NewInt(1)
	for _, argument := range(terms) {
		result.Mul(result, argument.(*Number).value)
	}
	return &Number{result}
}

func quotient_proc(environment *Environment, term Term) Term {
        terms := term.(Tuple)
	a := environment.evaluate(terms[0]).(*Number)
	b := environment.evaluate(terms[1]).(*Number)
	result := big.NewInt(0).Quo(a.value, b.value)
	return &Number{result}
}

func remainder_proc(environment *Environment, term Term) Term {
        terms := term.(Tuple)
	a := environment.evaluate(terms[0]).(*Number)
	b := environment.evaluate(terms[1]).(*Number)
	result := big.NewInt(0).Rem(a.value, b.value)
	return &Number{result}
}

func is_number_equal_proc(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	value := terms[0].(*Number).value
	for _, argument := range(terms[1:]) {
		num := environment.evaluate(argument).(*Number)
		if value.Cmp(num.value) != 0 {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func is_less_than_proc(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	num := environment.evaluate(terms[0])
	previous := num.(*Number).value
	for _, argument := range(terms[1:]) {
		num = environment.evaluate(argument)
		next := num.(*Number).value
		if previous.Cmp(next) == -1 {
			previous = next
		} else {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func is_greater_than_proc(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	num := environment.evaluate(terms[0])
	previous := num.(*Number).value
	for _, argument := range(terms[1:]) {
		num = environment.evaluate(argument)
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

	terms := term.(Tuple)

	switch v := terms[0].(type) {
	case Symbol:
		variable = v
		value = terms[1]
	case Tuple:
		variable = v[0].(Symbol)
		parameters := v[1]
		body := terms[1]
		value = lambda(environment, Tuple([]Term{parameters, body}))
	default:
		panic("unexpected type")

	}
	environment.Define(variable, value)
	return nil
}

func set(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	variable := terms[0]
	value := terms[1]
	environment.Set(variable.(Symbol), value)
	return nil
}

func begin(environment *Environment, term Term) Term {
	var result Term
	for _, expression := range(term.(Tuple)) {
		result = environment.evaluate(expression)
	}
	return result
}

 func and(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	for _, exp := range(terms) {
		result := environment.evaluate(exp).(Boolean)
		if !result {
			return result
		}
	}
	return Boolean(true)
}

func or(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	for _, exp := range(terms) {
		result := environment.evaluate(exp).(Boolean)
		if result {
			return result
		}
	}
	return Boolean(false)
}

func ifPrimitive(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	if_predicate := terms[0]
	if environment.evaluate(if_predicate).(Boolean) {
		if_consequent := terms[1]
 		term = if_consequent
 	} else {
 		var if_alternative Term
		if len(terms) < 3 {
 			if_alternative = Boolean(false)
 		} else {
			if_alternative = terms[2]
 		}
		term = if_alternative
 	}
 	return environment.evaluate(term)
}

func apply(environment *Environment, term Term) Term {
	return Application{term}
}

func evalPrimitive(environment *Environment, term Term) Term {
	return environment.evaluate(term.(Tuple)[0])
}

func let(environment *Environment, term Term) Term {
	terms := term.(Tuple)
	bindings := terms[0]
	body := terms[1]

	var parameters, arguments Tuple
	for _, binding := range(bindings.(Tuple)) {
		b := binding.(Tuple)
		parameters = append(parameters, b[0])
		arguments = append(arguments, b[1])
	}

	parameters = Tuple([]Term{parameters}) // TODO: ??

	operator := lambda(environment, Tuple([]Term{parameters, body}))
	operands := arguments

	return Application{Tuple([]Term{operator, operands})}
}
