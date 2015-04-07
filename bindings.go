package hu

import "strings"

func AddPrimitive(environment Environment, name string, function PrimitiveFunction) {
	environment.Define(Symbol(name), function)
}

func AddDefaultBindings(environment Environment) {
	environment.Define("true", Boolean(true))
	environment.Define("false", Boolean(false))

	AddPrimitive(environment, "lambda", lambda)
	AddPrimitive(environment, "operator", operator)

	AddPrimitive(environment, "=", is_number_equal_proc)
	AddPrimitive(environment, "<", is_less_than_proc)
	AddPrimitive(environment, ">", is_greater_than_proc)

	AddPrimitive(environment, "+", add_numbers)
	environment.Define("add_numbers", Primitive(add_numbersP))
	AddPrimitive(environment, "concat", add_lists)

	AddPrimitive(environment, "-", subtract_proc)

	AddPrimitive(environment, "*", multiply_proc)
	//AddPrimitive(environment, "quotient", quotient_proc)
	//AddPrimitive(environment, "remainder", remainder_proc)

	AddPrimitive(environment, "define", define)
	AddPrimitive(environment, "set", set)
	AddPrimitive(environment, "get", get)
	AddPrimitive(environment, "begin", begin)
	AddPrimitive(environment, "if", ifPrimitive)
	AddPrimitive(environment, "and", and)
	AddPrimitive(environment, "or", or)
	AddPrimitive(environment, "apply", apply)
	AddPrimitive(environment, "eval", evalPrimitive)
	AddPrimitive(environment, "let", let)

	Evaluate(environment, Read(strings.NewReader(`{define plus {operator ((lhs) (rhs)) {+ lhs rhs}}}
{define plus_list_operator {operator (lhs rhs) {concat lhs rhs}}} {1 2 plus 3 4}}
	`)))
}
