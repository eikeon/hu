package hu

import (
	"testing"
	"strings"
)

func TestInterpreter(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddPrimative("+", add_proc)

	reader := strings.NewReader("(+ 1 2)")
	expression := interpreter.Read(reader)
	result := interpreter.Evaluate(expression)
	if result.(*NumberObject).value != 3 {
		t.Error("unexpected result")
	}
}
