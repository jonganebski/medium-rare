package helper

import (
	"fmt"
	"math"
	"strings"
)

// ComputeReadTime computes read time of the story
func ComputeReadTime(totalText string) string {
	wordsCount := 0
	words := strings.Split(totalText, " ")
	for _, word := range words {
		if word == "" || strings.Contains(word, "=") || word == "()" || word == "(" || word == ")" || word == "{" || word == "}" || word == "<" || word == ">" {
			continue
		}
		wordsCount++
	}
	readTimeMinute := int(math.Ceil(float64(wordsCount) / 200))
	return fmt.Sprintf("%d min read", readTimeMinute)
}
