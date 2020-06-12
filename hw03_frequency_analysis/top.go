package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

var SplitFilter = regexp.MustCompile(`[\s?!.;,]`)
var IgnoredSymbols = regexp.MustCompile(`[-]`)

type WordFrequency struct {
	word          string
	RepeatCounter int
}

func Top10(inputLine string) []string {
	const MaxLength = 10
	mostFrequentWords := make([]string, 0, MaxLength)

	// prepare map with counter
	freqMap := make(map[string]int)
	for _, word := range SplitFilter.Split(inputLine, -1) {
		word = IgnoredSymbols.ReplaceAllString(word, "")
		if len(word) > 0 {
			freqMap[strings.ToLower(word)]++
		}
	}

	// convert map to slice wordFrequency for next order
	wordSlice := make([]WordFrequency, 0, len(freqMap))
	for word, counter := range freqMap {
		wordSlice = append(wordSlice, WordFrequency{word, counter})
	}
	sort.Slice(wordSlice, func(i, j int) bool {
		return wordSlice[i].RepeatCounter > wordSlice[j].RepeatCounter
	})

	// store data to result slice
	for _, value := range wordSlice {
		if len(mostFrequentWords) >= MaxLength {
			break
		}
		mostFrequentWords = append(mostFrequentWords, value.word)
	}

	return mostFrequentWords
}
