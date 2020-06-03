package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const EOF = "_" // append EOF to input string as tail processing

func Unpack(inputLine string) (string, error) {
	var builder strings.Builder
	var previous string
	repeatableChar := false
	escapedNext := false

	inputLine += EOF

	for idx, char := range inputLine {
		if idx == 0 && unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		// Escape logic
		switch {
		case escapedNext:
			previous = string(char)
			escapedNext = false
			continue
		case string(char) == `\`:
			builder.WriteString(previous)
			escapedNext = true
			continue
		}

		if unicode.IsDigit(char) && !escapedNext {
			counter, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidString
			}

			// Skip repeateble and countered symbols
			if counter > 0 && repeatableChar {
				return "", ErrInvalidString
			}

			builder.WriteString(strings.Repeat(previous, counter))
			previous = ""

			continue
		}

		// Repetable logic
		switch {
		case previous == string(char):
			repeatableChar = true
		case repeatableChar:
			repeatableChar = false
		}

		builder.WriteString(previous)
		previous = string(char)
	}

	return builder.String(), nil
}
