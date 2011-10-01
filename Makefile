include $(GOROOT)/src/Make.inc

TARG=hu
GOFILES=\
	lex.go\
	recipe.go\
	recipes.go\
	parser.go\
	dictionary.go\

include $(GOROOT)/src/Make.pkg
