package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const topCount = 10

type WordCount struct {
	word  string
	count int
}

func Top10(text string) []string {
	splitText := strings.Fields(text)
	mapCount := make(map[string]int)

	for _, word := range splitText {
		mapCount[word]++
	}

	wordCounts := make([]WordCount, 0, len(mapCount))
	for word, count := range mapCount {
		wordCounts = append(wordCounts, WordCount{
			word:  word,
			count: count,
		})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}
		return wordCounts[i].count > wordCounts[j].count
	})

	result := make([]string, 0, topCount)

	for i, wordCount := range wordCounts {
		result = append(result, wordCount.word)
		if i == topCount-1 {
			break
		}
	}

	return result
}
