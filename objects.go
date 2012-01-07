package hu

import (
	"fmt"
	"bytes"
	"big"
)

type Term interface {
	String() string
}

type Reducible interface {
	Reduce(*Environment) Term
}

type Rune int

func (rune Rune) String() string {
	return string(rune)
}

type Boolean bool

func (b Boolean) String() (result string) {
	if b {
		result  = "true"
	} else {
		result  = "false"
	}
	return
}

type Number struct {
	value *big.Int
}

func (n *Number) String() string {
	return n.value.String()
}

type Symbol string

func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) Reduce(environment *Environment) Term {
	return environment.Get(s)
}

type String string

func (s String) String() string {
	var out bytes.Buffer
	for _, rune := range s {
		switch rune {
		case '\n':
			out.WriteString("\\n")
			break
		case '\\':
			out.WriteString("\\\\")
			break
		case '"':
			out.WriteString("\\\"")
			break
		default:
			out.WriteRune(rune)
		}
	}
	return out.String()
}

type Pair struct {
	car, cdr Term
}

func (pair *Pair) String() string {
	return fmt.Sprintf("(%v %v)", pair.car, pair.cdr)
}

type Operator interface {
	apply(*Environment, Term) Term
}

type PrimitiveFunction func(*Environment, Term) Term

func (pf PrimitiveFunction) apply(environment *Environment, term Term) Term {
	return pf(environment, term)
}

func (pf PrimitiveFunction) String() string {
	return fmt.Sprintf("#<primitive-function> %p", pf)
}

type Primitive func(*Environment) Term

func (p Primitive) String() string {
	return fmt.Sprintf("#<primitive> %p", p)
}

func (p Primitive) Reduce(environment *Environment) Term {
	return p(environment)
}

type Application struct {
	term Term
}

func (application Application) String() string {
	return fmt.Sprintf("{%v}", application.term)
}

func (application Application) Reduce(environment *Environment) Term {
	var lhs, last *Pair
	term := application.term.(*Pair)
	next:
	switch operator := environment.evaluate(term.car).(type) {
	case Operator:
		var operands Term
		switch operator.(type) {
		case PrimitiveFunction:
			operands = term.cdr
		default:
			operands = &Pair{lhs, &Pair{term.cdr, nil}}
		}
		return operator.apply(environment, operands)
	default:
		e := &Pair{operator, nil}
		if lhs == nil {	lhs = e	} else { last.cdr = e }
		last = e
	}
	if term.cdr != nil {
		term = term.cdr.(*Pair)
		goto next
	}
	return lhs
}

type Abstraction struct {
	parameters Term
	term Term
}

func (a Abstraction) apply(environment *Environment, values Term) Term {
	e := environment.NewChildEnvironment()
	e.Extend(a.parameters, values)
	return Closure{a.term, e}
}

func (abstraction Abstraction) String() string {
	return fmt.Sprintf("#<abstraction> %v %v", abstraction.parameters, abstraction.term)
}

type Closure struct {
	term Term
	environment *Environment
}

func (closure Closure) String() string {
	return fmt.Sprintf("#<Closure> %v %v\n", closure.term, closure.environment)
}

func (closure Closure) Reduce(environment *Environment) Term {
	return closure.environment.evaluate(closure.term)
}

type Error string

func (error Error) String() string {
	return string(error)
}

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
top:
	switch vars := variables.(type) {
	case *Pair:
		vals := values.(*Pair)
		environment.Extend(vars.car, vals.car)
		variables, values = vars.cdr, vals.cdr
		goto top
	case Symbol:
		environment.Define(vars, environment.parent.Closure(values))
	}
}

func (environment *Environment) Define(variable Symbol, value Term) {
	environment.frame[variable] = value
}

func (environment *Environment) Set(variable Symbol, value Term) {
	_, ok := environment.frame[variable]
	if ok {
		environment.Define(variable, value)
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
	case Reducible:
		term = t.Reduce(environment)
		goto tailcall
	}
	return term
}

func car(term Term) Term {
	return term.(*Pair).car
}

func cdr(term Term) Term {
	return term.(*Pair).cdr
}

func list_from(list Term, selector func(Term) Term) (result Term) {
	switch l := list.(type) {
	case *Pair:
		result = &Pair{selector(l.car), list_from(l.cdr, selector)}
	}
	return
}

func concat(pairs ...*Pair) (result *Pair) {
	var last *Pair
	for _, pair := range pairs {
	next:
		if result == nil {
			result = &Pair{pair.car, nil}
			last = result
		} else {
			p := &Pair{pair.car, nil}
			last.cdr = p
			last = p
		}
		if pair.cdr != nil {
			pair = pair.cdr.(*Pair)
			goto next
		}
	}
	return result
}
