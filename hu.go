// Package hu implements (an interpreter for) a language optimized for
// humans.
package hu

import (
	"fmt"
	"math/big"
	"strings"
)

type Term interface {
	String() string
}

type Reducible interface {
	Term
	Reduce(Environment) Term
}

type Rune int

func (rune Rune) String() string {
	return string(rune)
}

type Boolean bool

func (b Boolean) String() (result string) {
	if b {
		result = "true"
	} else {
		result = "false"
	}
	return
}

type Number struct {
	value *big.Rat
}

func (n *Number) String() string {
	return n.value.RatString()
}

type Symbol string

func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) Reduce(environment Environment) Term {
	return environment.Get(s)
}

type String string

func (s String) String() string {
	return string(s)
}

type Tuple []Term

func (tuple Tuple) String() string {
	return fmt.Sprintf("(%v)", []Term(tuple))
}

type Set []Term

func (set Set) String() string {
	return fmt.Sprintf("{%v}", []Term(set))
}

type Part []Term

func (part Part) String() string {
	var terms []string
	for _, term := range part {
		terms = append(terms, term.String())
	}
	return strings.Join(terms, "")
}

type Operator interface {
	Term
	apply(Environment, Term) Term
}

type PrimitiveFunction func(Environment, Term) Term

func (pf PrimitiveFunction) apply(environment Environment, term Term) Term {
	return pf(environment, term)
}

func (pf PrimitiveFunction) String() string {
	return fmt.Sprintf("#<primitive-function> %p", pf)
}

type Primitive func(Environment) Term

func (p Primitive) String() string {
	return fmt.Sprintf("#<primitive> %p", p)
}

func (p Primitive) Reduce(environment Environment) Term {
	return p(environment)
}

type Application []Term

func (application Application) String() string {
	return fmt.Sprintf("{%v}", []Term(application))
}

func (application Application) Reduce(environment Environment) Term {
	for i, term := range application {
		switch operator := environment.Evaluate(term).(type) {
		case Operator:
			var operands Term
			switch operator.(type) {
			case PrimitiveFunction:
				operands = Tuple(application[i+1:])
			default:
				lhs := Tuple(application[0:i])
				rhs := Tuple(application[i+1:])
				operands = Tuple([]Term{lhs, rhs})
			}
			return operator.apply(environment, operands)
		}
	}
	return nil
}

type Abstraction struct {
	Parameters Term
	Term       Term
}

func (a Abstraction) apply(environment Environment, values Term) Term {
	e := environment.NewChildEnvironment()
	e.Extend(a.Parameters, values)
	return Closure{a.Term, e}
}

func (abstraction Abstraction) String() string {
	return fmt.Sprintf("#<abstraction> %v %v", abstraction.Parameters, abstraction.Term)
}

type Closure struct {
	Term        Term
	Environment Environment
}

func (closure Closure) String() string {
	return fmt.Sprintf("#<Closure> %v %v\n", closure.Term, closure.Environment)
}

func (closure Closure) Reduce(environment Environment) Term {
	return closure.Environment.Evaluate(closure.Term)
}

type Error string

func (error Error) String() string {
	return string(error)
}

type UnboundVariableError struct {
	variable  Term
	operation string
}

func (e UnboundVariableError) String() string {
	return "Unbound Variable: " + e.variable.String() + " operation: " + e.operation
}

type environment struct {
	frame  Frame
	parent Environment
}

func (environment *environment) String() string {
	return "#<environment>"
}

func NewEnvironment() Environment {
	return &environment{frame: make(LocalFrame)}
}

func NewEnvironmentWithFrame(f Frame) *environment {
	return &environment{frame: f}
}

func NewEnvironmentWithParent(p Environment) *environment {
	return &environment{frame: make(LocalFrame), parent: p}
}

// returns a new (child) environment from this environment extended
// with bindings given by variables, values.
func (e *environment) NewChildEnvironment() Environment {
	child := &environment{frame: make(LocalFrame)}
	child.parent = e
	return child
}

func (environment *environment) Extend(variables, values Term) {
	switch vars := variables.(type) {
	case Tuple:
		vals := values.(Tuple)
		if len(vals) != len(vars) {
			fmt.Println("type mismatch:", vals, vars)
		}
		for i, v := range vars {
			val := vals[i]
			environment.Extend(v, val)
		}
	case Symbol:
		if vars != Term(nil) {
			environment.Define(vars, Closure{values, environment.parent})
		}
	}
}

type Frame interface {
	Define(variable Symbol, value Term)
	Set(variable Symbol, value Term) bool
	Get(variable Symbol) (Term, bool)
}

type Environment interface {
	Define(variable Symbol, value Term)
	Set(variable Symbol, value Term)
	Get(variable Symbol) Term
	Evaluate(term Term) (result Term)
	NewChildEnvironment() Environment
	Extend(variables, values Term)
}

type LocalFrame map[Symbol]Term

func (frame LocalFrame) Define(variable Symbol, value Term) {
	frame[variable] = value
}

func (frame LocalFrame) Set(variable Symbol, value Term) bool {
	_, ok := frame[variable]
	return ok
}

func (frame LocalFrame) Get(variable Symbol) (Term, bool) {
	value, ok := frame[variable]
	return value, ok
}

func (environment *environment) Define(variable Symbol, value Term) {
	environment.frame.Define(variable, value)
}

func (environment *environment) Set(variable Symbol, value Term) {
	_, ok := environment.frame.Get(variable)
	if ok {
		environment.Define(variable, value)
	} else if environment.parent != nil {
		environment.parent.Set(variable, value)
	} else {
		panic(UnboundVariableError{variable, "set"})
	}
}

func (environment *environment) Get(variable Symbol) Term {
	value, ok := environment.frame.Get(variable)
	if ok {
		return value
	} else if environment.parent != nil {
		return environment.parent.Get(variable)
	} else {
		panic(UnboundVariableError{variable, "get"})
	}
	return nil
}

func (environment *environment) AddPrimitive(name string, function PrimitiveFunction) {
	environment.Define(Symbol(name), function)
}

func (environment *environment) _Evaluate(term Term) (result Term) {
	defer func() {
		switch x := recover().(type) {
		case Term:
			result = x
		case interface{}:
			result = Error(fmt.Sprintf("%v", x))
		}
	}()
	result = environment.Evaluate(term)
	return
}

func (environment *environment) Evaluate(term Term) Term {
tailcall:
	switch t := term.(type) {
	case Reducible:
		term = t.Reduce(environment)
		goto tailcall
	}
	return term
}
