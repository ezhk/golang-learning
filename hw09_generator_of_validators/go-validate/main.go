package main

import (
	"fmt"
	"os"
	"strings"
)

type Condition struct {
	Name  string
	Value string
}

type Field struct {
	Name       string
	Type       string
	IsSlice    bool
	Conditions []Condition
}

type ParsedStruct struct {
	Name   string
	Fields []Field
}

func main() {
	filename := os.Args[2]
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
