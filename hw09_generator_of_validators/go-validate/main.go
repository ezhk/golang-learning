package main

import (
	"fmt"
	"os"
	"strings"
)

type ParsedDocument struct {
	Structs []ParsedStruct
	Aliases AliasedStructs
}

type ParsedStruct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name       string
	Type       string
	IsSlice    bool
	Conditions []Condition
}

type Condition struct {
	Name  string
	Value string
}

type AliasedStructs map[string]string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-validate [module path]")
		os.Exit(1)
	}

	filename := os.Args[1]
	parsedStruct, err := ParseStruct(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	validatedStruct := ValidateStruct(parsedStruct)
	document, err := GenerateDocument(validatedStruct)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	savedFilename := strings.TrimSuffix(filename, ".go")
	f, err := os.Create(fmt.Sprintf("%s_validation_generated.go", savedFilename))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.Write(document)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
