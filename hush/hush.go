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

	var result hu.Term
	fmt.Printf("hu> ")
	for {
		expression := hu.Read(reader)
		if expression != nil {
			if expression == hu.Symbol("\n") {
				if result != nil {
					fmt.Fprintf(os.Stdout, "%v\n", result)
				}
				fmt.Printf("hu> ")
				continue
			} else {
				result = environment.Evaluate(expression)
			}
		} else {
			fmt.Fprintf(os.Stdout, "Goodbye!\n")
			break
		}
	}
}
