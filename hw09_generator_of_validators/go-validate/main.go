package main

import (
	"fmt"
	"os"
)

// type Term struct {
// 	Name     string
// 	Argument string
// }

type Field struct {
	Name  string
	Type  string
	Slice bool
	// Terms []Term
	Tag string
}

type ParsedStruct struct {
	StructName string
	Fields     []Field
}

func main() {
	filename := os.Args[2]
	parsedStruct, err := ParseStruct(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	validatedStruct := ValidateStruct(parsedStruct)
	fmt.Printf("%+v\n", validatedStruct)

	// for _, readStruct := range parsedStruct {
	// 	fmt.Printf("%#v\n", readStruct)
	// }
}
