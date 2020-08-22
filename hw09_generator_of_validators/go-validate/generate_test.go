package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("generate document and compare with file", func(t *testing.T) {
		defaultStruct := []ParsedStruct{
			ParsedStruct{
				Name: "User",
				Fields: []Field{
					Field{
						Name:    "ID",
						Type:    "string",
						IsSlice: false,
						Conditions: []Condition{
							Condition{
								Name:  "len",
								Value: "36",
							},
						},
					},
					Field{
						Name:    "Age",
						Type:    "int",
						IsSlice: false,
						Conditions: []Condition{
							Condition{
								Name:  "min",
								Value: "18",
							},
							Condition{
								Name:  "max",
								Value: "50",
							},
						},
					},
					Field{
						Name:    "Email",
						Type:    "string",
						IsSlice: false,
						Conditions: []Condition{
							Condition{
								Name:  "regexp",
								Value: "^\\w+@\\w+\\.\\w+$",
							},
						},
					},
					Field{
						Name:    "Role",
						Type:    "UserRole",
						IsSlice: false,
						Conditions: []Condition{
							Condition{
								Name:  "in",
								Value: "admin,stuff",
							},
						},
					},
					Field{
						Name:    "Phones",
						Type:    "string",
						IsSlice: true,
						Conditions: []Condition{
							Condition{
								Name:  "len",
								Value: "11",
							},
						},
					},
				},
			},
			ParsedStruct{
				Name: "App",
				Fields: []Field{
					Field{
						Name:    "Version",
						Type:    "string",
						IsSlice: false,
						Conditions: []Condition{
							Condition{
								Name:  "len",
								Value: "5",
							},
						},
					},
				},
			},
		}

		document, err := GenerateDocument(defaultStruct)
		require.Nil(t, err)

		fileBytes, err := ioutil.ReadFile("testdata/test_models_validation_generated.go")
		require.Nil(t, err)

		require.Equal(t, fileBytes, document)
	})
}
