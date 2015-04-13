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
	v, _ := environment.Get(s)
	return v
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
		switch operator := Evaluate(environment, term).(type) {
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

func (a Abstraction) apply(e Environment, values Term) Term {
	c := &NestedEnvironment{Environment: make(LocalEnvironment), Parent: e}
	Extend(c, a.Parameters, values)
	return Closure{a.Term, c}
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
	return Evaluate(closure.Environment, closure.Term)
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

func Extend(environment Environment, variables, values Term) {
	switch vars := variables.(type) {
	case Tuple:
		vals := values.(Tuple)
		if len(vals) != len(vars) {
			fmt.Println("type mismatch:", vals, vars)
		}
		for i, v := range vars {
			val := vals[i]
			Extend(environment, v, val)
		}
	case Symbol:
		if vars != Term(nil) {
			parent := environment.(*NestedEnvironment).Parent
			environment.Define(vars, Closure{values, parent})
		}
	}
}

type Environment interface {
	Define(variable Symbol, value Term)
	Set(variable Symbol, value Term) bool
	Get(variable Symbol) (Term, bool)
}

type LocalEnvironment map[Symbol]Term

func (environment LocalEnvironment) Define(variable Symbol, value Term) {
	environment[variable] = value
}

func (environment LocalEnvironment) Set(variable Symbol, value Term) bool {
	environment[variable] = value
	return true
}

func (environment LocalEnvironment) Get(variable Symbol) (Term, bool) {
	value, ok := environment[variable]
	if ok {
		return value, ok
	} else {
		return UnboundVariableError{variable, "get"}, false
	}

}

type NestedEnvironment struct {
	Environment Environment
	Parent      Environment
}

func (environment *NestedEnvironment) String() string {
	return "#<environment>"
}

func (ne *NestedEnvironment) Define(variable Symbol, value Term) {
	ne.Environment.Define(variable, value)
}

func (ne *NestedEnvironment) Set(variable Symbol, value Term) bool {
	_, ok := ne.Environment.Get(variable)
	if ok {
		ne.Define(variable, value)
		return true
	} else if ne.Parent != nil {
		return ne.Parent.Set(variable, value)
	} else {
		return false
	}
}

func (ne *NestedEnvironment) Get(variable Symbol) (Term, bool) {
	value, ok := ne.Environment.Get(variable)
	if ok {
		return value, true
	} else if ne.Parent != nil {
		return ne.Parent.Get(variable)
	} else {
		return UnboundVariableError{variable, "get"}, false
	}
}

type Property struct {
	Name   Symbol
	DidSet Abstraction
}

func (property Property) String() string {
	return fmt.Sprintf("#<property> %v", property.Name)
}

func Evaluate(environment Environment, term Term) Term {
tailcall:
	switch t := term.(type) {
	case Reducible:
		term = t.Reduce(environment)
		goto tailcall
	}
	return term
}

func GuardedEvaluate(environment Environment, expression Term) (result Term) {
	defer func() {
		switch x := recover().(type) {
		case Term:
			result = x
		case interface{}:
			result = Error(fmt.Sprintf("%v", x))
		}
	}()
	result = Evaluate(environment, expression)
	return
}
