package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("zero valid fields", func(t *testing.T) {
		defaultStruct := []ParsedStruct{
			ParsedStruct{
				Name: "App",
				Fields: []Field{
					Field{
						Name:    "InvalidFilter",
						Type:    "string",
						IsSlice: false,
						Conditions: []Condition{
							Condition{
								Name:  "unknownfilter",
								Value: "0",
							},
						},
					},
				},
			},
		}

		validatedStruct := ValidateStruct(defaultStruct)
		require.Equal(t, 0, len(validatedStruct[0].Fields))
	})
}
