package english

import (
	"github.com/JLugagne/snowball/snowballword"
)

// Step 0 is to strip off apostrophes and "s".
func step0(w *snowballword.SnowballWord) bool {
	suffix := w.FirstSuffix("'s'", "'s", "'")
	if suffix == "" {
		return false
	}
	suffixLength := snowballword.RuneLen(suffix)
	w.RemoveLastNRunes(suffixLength)
	return true
}
