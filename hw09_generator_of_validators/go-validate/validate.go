package main

import (
	"errors"
)

var AllowedValidators map[string]map[string]struct{} = map[string]map[string]struct{}{
	"int": {
		"min": struct{}{},
		"max": struct{}{},
		"in":  struct{}{},
	},
	"string": {
		"len":    struct{}{},
		"regexp": struct{}{},
		"in":     struct{}{},
	},
}

func ValidateStruct(parsedStructs []ParsedStruct) []ParsedStruct {
	// func read slices of structs and validate them;
	// validated data redefines earlier data
	validatedStructs := make([]ParsedStruct, 0)
	for _, parsedSt := range parsedStructs {
		vStruct := ParsedStruct{Name: parsedSt.Name}
		for _, field := range parsedSt.Fields {
			validatedField, err := ValidateField(field)
			if err != nil {
				continue
			}
			vStruct.Fields = append(vStruct.Fields, validatedField)
		}
		validatedStructs = append(validatedStructs, vStruct)
	}

	return validatedStructs
}

func ValidateField(field Field) (Field, error) {
	validators, ok := AllowedValidators[field.Type]
	if !ok {
		return Field{}, errors.New("expression not valid")
	}

	for _, condition := range field.Conditions {
		if _, ok := validators[condition.Name]; !ok {
			return Field{}, errors.New("field name not valid")
		}
	}

	return field, nil
}
