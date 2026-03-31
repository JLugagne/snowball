package english

import (
	"github.com/JLugagne/snowball/snowballword"
	"strings"
)

// Stem an English word.  This is the only exported
// function in this package.
//
func Stem(word string, stemStopwWords bool) string {

	word = strings.ToLower(strings.TrimSpace(word))

	// Return small words and stop words
	if len(word) <= 2 || (stemStopwWords == false && IsStopWord(word)) {
		return word
	}

	// Return special words immediately
	if specialVersion := StemSpecialWord(word); specialVersion != "" {
		word = specialVersion
		return word
	}

	w := snowballword.New(word)
	stemWord(&w)
	return w.String()

}

// StemWord stems w in place. The caller must have already lowercased
// and loaded the word into w. Returns true if the word was modified.
func StemWord(w *snowballword.SnowballWord) {
	stemWord(w)
}

func stemWord(w *snowballword.SnowballWord) {
	preprocess(w)
	step0(w)
	step1a(w)
	step1b(w)
	step1c(w)
	step2(w)
	step3(w)
	step4(w)
	step5(w)
	postprocess(w)
}
