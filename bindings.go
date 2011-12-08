package hu

func (interpreter *Interpreter) AddDefaultBindings() {
	interpreter.AddPrimitive("add", (*Interpreter).add_proc)

	interpreter.AddPrimitive("quote", (*Interpreter).quote)
	interpreter.AddPrimitive("unquote", (*Interpreter).unquote)
	interpreter.AddPrimitive("define", (*Interpreter).define)
	interpreter.AddPrimitive("set!", (*Interpreter).set)
	interpreter.AddPrimitive("lambda", (*Interpreter).lambda)
	interpreter.AddPrimitive("begin", (*Interpreter).begin)
	interpreter.AddPrimitive("if", (*Interpreter).ifPrimitive)
	interpreter.AddPrimitive("and", (*Interpreter).and)
	interpreter.AddPrimitive("or", (*Interpreter).or)
	interpreter.AddPrimitive("apply", (*Interpreter).apply)
	interpreter.AddPrimitive("eval", (*Interpreter).evalPrimitive)
	interpreter.AddPrimitive("let", (*Interpreter).let)

	interpreter.AddPrimitiveProcedure("null?", (*Interpreter).is_null_proc)
	interpreter.AddPrimitiveProcedure("boolean?", is_type_proc_tor(is_boolean))
	interpreter.AddPrimitiveProcedure("symbol?", is_type_proc_tor(is_symbol))
	interpreter.AddPrimitiveProcedure("integer?", is_type_proc_tor(is_number))
	interpreter.AddPrimitiveProcedure("char?", is_type_proc_tor(is_character))
	interpreter.AddPrimitiveProcedure("string?", is_type_proc_tor(is_string))
	interpreter.AddPrimitiveProcedure("pair?", is_type_proc_tor(is_pair))
	interpreter.AddPrimitiveProcedure("eof-object?", is_type_proc_tor(is_eof_object))

	interpreter.AddPrimitiveProcedure("+", (*Interpreter).add_proc)
	interpreter.AddPrimitiveProcedure("-", (*Interpreter).subtract_proc)
	interpreter.AddPrimitiveProcedure("*", (*Interpreter).multiply_proc)
	interpreter.AddPrimitiveProcedure("quotient" , (*Interpreter).quotient_proc)
	interpreter.AddPrimitiveProcedure("remainder", (*Interpreter).remainder_proc)
	interpreter.AddPrimitiveProcedure("=", (*Interpreter).is_number_equal_proc)
	interpreter.AddPrimitiveProcedure("<", (*Interpreter).is_less_than_proc)
	interpreter.AddPrimitiveProcedure(">", (*Interpreter).is_greater_than_proc)

	interpreter.AddPrimitiveProcedure("cons", (*Interpreter).cons_proc)
	interpreter.AddPrimitiveProcedure("car", (*Interpreter).car_proc)
	interpreter.AddPrimitiveProcedure("cdr", (*Interpreter).cdr_proc)
	interpreter.AddPrimitiveProcedure("set-car!", (*Interpreter).set_car_proc)
	interpreter.AddPrimitiveProcedure("set-cdr!", (*Interpreter).set_cdr_proc)
	interpreter.AddPrimitiveProcedure("list", (*Interpreter).list_proc)

	interpreter.AddPrimitiveProcedure("eq?", (*Interpreter).is_eq_proc)

	interpreter.AddPrimitiveProcedure("read", (*Interpreter).read_proc)

	interpreter.AddPrimitiveProcedure("write-char", (*Interpreter).write_char_proc)
	interpreter.AddPrimitiveProcedure("write", (*Interpreter).write_proc)

	interpreter.AddPrimitiveProcedure("error", (*Interpreter).error_proc)
}
