package spanish

import (
	"github.com/JLugagne/snowball/snowballword"
	"log"
	"strings"
)

func printDebug(debug bool, w *snowballword.SnowballWord) {
	if debug {
		log.Println(w.DebugString())
	}
}

// Stem an Spanish word.  This is the only exported
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
	step0(w)
	changeInStep1 := step1(w)
	if changeInStep1 == false {
		changeInStep2a := step2a(w)
		if changeInStep2a == false {
			step2b(w)
		}
	}
	step3(w)
	postprocess(w)
}
