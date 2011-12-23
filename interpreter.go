package hu

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{NewEnvironment()}
}

func (interpreter *Interpreter) AddPrimitive(name string, function PrimitiveFunction) {
	interpreter.environment.Define(Symbol(name), &PrimitiveFunctionObject{function})
}

func (interpreter *Interpreter) AddPrimitiveProcedure(name string, function PrimitiveFunction) {
	procedure := func(interpreter *Interpreter, object Object, environment *Environment) Object {
		return function(interpreter, interpreter.evaluate(object, environment), environment)
	}
	interpreter.environment.Define(Symbol(name), &PrimitiveFunctionObject{procedure})
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
		return cons(interpreter.evaluate(o.car, environment), interpreter.evaluate(o.cdr, environment))
	case *Application:
		switch operator := interpreter.evaluate(o.operator, environment).(type) {
		case *PrimitiveFunctionObject:
			object = operator.function(interpreter, o.operands, environment)
			goto tailcall
		case *Abstraction:
			operands := interpreter.evaluate(o.operands, environment)
			environment = operator.environment.Extend(operator.parameters, operands)
			object = operator.object
			goto tailcall
		}
	}
	return object
}
