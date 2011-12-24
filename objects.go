package hu

import (
	"fmt"
	"bytes"
	"big"
)

type Object interface {
	String() string
}

type RuneObject struct {
	rune int
}

func (o *RuneObject) String() string {
	return string(o.rune)
}

type BooleanObject bool

func (o BooleanObject) String() string {
	var v byte
	if o {
		v = 't'
	} else {
		v = 'f'
	}
	return fmt.Sprintf("#%c", v)
}

type NumberObject struct {
	value *big.Int
}

func (o *NumberObject) String() string {
	return o.value.String()
}

type SymbolObject struct {
	value string
}

func (o *SymbolObject) String() string {
	return o.value
}

type StringObject struct {
	value string
}

func Symbol(value string) Object {
	var element Object

	element = symbol_table
	for element != nil {
		if car(element).(*SymbolObject).value == value {
			return car(element)
		}
		element = cdr(element)
	}

	obj := &SymbolObject{value}
	symbol_table = cons(obj, symbol_table)
	return obj
}

func (o *StringObject) String() string {
	var out bytes.Buffer
	for _, rune := range o.value {
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

type PairObject struct {
	car, cdr Object
}

func (pair *PairObject) String() string {
	var out bytes.Buffer
	out.WriteRune('(')

	car_obj := car(pair)
	cdr_obj := cdr(pair)
	if car_obj == nil {
		fmt.Fprintf(&out, "Ø")
	} else {
		fmt.Fprintf(&out, "%v", car_obj)
	}
	if is_pair(cdr_obj) {
		fmt.Fprintf(&out, " %s", cdr_obj)
	} else if cdr_obj == nil {

	} else {
		if cdr_obj == nil {
			fmt.Fprintf(&out, " . Ø", cdr_obj)
		} else {
			fmt.Fprintf(&out, " . %s", cdr_obj)
		}
	}
	out.WriteRune(')')
	return out.String()
}

type Application struct {
	operator Object
	operands Object
}

func (expression *Application) String() string {
	var out bytes.Buffer
	out.WriteRune('{')

	car_obj := expression.operator
	cdr_obj := expression.operands
	if car_obj == nil {
		fmt.Fprintf(&out, "Ø")
	} else {
		fmt.Fprintf(&out, "%v", car_obj)
	}
	if is_pair(cdr_obj) {
		fmt.Fprintf(&out, " %s", cdr_obj)
	} else if cdr_obj == nil {

	} else {
		if cdr_obj == nil {
			fmt.Fprintf(&out, " . Ø", cdr_obj)
		} else {
			fmt.Fprintf(&out, " . %s", cdr_obj)
		}
	}
	out.WriteRune('}')
	return out.String()
}

type Abstraction struct {
	parameters Object
	object Object
	environment *Environment
}

func (expression *Abstraction) String() string {
	return fmt.Sprintf("#<abstraction> %v %v", expression.parameters, expression.object)
}

type EOFObject struct {

}

func (o *EOFObject) String() string {
	return "#<eof>"
}

type PrimitiveFunction func(*Interpreter, Object, *Environment) Object

type PrimitiveFunctionObject struct {
	function PrimitiveFunction
}

func (o *PrimitiveFunctionObject) String() string {
	return fmt.Sprintf("#<primitive-function> %v", o.function)
}

var (
	TRUE, FALSE             BooleanObject
	eof_object              Object
	quote_symbol            Object
	symbol_table            Object
)

func init() {
	TRUE = true
	FALSE = false
	quote_symbol = Symbol("quote")
	eof_object = &EOFObject{}
}

func is_pair(obj Object) bool {
	_, ok := obj.(*PairObject)
	return ok
}

func is_boolean(obj Object) bool {
	_, ok := obj.(*BooleanObject)
	return ok
}

func is_symbol(obj Object) bool {
	_, ok := obj.(*SymbolObject)
	return ok
}

func is_number(obj Object) bool {
	_, ok := obj.(*NumberObject)
	return ok
}

func is_string(obj Object) bool {
	_, ok := obj.(*StringObject)
	return ok
}

func is_character(obj Object) bool {
	_, ok := obj.(*RuneObject)
	return ok
}

func is_eof_object(obj Object) bool {
	_, ok := obj.(*EOFObject)
	return ok
}

func is_the_empty_list(obj Object) bool {
	return obj == nil
}

func is_false(obj Object) bool {
	return obj == FALSE
}

func is_true(obj Object) bool {
	return is_false(obj) == false
}

func cons(car, cdr Object) Object {
	return &PairObject{car, cdr}
}

func car(object Object) Object {
	return object.(*PairObject).car
}

func cdr(object Object) Object {
	return object.(*PairObject).cdr
}

func set_car(object, value Object) {
	object.(*PairObject).car = value
}

func set_cdr(object, value Object) {
	object.(*PairObject).cdr = value
}

func is_last(seq Object) bool {
	return is_the_empty_list(cdr(seq))
}

func list_from(list Object, selector func(Object) Object) (result Object) {
	if list != nil {
		result = &PairObject{selector(car(list)), list_from(cdr(list), selector)}
	}
	return
}
