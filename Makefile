include $(GOROOT)/src/Make.inc

TARG=hu
GOFILES=\
	objects.go\
	read.go\
	primitives.go\
	bindings.go\
	lex.go\
	recipe.go\
	parser.go\
	dictionary.go\

include $(GOROOT)/src/Make.pkg
