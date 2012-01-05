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

// returns a new (child) environment from this environment extended
// with bindings given by variables, values.
func (environment *Environment) NewChildEnvironment() *Environment {
	child := NewEnvironment()
	child.parent = environment
	return child
}

func (environment *Environment) Closure(term Term) Term {
	switch v := term.(type) {
	case Application:
		return Closure{term, environment}
	}
	return term
}

func (environment *Environment) Extend(variables, values Term) {
	for ; variables != nil && values != nil; variables, values = cdr(variables), cdr(values) {
		switch variables.(type) {
		case *Pair:
			environment.Extend(car(variables), car(values))
		default:
			environment.frame[variables.(Symbol)] = environment.parent.Closure(values)
			return
		}
	}
}

func (environment *Environment) Define(variable Symbol, value Term) {
	environment.frame[variable] = environment.Closure(value)
}

func (environment *Environment) Set(variable Symbol, value Term) {
	_, ok := environment.frame[variable]
	if ok {
		environment.frame[variable] = environment.Closure(value)
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
	switch t := term.(type) {
	case Evaluable:
		term = t.Evaluate(environment)
		goto tailcall
	}
	return term
}
