package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type wordCount struct {
	word  string
	count int
}

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	freq := make(map[string]int, len(words))

	for _, word := range words {
		word = strings.ToLower(word)
		word = strings.TrimRightFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		if word == "" {
			continue
		}
		freq[word]++
	}

	pairs := make([]wordCount, 0, len(freq))
	for word, count := range freq {
		pairs = append(pairs, wordCount{word: word, count: count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].word < pairs[j].word
	})

	resultLen := 10
	if len(pairs) < resultLen {
		resultLen = len(pairs)
	}

	result := make([]string, resultLen)
	for i := 0; i < resultLen; i++ {
		result[i] = pairs[i].word
	}

	return result
}
