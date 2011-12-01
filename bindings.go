package hu

func (interpreter *Interpreter) AddDefaultBindings() {
	interpreter.AddMacro("quote", (*Interpreter).quote)
	interpreter.AddMacro("define", (*Interpreter).define)
	interpreter.AddMacro("set!", (*Interpreter).set)
	interpreter.AddMacro("lambda", (*Interpreter).lambda)
	interpreter.AddMacro("begin", (*Interpreter).begin)
	interpreter.AddMacro("if", (*Interpreter).ifMacro)
	interpreter.AddMacro("and", (*Interpreter).and)
	interpreter.AddMacro("or", (*Interpreter).or)
	interpreter.AddMacro("apply", (*Interpreter).apply)
	interpreter.AddMacro("eval", (*Interpreter).evalMacro)
	interpreter.AddMacro("let", (*Interpreter).let)

	interpreter.AddPrimative("null?", is_null_proc)
	interpreter.AddPrimative("boolean?", is_type_proc_tor(is_boolean))
	interpreter.AddPrimative("symbol?", is_type_proc_tor(is_symbol))
	interpreter.AddPrimative("integer?", is_type_proc_tor(is_number))
	interpreter.AddPrimative("char?", is_type_proc_tor(is_character))
	interpreter.AddPrimative("string?", is_type_proc_tor(is_string))
	interpreter.AddPrimative("pair?", is_type_proc_tor(is_pair))
	interpreter.AddPrimative("eof-object?", is_type_proc_tor(is_eof_object))

	interpreter.AddPrimative("+", add_proc)
	interpreter.AddPrimative("-", subtract_proc)
	interpreter.AddPrimative("*", multiply_proc)
	// TODO:
	// interpreter.AddPrimative("quotient" , quotient_proc)
	// interpreter.AddPrimative("remainder", remainder_proc)
	interpreter.AddPrimative("=", is_number_equal_proc)
	interpreter.AddPrimative("<", is_less_than_proc)
	interpreter.AddPrimative(">", is_greater_than_proc)

	interpreter.AddPrimative("cons", cons_proc)
	interpreter.AddPrimative("car", car_proc)
	interpreter.AddPrimative("cdr", cdr_proc)
	interpreter.AddPrimative("set-car!", set_car_proc)
	interpreter.AddPrimative("set-cdr!", set_cdr_proc)
	interpreter.AddPrimative("list", list_proc)

	interpreter.AddPrimative("eq?", is_eq_proc)

	// TODO: is there a more idiomatic way of doing this in go?
	read_proc := func(object Object) Object { return interpreter.read_proc(object) }

	interpreter.AddPrimative("read", read_proc)

	interpreter.AddPrimative("write-char", write_char_proc)
	interpreter.AddPrimative("write", write_proc)

	interpreter.AddPrimative("error", error_proc)
}
