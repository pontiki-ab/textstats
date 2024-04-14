package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
	file, err := os.Open("assets/cmudict-0.7b.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// loop over all lines
	results := map[string][]string{}
	for line, isPrefix, err := reader.ReadLine(); err == nil; line, isPrefix, err = reader.ReadLine() {
		if isPrefix {
			log.Fatal("buffer size to small")
		}

		// skip empty lines
		if len(line) == 0 {
			continue
		}

		// split line into words
		words := strings.Split(string(line), " ")

		// if first char doesn't start with a letter, skip this line (skipping comments, special characters, e.g. "END-QUOTE  EH1 N D K W OW1 T", etc.)
		if words[0][0] < 'A' || words[0][0] > 'Z' {
			continue
		}

		// if first word contains a "(" it's a pronunciation variant, skip this line
		if strings.Contains(words[0], "(") {
			continue
		}

		// add word and its phonemes to results
		results[words[0]] = words[1:]
	}

	// loop over all words and count syllables by their phonemes that constitute vowels
	wordsPerSyllable := map[int]int{}
	syllables := map[string]int{}
	for word, phonemes := range results {
		syllables[word] = 0
		for _, phoneme := range phonemes {
			if len(phoneme) == 0 {
				continue
			}

			// if a phoneme starts with a vowel, increment syllable count
			if slices.Contains([]rune{'A', 'E', 'I', 'O', 'U'}, rune(phoneme[0])) {
				syllables[word]++
			}
		}
		if _, ok := wordsPerSyllable[syllables[word]]; ok {
			wordsPerSyllable[syllables[word]]++
		} else {
			wordsPerSyllable[syllables[word]] = 1
		}
	}

	jsonString, err := json.MarshalIndent(syllables, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s\n", jsonString)
}
