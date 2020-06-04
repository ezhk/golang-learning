package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const EOF = "_" // append EOF to input string as tail processing

func writeNonZeroRune(sb *strings.Builder, ch rune) {
	// Specially skip nil value
	if ch > 0 {
		sb.WriteRune(ch)
	}
}

func RepeatRune(sb *strings.Builder, ch rune, count int) string {
	for i := 0; i < count; i++ {
		writeNonZeroRune(sb, ch)
	}
	return sb.String()
}

func Unpack(inputLine string) (string, error) {
	var builder strings.Builder
	var factor strings.Builder

	repeatableChar := false
	escapedNext := false
	var previous rune

	extendedInputLine := inputLine + EOF
	for idx, char := range extendedInputLine {
		if idx == 0 && unicode.IsDigit(char) {
			return "", ErrInvalidString
		}
		// Escape next logic
		if escapedNext {
			previous = char
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

			// Possible to use string.Repeat here
			RepeatRune(&builder, previous, intFactor)
			previous = 0
			factor.Reset()
		}
		// Register escaping symbol
		if char == '\\' {
			writeNonZeroRune(&builder, previous)
			escapedNext = true
			continue
		}
		// Register repeatable char
		if previous == char {
			repeatableChar = true
		} else {
			repeatableChar = false
		}

		writeNonZeroRune(&builder, previous)
		previous = char
	}

	return builder.String(), nil
}
