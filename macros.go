package hu

func (interpreter *Interpreter) quote(object Object, environment *Environment) Object {
	return cons(quote_symbol, object)
}

func (interpreter *Interpreter) unquote(object Object, environment *Environment) Object {
	// TODO check that car is quote
	return cdr(object)
}

func (interpreter *Interpreter) define(object Object, environment *Environment) Object {
	var variable, value Object

	if is_symbol(car(object)) {
		variable = car(object)
		value = interpreter.evaluate(car(cdr(object)), environment)
	} else {
		variable = car(car(object))
		parameters := cdr(car(object))
		body := cdr(object)
		value = interpreter.lambda(cons(parameters, body), environment)
	}

	environment.Define(variable, value)
	return nil
}

func (interpreter *Interpreter) set(object Object, environment *Environment) Object {
	variable := car(object)
	value := car(cdr(object))
	value = interpreter.evaluate(value, environment)
	environment.Set(variable, value)
	return nil
}

func (interpreter *Interpreter) lambda(object Object, environment *Environment) Object {
	parameters := car(object)
	object = cdr(object)
	return &Abstraction{parameters, object, environment}
}

func (interpreter *Interpreter) begin(object Object, environment *Environment) Object {
	var result Object
	for expressions := object; expressions != nil; expressions = cdr(expressions) {
		expression := car(expressions)
		result = interpreter.evaluate(expression, environment)
	}
	return result
}

func (interpreter *Interpreter) and(object Object, environment *Environment) Object {
	tests := object
	for exp := tests; exp != nil; exp = cdr(exp) {
		first_exp := car(exp)
		result := interpreter.evaluate(first_exp, environment)
		if is_false(result) {
			return result
		}
	}
	return BooleanObject(true)
}

func (interpreter *Interpreter) or(object Object, environment *Environment) Object {
	tests := object
	for exp := tests; exp != nil; exp = cdr(exp) {
		first_exp := car(exp)
		result := interpreter.evaluate(first_exp, environment)
		if is_true(result) {
			return result
		}
	}
	return BooleanObject(false)
}

func (interpreter *Interpreter) ifPrimitive(object Object, environment *Environment) Object {
	if_predicate := car(object)
	if is_true(interpreter.evaluate(if_predicate, environment)) {
		if_consequent := car(cdr(object))
		object = if_consequent
	} else {
		var if_alternative Object
		if is_the_empty_list(cdr(cdr(object))) {
			if_alternative = BooleanObject(false)
		} else {
			if_alternative = car(cdr(cdr(object)))
		}
		object = if_alternative
	}
	return interpreter.evaluate(object, environment)
}

func (interpreter *Interpreter) apply(object Object, environment *Environment) Object {
	operator := car(object)
	operands := cdr(object)
	return &Application{operator, operands}
}

func (interpreter *Interpreter) evalPrimitive(object Object, environment *Environment) Object {
	expression := car(object)
	environment = car(cdr(object)).(*Environment)
	return interpreter.evaluate(expression, environment)
}

func (interpreter *Interpreter) let(object Object, environment *Environment) Object {
	bindings := car(object)
	body := cdr(object)

	binding_parameter := func(binding Object) Object { return car(binding) }
	parameters := list_from(bindings, binding_parameter)

	binding_arguments := func(binding Object) Object { return car(cdr(binding)) }
	arguments := list_from(bindings, binding_arguments)

	operator := interpreter.lambda(cons(parameters, body), environment)
	operands := arguments

	return &Application{operator, operands}
}
