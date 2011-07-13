include $(GOROOT)/src/Make.inc

TARG=hu_server
GOFILES=\
	minify.go\
	page.go\
	hu_server.go\

include $(GOROOT)/src/Make.cmd
