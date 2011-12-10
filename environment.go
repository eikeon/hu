package hu

import (
	"fmt"
)

type Environment struct {
	frame  map[Object]Object
	parent *Environment
}

func (environment *Environment) String() string {
	return "#<environment>"
}

func make_frame(variables, values Object) map[Object]Object {
	frame := make(map[Object]Object)
	for ; variables != nil && values != nil; variables, values = cdr(variables), cdr(values) {
		if is_pair(variables) {
			variable, value := car(variables), car(values)
			frame[variable] = value
		} else {
			// TODO: needs a test case
			frame[variables] = values
			break
		}
	}
	return frame
}

func NewEnvironment() *Environment {
	return &Environment{frame: make(map[Object]Object)}
}

// returns a new (child) environment from this environment extended
// with bindings given by variables, values.
func (environment *Environment) Extend(variables, values Object) *Environment {
	return &Environment{make_frame(variables, values), environment}

}

func (environment *Environment) Define(variable, value Object) {
	environment.frame[variable] = value
}

func (environment *Environment) Set(variable, value Object) {
	_, ok := environment.frame[variable]
	if ok {
		environment.frame[variable] = value
	} else if environment.parent != nil {
		environment.parent.Set(variable, value)
	} else {
		panic(fmt.Sprintf("unbound variable '%s'\n", variable))
	}
}

func (environment *Environment) Get(variable Object) Object {
	value, ok := environment.frame[variable]
	if ok {
		return value
	} else if environment.parent != nil {
		return environment.parent.Get(variable)
	} else {
		fmt.Printf("unbound variable '%s'\n", variable)
	}
	return nil
}
