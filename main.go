package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Griffon!")
	// read from a file
	b, err := os.ReadFile("testdata/test1.hcl")
	if err != nil {
		panic(err)
	}
	// parse the file
	config, err := ParseHCLUsingBodySchema("testdata/test1.hcl", b, getEvalContext())
	if err != nil {
		panic(err)
	}
	_ = config
	// spew.Dump(config)
}
