package swedish

import (
	"strings"

	"github.com/JLugagne/snowball/snowballword"
)

// Stem a Swedish word. This is the only exported
// function in this package.
//
func Stem(word string, stemStopwWords bool) string {

	word = strings.ToLower(strings.TrimSpace(word))

	// Return small words and stop words
	if len(word) <= 2 || (stemStopwWords == false && IsStopWord(word)) {
		return word
	}

	w := snowballword.New(word)
	stemWord(&w)
	return w.String()

}

// StemWord stems w in place.
func StemWord(w *snowballword.SnowballWord) {
	stemWord(w)
}

func stemWord(w *snowballword.SnowballWord) {
	preprocess(w)
	step1(w)
	step2(w)
	step3(w)
}
