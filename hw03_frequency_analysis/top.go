package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var ignoredSymbols = regexp.MustCompile(`[-]`)

type WordFrequency struct {
	word          string
	RepeatCounter int
}

func Top10(inputLine string) []string {
	const MaxLength = 10
	mostFrequentWords := make([]string, 0, MaxLength)

	// prepare map with counter
	freqMap := make(map[string]int)
	inputLine = ignoredSymbols.ReplaceAllString(inputLine, "")
	for _, word := range strings.FieldsFunc(inputLine, splitFunc) {
		if len(word) > 0 {
			freqMap[strings.ToLower(word)]++
		}
	}

	// convert map to slice wordFrequency for next sorting
	wordSlice := make([]WordFrequency, 0, len(freqMap))
	for word, counter := range freqMap {
		wordSlice = append(wordSlice, WordFrequency{word, counter})
	}
	sort.Slice(wordSlice, func(i, j int) bool {
		return wordSlice[i].RepeatCounter > wordSlice[j].RepeatCounter
	})

	// store data
	for _, value := range wordSlice {
		if len(mostFrequentWords) >= MaxLength {
			break
		}
		mostFrequentWords = append(mostFrequentWords, value.word)
	}

	return mostFrequentWords
}

func splitFunc(char rune) bool {
	return unicode.IsPunct(char) || unicode.IsSpace(char)
}
