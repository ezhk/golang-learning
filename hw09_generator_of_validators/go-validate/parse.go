package main

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

const FilterTag string = "validate"

func ParseStruct(filePath string) ([]ParsedStruct, error) {
	parsedStructs := make([]ParsedStruct, 0)
	aliasedTypes := make(AliasedStructs)

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// range decalaration block, like a "type ()"
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// read each specs and skip block without structs and names
		for _, spec := range genDecl.Specs {
			var declStruct ParsedStruct

			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			var (
				structName string
				structType *ast.StructType
			)
			switch v := typeSpec.Type.(type) {
			case *ast.Ident:
				structName = v.Name
				aliasedTypes[typeSpec.Name.Name] = structName

				continue
			case *ast.StructType:
				structType = v
				structName = typeSpec.Name.Name
			default:
				continue
			}

			// read block fields and fetch tags
			fields := make([]Field, 0)
			for _, field := range structType.Fields.List {
				field, err := parseField(field, fset)
				if err != nil {
					continue
				}

				fields = append(fields, field)
			}
			declStruct = ParsedStruct{Name: structName, Fields: fields}
			// extend slice of parsed structs
			parsedStructs = append(parsedStructs, declStruct)
		}
	}

	return ApplyAliases(ParsedDocument{Structs: parsedStructs, Aliases: aliasedTypes}), nil
}

func parseField(field *ast.Field, fset *token.FileSet) (Field, error) {
	// tag, type and names must be defined to Tag struct
	if field.Tag == nil || field.Type == nil || field.Names == nil || len(field.Names) < 1 {
		return Field{}, errors.New("invalid field structure")
	}

	// convert field type to string
	var fieldType bytes.Buffer
	err := format.Node(&fieldType, fset, field.Type)
	if err != nil {
		return Field{}, err
	}

	// Check that variable is slice or not
	var (
		isSlice      bool
		variableType string
	)
	if readType := fieldType.String(); readType[:2] == "[]" {
		isSlice = true
		variableType = readType[2:]
	} else {
		variableType = fieldType.String()
	}

	fieldName := field.Names[0].Name

	conditions, err := prepareConditions(field.Tag.Value)
	if err != nil {
		return Field{}, err
	}

	// skip first and last quotes ` and `
	tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
	if tag.Get(FilterTag) == "" {
		return Field{}, errors.New("empty tag string")
	}

	return Field{Name: fieldName, Type: variableType, IsSlice: isSlice, Conditions: conditions}, nil
}

func prepareConditions(fieldTag string) ([]Condition, error) {
	conditions := make([]Condition, 0)

	tag := reflect.StructTag(fieldTag[1 : len(fieldTag)-1])
	if tag.Get(FilterTag) == "" {
		return conditions, errors.New("incorrect filter tag")
	}

	rules := strings.Split(tag.Get(FilterTag), "|")
	for _, r := range rules {
		conditionSlice := strings.SplitN(r, ":", 2)
		cond := Condition{Name: conditionSlice[0], Value: conditionSlice[1]}
		conditions = append(conditions, cond)
	}

	if len(conditions) < 1 {
		return conditions, errors.New("empty filter tags")
	}

	return conditions, nil
}

func ApplyAliases(parsedDocument ParsedDocument) []ParsedStruct {
	// function parse ParsedDocument and convert aliased types into based

	validatedStructs := make([]ParsedStruct, 0)
	for _, parsedSt := range parsedDocument.Structs {
		vStruct := ParsedStruct{Name: parsedSt.Name}
		for _, field := range parsedSt.Fields {
			if value, ok := parsedDocument.Aliases[field.Type]; ok {
				field.Type = value
			}
			vStruct.Fields = append(vStruct.Fields, field)
		}
		validatedStructs = append(validatedStructs, vStruct)
	}

	return validatedStructs
}
