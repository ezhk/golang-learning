package main

import (
	"errors"
	"fmt"
	"strings"
)

var AllowedValidators map[string][]string = map[string][]string{
	"int":    {"min", "max", "in"},
	"string": {"len", "regexp", "in"},
}

func ValidateStruct(parsedStructs []ParsedStruct) []ParsedStruct {
	// func read slices of structs and validate them;
	// validated data redefines earlier data
	validatedStructs := make([]ParsedStruct, 0)
	for _, parsedSt := range parsedStructs {
		vStruct := ParsedStruct{StructName: parsedSt.StructName}
		for _, field := range parsedSt.Fields {
			validatedField, err := ValidateField(field)
			if err != nil {
				continue
			}

			fmt.Println(GeneratorField(field))
			vStruct.Fields = append(vStruct.Fields, validatedField)
		}
		validatedStructs = append(validatedStructs, vStruct)
	}

	return validatedStructs
}

func ValidateField(field Field) (Field, error) {
	validateCondition := strings.SplitN(field.Tag, ":", 2)
	validateKeyword := validateCondition[0]

	validators, ok := AllowedValidators[field.Type]
	if !ok {
		return Field{}, errors.New("empression not valid")
	}

	for _, v := range validators {
		if v == validateKeyword {
			return field, nil
		}
	}

	return field, errors.New("not found valid validators")
}
