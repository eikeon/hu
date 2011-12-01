include $(GOROOT)/src/Make.inc

TARG=hu
GOFILES=\
	objects.go\
	read.go\
	environment.go\
	interpreter.go\
	primitives.go\
	macros.go\
	bindings.go\
	lex.go\
	recipe.go\
	recipes.go\
	parser.go\
	resize.go\
	dictionary.go\

include $(GOROOT)/src/Make.pkg
