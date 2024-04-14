package textstats

import (
	"bufio"
	"fmt"
	"github.com/pontiki-ab/textstats/pkg/constants"
	"io"
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mtso/syllables"
)

var (
	pluralRegexp = regexp.MustCompile("([a-zA-Z]+?)(s\\b|\\b)")
)

type SyllableAnalysis struct {
	CMUSyllableCount             *int // http://www.speech.cs.cmu.edu/cgi-bin/cmudict
	DefaultAlgoSyllableCount     int
	AlternativeAlgoSyllableCount int // github.com/mtso/syllables
	MostLikelySyllableCount      int // CMU syllables if available, otherwise Alternative method, unless (see chooseDefaultAlgoOverAlternative)
}

type Word struct {
	Word        string
	Occurrences int
	Syllables   SyllableAnalysis
}

// Results is a struct containing the results of an analysis
type Results struct {
	CountTotalWords  int
	CountUniqueWords int
	Sentences        int
	Letters          int
	Punctuation      int
	Spaces           int
	Syllables        int
	DifficultWords   int

	WordCountPerSyllableCount       map[int]int
	UniqueWordCountPerSyllableCount map[int]int
	WordList                        map[string]*Word
}

// AverageLettersPerWord returns the average number of letters per word in the
// text
func (r *Results) AverageLettersPerWord() float64 {
	return float64(r.Letters) / float64(r.CountTotalWords)
}

// AverageSyllablesPerWord returns the average number of syllables per word in
// the text
func (r *Results) AverageSyllablesPerWord() float64 {
	return float64(r.Syllables) / float64(r.CountTotalWords)
}

// AverageWordsPerSentence returns the avergae number of words per sentence in
// the text
func (r *Results) AverageWordsPerSentence() float64 {
	if r.Sentences == 0 {
		return float64(r.CountTotalWords)
	}
	return float64(r.CountTotalWords) / float64(r.Sentences)
}

// WordsWithAtLeastNSyllables returns the number of words with at least N
// syllables, including or excluding proper nouns, in the text
func (r *Results) WordsWithAtLeastNSyllables(n int) int {
	var total int
	for sCount, wCount := range r.WordCountPerSyllableCount {
		if sCount >= n {
			total += wCount
		}
	}

	if total < 0 {
		return 0
	}

	return total
}

// PercentageWordsWithAtLeastNSyllables returns the percentage of words with at
// least N syllables, including or excluding proper nouns, in the text
func (r *Results) PercentageWordsWithAtLeastNSyllables(n int) float64 {
	return (float64(r.WordsWithAtLeastNSyllables(n)) / float64(r.CountTotalWords)) * 100.0
}

// FleschKincaidReadingEase returns the Flesch-Kincaid reading ease score for
// given text
func (r *Results) FleschKincaidReadingEase() float64 {
	return 206.835 - (1.015 * r.AverageWordsPerSentence()) - (84.6 * r.AverageSyllablesPerWord())
}

// FleschKincaidGradeLevel returns the Flesch-Kincaid grade level for the given text
func (r *Results) FleschKincaidGradeLevel() float64 {
	return (0.39 * r.AverageWordsPerSentence()) + (11.8 * r.AverageSyllablesPerWord()) - 15.59
}

// GunningFogScore returns the Gunning-Fog score for the given text
func (r *Results) GunningFogScore() float64 {
	return (r.AverageWordsPerSentence() + r.PercentageWordsWithAtLeastNSyllables(3)) * 0.4
}

// ColemanLiauIndex returns the Coleman-Liau index for the given text
func (r *Results) ColemanLiauIndex() float64 {
	sentences := float64(r.Sentences)
	if sentences == 0 {
		sentences = 1
	}

	return (5.89 * (float64(r.Letters) / float64(r.CountTotalWords))) - (0.3 * (sentences / float64(r.CountTotalWords))) - 15.8
}

// SMOGIndex returns the SMOG index for the given text
func (r *Results) SMOGIndex() float64 {
	sentences := float64(r.Sentences)
	if sentences == 0 {
		sentences = 1
	}

	return 1.0430 * math.Sqrt((float64(r.WordsWithAtLeastNSyllables(3))*(30/sentences))+3.1291)
}

// AutomatedReadabilityIndex returns the Automated Readability index for the given text
func (r *Results) AutomatedReadabilityIndex() float64 {
	sentences := float64(r.Sentences)
	if sentences == 0 {
		sentences = 1
	}

	return (4.71 * (float64(r.Letters) / float64(r.CountTotalWords))) + (0.5 * (float64(r.CountTotalWords) / sentences)) - 21.43
}

// DaleChallReadabilityScore returns the Dale-Chall readability score for the given text
func (r *Results) DaleChallReadabilityScore() float64 {
	difficultyPercentage := (float64(r.DifficultWords) / float64(r.CountTotalWords)) * 100

	sentences := float64(r.Sentences)
	if sentences == 0 {
		sentences = 1
	}

	score := (0.1579 * difficultyPercentage) + (0.0496 * (float64(r.CountTotalWords) / sentences))
	if difficultyPercentage > 5 {
		score += 3.6365
	}

	return score
}

func SyllableCount(word string) (sCount int) {
	word = strings.ToLower(word)

	// return early if we have a problem word
	sCount, ok := constants.ProblemWords[word]
	if ok {
		return
	}

	var prefixSuffixCount int
	for _, regex := range constants.PrefixSuffixes[:] {
		if regex.MatchString(word) {
			word = regex.ReplaceAllString(word, "")
			prefixSuffixCount++
		}
	}

	var wordPartCount int
	for _, wordPart := range constants.ConsonantsRegexp.Split(word, -1) {
		if len(wordPart) > 0 {
			wordPartCount++
		}
	}

	sCount = wordPartCount + prefixSuffixCount

	for _, regex := range constants.SubSyllables[:] {
		if regex.MatchString(word) {
			sCount--
		}
	}

	for _, regex := range constants.AddSyllables[:] {
		if regex.MatchString(word) {
			sCount++
		}
	}

	return
}

func SyllableCountAlternative(word string) (sCount int) {
	return syllables.In(word)
}

func analyseWord(word string, res *Results) {
	res.CountTotalWords++

	defaultAlgoSyllableCount := SyllableCount(word)
	alternativeAlgoSyllableCount := SyllableCountAlternative(word)
	mapKey := fmt.Sprintf("%d-%d", defaultAlgoSyllableCount, alternativeAlgoSyllableCount)

	var cmuSyllableCountPtr *int
	if cmuSyllableCount, ok := constants.CMUDictSyllableCountPerWord[word]; ok {
		cmuSyllableCountPtr = &cmuSyllableCount
	}

	var mostLikelySyllableCount int
	switch {
	case cmuSyllableCountPtr != nil: // CMU syllable count is most likely the correct one
		mostLikelySyllableCount = *cmuSyllableCountPtr
	case defaultAlgoSyllableCount == alternativeAlgoSyllableCount: // default and alternative algo agree, use that
		mostLikelySyllableCount = defaultAlgoSyllableCount
	case constants.ChooseDefaultAlgoOverAlternative[mapKey]: // default algo is statistically more often correct in a few cases (see var chooseDefaultAlgoOverAlternative)
		mostLikelySyllableCount = defaultAlgoSyllableCount
	}

	if _, ok := res.WordList[word]; ok {
		res.WordList[word].Occurrences++
	} else {
		res.WordList[word] = &Word{
			Word:        word,
			Occurrences: 1,
			Syllables: SyllableAnalysis{
				CMUSyllableCount:             cmuSyllableCountPtr,
				DefaultAlgoSyllableCount:     defaultAlgoSyllableCount,
				AlternativeAlgoSyllableCount: alternativeAlgoSyllableCount,
				MostLikelySyllableCount:      mostLikelySyllableCount,
			},
		}

		if _, ok := res.UniqueWordCountPerSyllableCount[mostLikelySyllableCount]; ok {
			res.UniqueWordCountPerSyllableCount[mostLikelySyllableCount]++
		} else {
			res.UniqueWordCountPerSyllableCount[mostLikelySyllableCount] = 1
		}

		res.CountUniqueWords++
	}

	res.Syllables += mostLikelySyllableCount

	if _, ok := res.WordCountPerSyllableCount[mostLikelySyllableCount]; ok {
		res.WordCountPerSyllableCount[mostLikelySyllableCount]++
	} else {
		res.WordCountPerSyllableCount[mostLikelySyllableCount] = 1
	}

	if _, ok := constants.DaleChallWordList[word]; !ok {
		matches := pluralRegexp.FindStringSubmatch(word)
		if len(matches) >= 2 {
			if _, ok := constants.DaleChallWordList[matches[1]]; !ok {
				res.DifficultWords++
			}
		} else {
			res.DifficultWords++
		}
	}
}

// Analyse scans a reader and outputs an analysis
func Analyse(r io.Reader) (res *Results, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	res = &Results{
		WordList:                        map[string]*Word{},
		WordCountPerSyllableCount:       map[int]int{},
		UniqueWordCountPerSyllableCount: map[int]int{},
	}

	var word string
	var endWord bool
	for scanner.Scan() {
		str := scanner.Text()
		letter, _ := utf8.DecodeRuneInString(str)
		switch {
		case unicode.IsLetter(letter):
			res.Letters++
			word += str
			endWord = false
		case unicode.IsSpace(letter):
			endWord = true
			res.Spaces++
		case unicode.IsPunct(letter):
			endWord = true
			switch str {
			case ".", "!", "?":
				res.Sentences++
			}
			res.Punctuation++
		}

		if endWord && len(word) > 0 {
			analyseWord(word, res)
			endWord = false
			word = ""
		}
	}

	if len(word) > 0 {
		analyseWord(word, res)
	}

	// Return scanner error if any
	err = scanner.Err()

	return
}
