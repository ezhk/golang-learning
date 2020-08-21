package main

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
)

const FilterTag string = "validate"

func ParseStruct(filePath string) ([]ParsedStruct, error) {
	parsedStructs := make([]ParsedStruct, 0)
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

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structName := typeSpec.Name.Name

			// read block fields and fetch tags
			fields := make([]Field, 0)
			for _, field := range structType.Fields.List {
				field, err := parseField(field, fset)
				if err != nil {
					continue
				}

				fields = append(fields, field)
			}

			declStruct = ParsedStruct{StructName: structName, Fields: fields}

			// extend slice of parsed structs
			parsedStructs = append(parsedStructs, declStruct)
		}
	}

	return parsedStructs, nil
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

	// skip first and last quotes ` and `
	tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
	if tag.Get(FilterTag) == "" {
		return Field{}, errors.New("empty tag string")
	}

	return Field{Name: fieldName, Type: variableType, Slice: isSlice, Tag: tag.Get(FilterTag)}, nil
}
