package main

import (
	"hu"
	"fmt"
	"os"
	"bufio"
	"flag"
	"io"
	"log"
)

func main() {
	filenameFlag := flag.String("filename", "-", "filename from which to read and execute program")
	flag.Parse()
	filename := *filenameFlag

	var reader io.RuneScanner
	if filename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalln(err)
		}
		reader = bufio.NewReader(f)
	}


	interpreter := hu.NewInterpreter()
	interpreter.AddDefaultBindings()

	for {
		fmt.Printf("hu> ")
		result := interpreter.Read(reader)
		if result != nil {
			fmt.Fprintf(os.Stdout, "%v\n", result)
		}
	}
}
