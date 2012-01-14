package main

import (
	"github.com/eikeon/hu"
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


	environment := hu.NewEnvironment()
	hu.AddDefaultBindings(environment)

	for {
		fmt.Printf("hu> ")
		expression := hu.Read(reader)
		if expression != nil {
			result := environment.Evaluate(expression)
			if result != nil {
				fmt.Fprintf(os.Stdout, "%v\n", result)
			}
		} else {
			fmt.Fprintf(os.Stdout, "Goodbye!\n")
			break
		}
	}
}
