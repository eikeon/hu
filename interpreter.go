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
		return function(interpreter, interpreter.evalList(object, environment), environment)
	}
	interpreter.environment.Define(Symbol(name), &PrimitiveFunctionObject{procedure})
}

// TODO: refactor to take function and parameters as one argument of type Object
func (interpreter *Interpreter) AddClosure(function, parameters Object, outer *Environment) Object {
	f := func(interpreter *Interpreter, object Object, environment *Environment) Object {
		operands := interpreter.evalList(object, environment)
		environment = outer.Extend(parameters, operands)
		return interpreter.begin(function, environment)
	}
	return &PrimitiveFunctionObject{f}
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
		operands := o.cdr
		switch operator := interpreter.evaluate(o.car, environment).(type) {
		case *PrimitiveFunctionObject:
			object = operator.function(interpreter, operands, environment)
			goto tailcall
		}
	}
	return object
}
