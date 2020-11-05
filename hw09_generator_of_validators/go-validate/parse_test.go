package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("parse test file", func(t *testing.T) {
		waitingStruct := []ParsedStruct{
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
		parsedStruct, err := ParseStruct("testdata/test_models.go")

		require.Nil(t, err)
		require.Equal(t, waitingStruct, parsedStruct)
	})
}
