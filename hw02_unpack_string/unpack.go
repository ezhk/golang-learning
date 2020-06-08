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
	var (
		builder strings.Builder
		factor  strings.Builder

		isRepeatableChar bool
		isEscapedNext    bool
		previous         rune
	)

	extendedInputLine := inputLine + EOF
	for idx, char := range extendedInputLine {
		if idx == 0 && unicode.IsDigit(char) {
			return "", ErrInvalidString
		}
		// Escape next logic
		if isEscapedNext {
			previous = char
			isEscapedNext = false
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
			if intFactor > 0 && isRepeatableChar {
				return "", ErrInvalidString
			}

			// Possible to use string.Repeat here, but rune is using
			repeatRune(&builder, previous, intFactor)
			previous = 0
			factor.Reset()
		}
		// Register escaping symbol
		if char == '\\' {
			writeNonZeroRune(&builder, previous)
			isEscapedNext = true
			continue
		}
		// Register repeatable char
		if previous == char {
			isRepeatableChar = true
		} else {
			isRepeatableChar = false
		}

		writeNonZeroRune(&builder, previous)
		previous = char
	}

	return builder.String(), nil
}

func writeNonZeroRune(sb *strings.Builder, ch rune) {
	// Specially skip nil value
	if ch > 0 {
		sb.WriteRune(ch)
	}
}

func repeatRune(sb *strings.Builder, ch rune, count int) *strings.Builder {
	for i := 0; i < count; i++ {
		writeNonZeroRune(sb, ch)
	}

	return sb
}
