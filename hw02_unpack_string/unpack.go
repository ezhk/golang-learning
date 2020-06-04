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
	var factor strings.Builder

	repeatableChar := false
	escapedNext := false
	previous := ""

	extendedInputLine := inputLine + EOF
	for idx, char := range extendedInputLine {
		if idx == 0 && unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		// Escape next logic
		if escapedNext {
			previous = string(char)
			escapedNext = false
			continue
		}

		// Multiplicator might be two and more symbols
		if unicode.IsDigit(char) {
			factor.WriteRune(char)
			continue
		}

		// Meet char and defined repeateble factor
		if factor.Len() > 0 {
			intFactor, err := strconv.Atoi(factor.String())
			if err != nil {
				return "", ErrInvalidString
			}

			// Skip repeateble and countered symbols
			if intFactor > 0 && repeatableChar {
				return "", ErrInvalidString
			}
			builder.WriteString(strings.Repeat(previous, intFactor))

			previous = ""
			factor.Reset()
		}

		// Register escaping symbol
		if string(char) == `\` {
			builder.WriteString(previous)
			escapedNext = true
			continue
		}

		// Register repeatable char
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
