package hu

func (interpreter *Interpreter) quote(object Object, environment *Environment) Object {
	return cons(quote_symbol, object)
}

func (interpreter *Interpreter) unquote(object Object, environment *Environment) Object {
	// TODO check that car is quote
	return cdr(object)
}

func (interpreter *Interpreter) evalList(list Object, environment *Environment) Object {
	eval := func(object Object) Object { return interpreter.evaluate(object, environment) }
	return list_from(list, eval)
}

func (interpreter *Interpreter) define(object Object, environment *Environment) Object {
	var variable, value Object

	if is_symbol(car(object)) {
		variable = car(object)
		value = car(cdr(object))
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

func (interpreter *Interpreter) lambda(object Object, outer *Environment) Object {
	parameters := car(object)
	function := cdr(object)
	f := func(interpreter *Interpreter, object Object, environment *Environment) Object {
		operands := interpreter.evalList(object, environment)
		environment = outer.Extend(parameters, operands)
		function := interpreter.evalList(function, environment)
		return interpreter.begin(function, environment)
	}
	return &PrimitiveFunctionObject{f}
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
	result := TRUE
	tests := object
	for exp := tests; exp != nil; exp = cdr(exp) {
		first_exp := car(exp)
		result = interpreter.evaluate(first_exp, environment)
		if is_false(result) {
			return result
		}
	}
	return result
}

func (interpreter *Interpreter) or(object Object, environment *Environment) Object {
	result := FALSE
	tests := object
	for exp := tests; exp != nil; exp = cdr(exp) {
		first_exp := car(exp)
		result = interpreter.evaluate(first_exp, environment)
		if is_true(result) {
			return result
		}
	}
	return result
}

func (interpreter *Interpreter) ifPrimitive(object Object, environment *Environment) Object {
	if_predicate := car(object)
	if is_true(interpreter.evaluate(if_predicate, environment)) {
		if_consequent := car(cdr(object))
		object = if_consequent
	} else {
		var if_alternative Object
		if is_the_empty_list(cdr(cdr(object))) {
			if_alternative = FALSE
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
	return &ExpressionObject{operator, operands}
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
	application := cons(operator, operands)
	return interpreter.evaluate(application, environment)
}
