package hu

type Interpreter struct {
	environment                                          *Environment
	quote_symbol, lambda_symbol, begin_symbol Object
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{}
	interpreter.environment = NewEnvironment()

	interpreter.quote_symbol = Symbol("quote")
	interpreter.lambda_symbol = Symbol("lambda")
	interpreter.begin_symbol = Symbol("begin")

	return interpreter
}

func (interpreter *Interpreter) AddPrimative(name string, function PrimativeProcedure) {
	interpreter.environment.Define(Symbol(name), &PrimativeProcedureObject{function})
}

func (interpreter *Interpreter) AddMacro(name string, macro Macro) {
	interpreter.environment.Define(Symbol(name), &MacroObject{macro})
}

func (interpreter *Interpreter) Evaluate(object Object) Object {
	return interpreter.evaluate(object, interpreter.environment)
}

func (interpreter *Interpreter) evaluate(object Object, environment *Environment) Object {
tailcall:
	switch o := object.(type) {
	case *SymbolObject:
		return environment.Get(o)
	case *PairObject:
		switch operator := interpreter.evaluate(o.car, environment).(type) {
		case *PrimativeProcedureObject:
			operands := o.cdr
			eval := func(operand Object) Object { return interpreter.evaluate(operand, environment) }
			arguments := list_from(operands, eval)
			object = operator.function(arguments)
			return object
		case *MacroObject:
			operands := o.cdr
			object = operator.expand(interpreter, operands, environment)
			return object
		case *CompoundProcedureObject:
			operands := o.cdr
			eval := func(operand Object) Object { return interpreter.evaluate(operand, environment) }
			arguments := list_from(operands, eval)
			environment = operator.environment.Extend(operator.parameters, arguments)
			object = &PairObject{interpreter.begin_symbol, operator.body}
			goto tailcall
		}
	}
	return object
}
