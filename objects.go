package hu

import (
	"fmt"
	"bytes"
	"big"
)

type Term interface {
	String() string
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

type Application struct {
	term Term
}

func (application Application) String() string {
	return fmt.Sprintf("{%v}", application.term)
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

type Error string

func (error Error) String() string {
	return string(error)
}

func cons(car, cdr Term) Term {
	return &Pair{car, cdr}
}

func car(term Term) Term {
	return term.(*Pair).car
}

func cdr(term Term) Term {
	return term.(*Pair).cdr
}

func list_from(list Term, selector func(Term) Term) (result Term) {
	if list != nil {
		result = &Pair{selector(car(list)), list_from(cdr(list), selector)}
	}
	return
}

func concat(pairs ...*Pair) (result *Pair) {
	var last *Pair
	for _, pair := range pairs {
		var term Term = pair
		for ; term != nil; term = cdr(term) {
			if result == nil {
				result = &Pair{car(term), nil}
				last = result
			} else {
				p := &Pair{car(term), nil}
				last.cdr = p
				last = p
			}
		}
	}
	return result
}
