package main

import (
	"encoding/json"
	"fmt"
	"github.com/pontiki-ab/textstats/pkg/textstats"
	"log"
	"os"
)

type Results struct {
	BothIncorrect  int
	OldAlgoCorrect map[string]int
	NewAlgoCorrect map[string]int
	TotalWords     int
}

func main() {
	wordListString, err := os.ReadFile("assets/cmudict-0.7b-syllable-count.json")
	if err != nil {
		log.Fatal(err)
	}

	wordListMap := map[string]int{}
	err = json.Unmarshal(wordListString, &wordListMap)
	if err != nil {
		log.Fatal(err)
	}

	results := Results{}
	for word, syllables := range wordListMap {
		if syllables == 0 {
			continue
		}

		results.TotalWords++
		if syllables != textstats.SyllableCount(word) && syllables != textstats.SyllableCountAlternative(word) {
			results.BothIncorrect++
		}

		if syllables == textstats.SyllableCount(word) && syllables != textstats.SyllableCountAlternative(word) {
			key := fmt.Sprintf("%d-%d", textstats.SyllableCount(word), textstats.SyllableCountAlternative(word))
			if results.OldAlgoCorrect == nil {
				results.OldAlgoCorrect = map[string]int{}
			}

			if _, ok := results.OldAlgoCorrect[key]; !ok {
				results.OldAlgoCorrect[key] = 1
			} else {
				results.OldAlgoCorrect[key]++
			}
		}

		if syllables != textstats.SyllableCount(word) && syllables == textstats.SyllableCountAlternative(word) {
			key := fmt.Sprintf("%d-%d", textstats.SyllableCount(word), textstats.SyllableCountAlternative(word))
			if results.NewAlgoCorrect == nil {
				results.NewAlgoCorrect = map[string]int{}
			}

			if _, ok := results.NewAlgoCorrect[key]; !ok {
				results.NewAlgoCorrect[key] = 1
			} else {
				results.NewAlgoCorrect[key]++
			}
		}
	}

	jsonString, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Result:\n%s\n", jsonString)
}
