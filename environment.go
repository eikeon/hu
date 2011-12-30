package hu

import "fmt"

type UnboundVariableError struct {
	variable Term
	operation string
}

func (e UnboundVariableError) String() string {
	return "Unbound Variable: " + e.variable.String() + " operation: " + e.operation
}

type Environment struct {
	frame  map[Symbol]Term
	parent *Environment
}

func (environment *Environment) String() string {
	return "#<environment>"
}

func NewEnvironment() *Environment {
	return &Environment{frame: make(map[Symbol]Term)}
}

func ClosureIfNeeded(term Term, environment *Environment) Term {
	switch v := term.(type) {
	case Application:
		return Closure{term, environment}
	}
	return term
}

// returns a new (child) environment from this environment extended
// with bindings given by variables, values.
func (environment *Environment) NewChildEnvironment(variables, values Term) *Environment {
	child := NewEnvironment()
	child.parent = environment
	for ; variables != nil && values != nil; variables, values = cdr(variables), cdr(values) {
		switch variables.(type) {
		case *Pair:
			variable, value := car(variables), car(values)
			child.frame[variable.(Symbol)] = ClosureIfNeeded(value, environment)
		default:
			panic("TODO: needs a test case")
			child.frame[variables.(Symbol)] = ClosureIfNeeded(values, environment)
			return child
		}
	}
	return child
}

func (environment *Environment) Define(variable Symbol, value Term) {
	environment.frame[variable] = ClosureIfNeeded(value, environment)
}

func (environment *Environment) Set(variable Symbol, value Term) {
	_, ok := environment.frame[variable]
	if ok {
		environment.frame[variable] = ClosureIfNeeded(value, environment)
	} else if environment.parent != nil {
		environment.parent.Set(variable, value)
	} else {
		panic(UnboundVariableError{variable, "set"})
	}
}

func (environment *Environment) Get(variable Symbol) Term {
	value, ok := environment.frame[variable]
	if ok {
		return value
	} else if environment.parent != nil {
		return environment.parent.Get(variable)
	} else {
		panic(UnboundVariableError{variable, "get"})
	}
	return nil
}

func (environment *Environment) AddPrimitive(name string, function PrimitiveFunction) {
	environment.Define(Symbol(name), function)
}

func (environment *Environment) Evaluate(term Term) (result Term) {
	defer func() {
		switch x := recover().(type) {
		case Term:
			result = x
		case interface{}:
      			result = Error(fmt.Sprintf("%v", x))
		}
	}()
	result = environment.evaluate(term)
	return
}

func (environment *Environment) evaluate(term Term) Term {
tailcall:
	switch o := term.(type) {
	case Symbol:
		term = environment.Get(o)
		goto tailcall
	case Closure:
		term = o.environment.evaluate(o.term)
		goto tailcall
	case Application:
	 	switch operator:= environment.evaluate(o.operator).(type) {
	 	case PrimitiveFunction:
	 		term = operator(environment, o.operands)
	 		goto tailcall
		case Abstraction:
			environment = environment.NewChildEnvironment(operator.parameters, o.operands)
			term = environment.evaluate(operator.term)
			goto tailcall
		default:
			panic(fmt.Sprintf("operator %v of unknown type %T", operator, operator))
	 	}
	}
	return term
}
