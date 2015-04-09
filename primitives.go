package hu

import (
	"math/big"

	"fmt"
	"log"
)

func lambda(environment Environment, term Term) Term {
	terms := term.(Tuple)
	parameters := Tuple([]Term{nil, terms[0]})
	//parameters := Tuple([]Term{nil, Tuple([]Term{terms[0]})})
	term = terms[1]
	return Abstraction{parameters, term}
}

func operator(environment Environment, term Term) Term {
	terms := term.(Tuple)
	parameters := terms[0]
	term = terms[1]
	return Abstraction{parameters, term}
}

func add_numbers(environment Environment, term Term) Term {
	var result = big.NewRat(0, 1)
	for i, argument := range Evaluate(environment, term).(Tuple) {
		num, ok := Evaluate(environment, argument).(*Number)
		if ok {
			result.Add(result, num.value)
		} else {
			error := fmt.Sprintf("argument %d ( %s ) to add_numbers not a number", i, argument)
			return Error(error)
		}
	}
	return &Number{result}
}

func add_numbersP(environment Environment) Term {
	var result = big.NewRat(0, 1)
	numbersExp, _ := environment.Get(Symbol("numbers"))
	numbers := Evaluate(environment, numbersExp)
	for _, number := range numbers.(Tuple) {
		num := Evaluate(environment, number).(*Number)
		result.Add(result, num.value)
	}
	return &Number{result}
}

func add_lists(environment Environment, arguments Term) Term {
	var terms []Term
	for _, argument := range arguments.(Tuple) {
		for _, term := range Evaluate(environment, argument).(Tuple) {
			terms = append(terms, term)
		}
	}
	return Tuple(terms)
}

func subtract_proc(environment Environment, term Term) Term {
	terms := term.(Tuple)
	// TODO: implement uniary negation
	num := Evaluate(environment, terms[0]).(*Number)
	result := big.NewRat(0, 1).Set(num.value)
	for _, argument := range terms[1:] {
		num = Evaluate(environment, argument).(*Number)
		result.Sub(result, num.value)
	}
	return &Number{result}
}

func multiply_proc(environment Environment, term Term) Term {
	terms := term.(Tuple)
	var result = big.NewRat(1, 1)
	log.Println(fmt.Sprintf("mult %#v\n", term))
	for _, argument := range terms {
		log.Println(fmt.Sprintf("'%#v'", argument))
		result.Mul(result, argument.(*Number).value)
	}
	return &Number{result}
}

func quotient_proc(environment Environment, term Term) Term {
	terms := term.(Tuple)
	a := Evaluate(environment, terms[0]).(*Number)
	b := Evaluate(environment, terms[1]).(*Number)
	result := big.NewRat(0, 1).Quo(a.value, b.value)
	return &Number{result}
}

// func remainder_proc(environment Environment, term Term) Term {
// 	terms := term.(Tuple)
// 	a := Evaluate(environment, terms[0]).(*Number)
// 	b := Evaluate(environment, terms[1]).(*Number)
// 	result := big.NewRat(0, 1).Rem(a.value, b.value)
// 	return &Number{result}
// }

func is_number_equal_proc(environment Environment, term Term) Term {
	terms := term.(Tuple)
	value := terms[0].(*Number).value
	for _, argument := range terms[1:] {
		num := Evaluate(environment, argument).(*Number)
		if value.Cmp(num.value) != 0 {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func is_less_than_proc(environment Environment, term Term) Term {
	terms := term.(Tuple)
	num := Evaluate(environment, terms[0])
	previous := num.(*Number).value
	for _, argument := range terms[1:] {
		num = Evaluate(environment, argument)
		next := num.(*Number).value
		if previous.Cmp(next) == -1 {
			previous = next
		} else {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func is_greater_than_proc(environment Environment, term Term) Term {
	terms := term.(Tuple)
	num := Evaluate(environment, terms[0])
	previous := num.(*Number).value
	for _, argument := range terms[1:] {
		num = Evaluate(environment, argument)
		next := num.(*Number).value
		if previous.Cmp(next) == 1 {
			previous = next
		} else {
			return Boolean(false)
		}
	}
	return Boolean(true)
}

func define(environment Environment, term Term) Term {
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
		//TODO: value = Closure{value, environment}
	default:
		panic("unexpected type")

	}
	environment.Define(variable, value)
	return nil
}

func variable(environment Environment, term Term) Term {
	// schedule () {lambda (newSchedule) {runSchedule newSchedule}}
	terms := term.(Tuple)

	name, ok := terms[0].(Symbol)
	if !ok {
		return Error("unexpected type for name")
	}
	didSet, ok := Evaluate(environment, terms[1]).(Abstraction)
	if !ok {
		return Error("unexpected type for didSet")
	}
	environment.Define(name, &Property{name, didSet})
	environment.Define(Symbol(name+"^didSet"), didSet)
	return nil
}

func set(environment Environment, term Term) Term {
	terms := term.(Tuple)
	variable := terms[0]
	value := Evaluate(environment, terms[1])
	name := variable.(Symbol)
	environment.Set(name, value)
	didSet, ok := environment.Get(Symbol(name + "^didSet"))
	log.Println("didSet", didSet, ok)
	if ok {
		Evaluate(environment, Application([]Term{didSet, value}))
	}
	return nil
}

func get(environment Environment, term Term) Term {
	terms := term.(Tuple)
	variable := terms[0]
	value, ok := environment.Get(variable.(Symbol))
	if ok {
		return value
	} else {
		return UnboundVariableError{variable, "get"}
	}
}

func begin(environment Environment, term Term) Term {
	var result Term
	for _, expression := range term.(Tuple) {
		result = Evaluate(environment, expression)
	}
	return result
}

func and(environment Environment, term Term) Term {
	terms := term.(Tuple)
	for _, exp := range terms {
		result := Evaluate(environment, exp).(Boolean)
		if !result {
			return result
		}
	}
	return Boolean(true)
}

func or(environment Environment, term Term) Term {
	terms := term.(Tuple)
	for _, exp := range terms {
		result := Evaluate(environment, exp).(Boolean)
		if result {
			return result
		}
	}
	return Boolean(false)
}

func ifPrimitive(environment Environment, term Term) Term {
	terms := term.(Tuple)
	if_predicate := terms[0]
	if Evaluate(environment, if_predicate).(Boolean) {
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
	return Evaluate(environment, term)
}

func apply(environment Environment, term Term) Term {
	return Application(term.(Tuple))
}

func evalPrimitive(environment Environment, term Term) Term {
	return Evaluate(environment, term.(Tuple)[0])
}

func let(environment Environment, term Term) Term {
	terms := term.(Tuple)
	bindings := terms[0]
	body := terms[1]

	var parameters, arguments Tuple
	for _, binding := range bindings.(Tuple) {
		b := binding.(Tuple)
		parameters = append(parameters, b[0])
		arguments = append(arguments, b[1])
	}

	parameters = Tuple([]Term{parameters}) // TODO: ??

	operator := lambda(environment, Tuple([]Term{parameters, body}))
	operands := arguments

	return Application([]Term{operator, operands})
}
