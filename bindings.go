package hu

func AddDefaultBindings(environment *Environment) {
	environment.Define("true", Boolean(true))
	environment.Define("false", Boolean(false))

	environment.AddPrimitive("=", is_number_equal_proc)
	environment.AddPrimitive("<", is_less_than_proc)
	environment.AddPrimitive(">", is_greater_than_proc)
	environment.AddPrimitive("+", add_proc)
	environment.AddPrimitive("-", subtract_proc)

	environment.AddPrimitive("*", multiply_proc)
	environment.AddPrimitive("quotient", quotient_proc)
	environment.AddPrimitive("remainder", remainder_proc)

	environment.AddPrimitive("define", define)
	environment.AddPrimitive("set", set)
	environment.AddPrimitive("lambda", lambda)
	environment.AddPrimitive("begin", begin)
	environment.AddPrimitive("if", ifPrimitive)
	environment.AddPrimitive("and", and)
	environment.AddPrimitive("or", or)
	environment.AddPrimitive("apply", apply)
	environment.AddPrimitive("eval", evalPrimitive)
	environment.AddPrimitive("let", let)
}
