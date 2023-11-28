package main

import (
	"fmt"
	"os"

	griffonParser "github.com/bensooraj/griffon/parser"
)

func main() {
	fmt.Println("Griffon!")
	// read from a file
	b, err := os.ReadFile("testdata/test1.hcl")
	if err != nil {
		panic(err)
	}
	// parse the file
	config, err := griffonParser.ParseWithBodySchema("testdata/test1.hcl", b, griffonParser.GetEvalContext(), nil)
	if err != nil {
		panic(err)
	}
	_ = config
	// spew.Dump(config)
}
