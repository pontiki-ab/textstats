package constants

import "regexp"

// statistical benchmark resulted in the following situations in which the default algorithm
// is statistically more often correct than the alternative algorithm (otherwise it's vice versa)
var ChooseDefaultAlgoOverAlternative = map[string]bool{
	"2-3": true, // default algo says 2 syllables, alternative says 3 syllables
	"3-4": true, // default algo says 3 syllables, alternative says 4 syllables
	"4-5": true, // default algo says 4 syllables, alternative says 5 syllables
	"5-4": true, // default algo says 5 syllables, alternative says 4 syllables
	"6-5": true, // default algo says 6 syllables, alternative says 5 syllables
	"6-7": true, // default algo says 6 syllables, alternative says 7 syllables
	"7-6": true, // default algo says 7 syllables, alternative says 6 syllables
}

var ConsonantsRegexp = regexp.MustCompile("[^aeiouy]+")

// ProblemWords are words that don't follow typical syllable counting rules
var ProblemWords = map[string]int{
	"simile":    3,
	"forever":   3,
	"shoreline": 2,
	"forest":    2,
}

// SubSyllables are syllables that would be counted as two but should be one
var SubSyllables = [...]*regexp.Regexp{
	regexp.MustCompile("cial"),
	regexp.MustCompile("tia"),
	regexp.MustCompile("cius"),
	regexp.MustCompile("cious"),
	regexp.MustCompile("giu"),
	regexp.MustCompile("ion"),
	regexp.MustCompile("ise"),
	regexp.MustCompile("iou"),
	regexp.MustCompile("sia$"),
	regexp.MustCompile("[^aeiouyt]{2,}ed$"),
	regexp.MustCompile(".ely$"),
	regexp.MustCompile("[cg]h?e[rsd]?$"),
	regexp.MustCompile("rved?$"),
	regexp.MustCompile("[aeiouy][dt]es?$"),
	regexp.MustCompile("[aeiouy][^aeiouydt]e[rsd]?$"),
	regexp.MustCompile("^[dr]e[aeiou][^aeiou]+$"),
	regexp.MustCompile("[aeiouy]rse$"),
}

// AddSyllables are syllables that would be counted as one but should be two
var AddSyllables = [...]*regexp.Regexp{
	regexp.MustCompile("ia"),
	regexp.MustCompile("riet"),
	regexp.MustCompile("dien"),
	regexp.MustCompile("iu"),
	regexp.MustCompile("io"),
	regexp.MustCompile("ii"),
	regexp.MustCompile("[aeiouym]bl$"),
	regexp.MustCompile("[aeiou]{3}"),
	regexp.MustCompile("^mc"),
	regexp.MustCompile("ism$"),
	regexp.MustCompile("[^aeiouy]{2}l$"),
	regexp.MustCompile("[^l]lien"),
	regexp.MustCompile("^coa[dglx]."),
	regexp.MustCompile("[^gq]ua[^auieo]"),
	regexp.MustCompile("dnt$"),
	regexp.MustCompile("uity$"),
	regexp.MustCompile("ie(r|st)$"),
	regexp.MustCompile("yee$"),
}

// PrefixSuffixes are single syllable prefixes and suffixes
var PrefixSuffixes = [...]*regexp.Regexp{
	regexp.MustCompile("^un"),
	regexp.MustCompile("^fore"),
	regexp.MustCompile("ly$"),
	regexp.MustCompile("less$"),
	regexp.MustCompile("ful$"),
	regexp.MustCompile("ers?$"),
	regexp.MustCompile("ings?$"),
}
