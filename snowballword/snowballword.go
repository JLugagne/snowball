/*
This package defines a SnowballWord struct that is used
to encapsulate most of the "state" variables we must track
when stemming a word.  The SnowballWord struct also has
a few methods common to stemming in a variety of languages.
*/
package snowballword

import (
	"fmt"
	"unicode/utf8"
)

// SnowballWord represents a word that is going to be stemmed.
type SnowballWord struct {

	// A slice of runes
	RS []rune

	// The index in RS where the R1 region begins
	R1start int

	// The index in RS where the R2 region begins
	R2start int

	// The index in RS where the RV region begins
	RVstart int
}

// Create a new SnowballWord struct
func New(in string) SnowballWord {
	n := RuneLen(in)
	rs := make([]rune, n, n+8)
	i := 0
	for _, r := range in {
		rs[i] = r
		i++
	}
	return SnowballWord{
		RS:      rs,
		R1start: n,
		R2start: n,
		RVstart: n,
	}
}

// Reset reuses the existing backing array if it has enough capacity,
// avoiding a heap allocation. Returns the word ready for stemming.
func (w *SnowballWord) Reset(in string) {
	n := RuneLen(in)
	if cap(w.RS) >= n {
		w.RS = w.RS[:n]
	} else {
		w.RS = make([]rune, n, n+8)
	}
	i := 0
	for _, r := range in {
		w.RS[i] = r
		i++
	}
	w.R1start = n
	w.R2start = n
	w.RVstart = n
}

// Replace a suffix and adjust R1start and R2start as needed.
// If `force` is false, check to make sure the suffix exists first.
func (w *SnowballWord) ReplaceSuffix(suffix, replacement string, force bool) bool {
	if !force {
		foundSuffix := w.FirstSuffix(suffix)
		if foundSuffix != suffix {
			return false
		}
	}
	return w.ReplaceSuffixString(suffix, replacement, true)
}

// ReplaceSuffixString replaces a suffix with a replacement string without
// allocating []rune slices. If `force` is false, check to make sure the
// suffix exists first.
func (w *SnowballWord) ReplaceSuffixString(suffix, replacement string, force bool) bool {
	suffixLen := RuneLen(suffix)
	if !force && !w.hasSuffixStringIn(0, len(w.RS), suffix) {
		return false
	}
	w.RS = w.RS[:len(w.RS)-suffixLen]
	for _, r := range replacement {
		w.RS = append(w.RS, r)
	}
	w.resetR1R2()
	return true
}

// RuneLen returns the number of runes in s.
// For ASCII-only strings it avoids the overhead of utf8.RuneCountInString.
func RuneLen(s string) int {
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] >= 0x80 {
			return utf8.RuneCountInString(s)
		}
	}
	return n
}

// hasSuffixStringIn checks if w.RS[startPos:endPos] ends with the runes of suffix,
// and that the suffix starts at or after startPos.
func (w *SnowballWord) hasSuffixStringIn(startPos, endPos int, suffix string) bool {
	rsLen := endPos - startPos
	suffixRuneLen := RuneLen(suffix)
	if suffixRuneLen > rsLen {
		return false
	}
	j := endPos - suffixRuneLen
	for _, r := range suffix {
		if w.RS[j] != r {
			return false
		}
		j++
	}
	return true
}

// HasSuffixString returns true if `w` ends with the runes of `suffix`,
// without allocating a []rune slice.
func (w *SnowballWord) HasSuffixString(suffix string) bool {
	return w.hasSuffixStringIn(0, len(w.RS), suffix)
}

// Remove the last `n` runes from the SnowballWord.
func (w *SnowballWord) RemoveLastNRunes(n int) {
	w.RS = w.RS[:len(w.RS)-n]
	w.resetR1R2()
}

// Replace a suffix and adjust R1start and R2start as needed.
// If `force` is false, check to make sure the suffix exists first.
func (w *SnowballWord) ReplaceSuffixRunes(suffixRunes []rune, replacementRunes []rune, force bool) bool {

	if force || w.HasSuffixRunes(suffixRunes) {
		lenWithoutSuffix := len(w.RS) - len(suffixRunes)
		w.RS = append(w.RS[:lenWithoutSuffix], replacementRunes...)

		// If R, R2, & RV are now beyond the length
		// of the word, they are set to the length
		// of the word.  Otherwise, they are left
		// as they were.
		w.resetR1R2()
		return true
	}
	return false
}

// Resets R1start and R2start to ensure they
// are within bounds of the current rune slice.
func (w *SnowballWord) resetR1R2() {
	rsLen := len(w.RS)
	if w.R1start > rsLen {
		w.R1start = rsLen
	}
	if w.R2start > rsLen {
		w.R2start = rsLen
	}
	if w.RVstart > rsLen {
		w.RVstart = rsLen
	}
}

// Return a slice of w.RS, allowing the start
// and stop to be out of bounds.
func (w *SnowballWord) slice(start, stop int) []rune {
	startMin := 0
	if start < startMin {
		start = startMin
	}
	max := len(w.RS) - 1
	if start > max {
		start = max
	}
	if stop > max {
		stop = max
	}
	return w.RS[start:stop]
}

// Returns true if `x` runes would fit into R1.
func (w *SnowballWord) FitsInR1(x int) bool {
	return w.R1start <= len(w.RS)-x
}

// Returns true if `x` runes would fit into R2.
func (w *SnowballWord) FitsInR2(x int) bool {
	return w.R2start <= len(w.RS)-x
}

// Returns true if `x` runes would fit into RV.
func (w *SnowballWord) FitsInRV(x int) bool {
	return w.RVstart <= len(w.RS)-x
}

// Return the R1 region as a slice of runes
func (w *SnowballWord) R1() []rune {
	return w.RS[w.R1start:]
}

// Return the R1 region as a string
func (w *SnowballWord) R1String() string {
	return string(w.R1())
}

// Return the R2 region as a slice of runes
func (w *SnowballWord) R2() []rune {
	return w.RS[w.R2start:]
}

// Return the R2 region as a string
func (w *SnowballWord) R2String() string {
	return string(w.R2())
}

// Return the RV region as a slice of runes
func (w *SnowballWord) RV() []rune {
	return w.RS[w.RVstart:]
}

// Return the RV region as a string
func (w *SnowballWord) RVString() string {
	return string(w.RV())
}

// Return the SnowballWord as a string
func (w *SnowballWord) String() string {
	return string(w.RS)
}

func (w *SnowballWord) DebugString() string {
	return fmt.Sprintf("{\"%s\", %d, %d, %d}", w.String(), w.R1start, w.R2start, w.RVstart)
}

// Return the first prefix found or the empty string.
func (w *SnowballWord) FirstPrefix(prefixes ...string) (foundPrefix string) {
	rsLen := len(w.RS)

	for _, prefix := range prefixes {
		prefixLen := RuneLen(prefix)
		if prefixLen > rsLen {
			continue
		}

		found := true
		i := 0
		for _, r := range prefix {
			if w.RS[i] != r {
				found = false
				break
			}
			i++
		}
		if found {
			foundPrefix = prefix
			break
		}
	}
	return
}

// Return true if `w.RS[startPos:endPos]` ends with runes from `suffixRunes`.
// That is, the slice of runes between startPos and endPos have a suffix of
// suffixRunes.
func (w *SnowballWord) HasSuffixRunesIn(startPos, endPos int, suffixRunes []rune) bool {
	maxLen := endPos - startPos
	suffixLen := len(suffixRunes)
	if suffixLen > maxLen {
		return false
	}

	numMatching := 0
	for i := 0; i < maxLen && i < suffixLen; i++ {
		if w.RS[endPos-i-1] != suffixRunes[suffixLen-i-1] {
			break
		} else {
			numMatching += 1
		}
	}
	if numMatching == suffixLen {
		return true
	}
	return false
}

// Return true if `w` ends with `suffixRunes`
func (w *SnowballWord) HasSuffixRunes(suffixRunes []rune) bool {
	return w.HasSuffixRunesIn(0, len(w.RS), suffixRunes)
}

// Find the first suffix that ends at `endPos` in the word among
// those provided; then,
// check to see if it begins after startPos.  If it does, return
// it, else return the empty string and empty rune slice.  This
// may seem a counterintuitive manner to do this.  However, it
// matches what is required most of the time by the Snowball
// stemmer steps.
func (w *SnowballWord) FirstSuffixIfIn(startPos, endPos int, suffixes ...string) (suffix string) {
	for _, suffix := range suffixes {
		if w.hasSuffixStringIn(0, endPos, suffix) {
			suffixLen := RuneLen(suffix)
			if endPos-suffixLen >= startPos {
				return suffix
			}
			return ""
		}
	}

	return ""
}

func (w *SnowballWord) FirstSuffixIn(startPos, endPos int, suffixes ...string) (suffix string) {
	for _, suffix := range suffixes {
		if w.hasSuffixStringIn(startPos, endPos, suffix) {
			return suffix
		}
	}

	return ""
}

// Find the first suffix in the word among those provided; then,
// check to see if it begins after startPos.  If it does,
// remove it.
func (w *SnowballWord) RemoveFirstSuffixIfIn(startPos int, suffixes ...string) (suffix string) {
	suffix = w.FirstSuffixIfIn(startPos, len(w.RS), suffixes...)
	suffixLength := RuneLen(suffix)
	if suffix != "" {
		w.RemoveLastNRunes(suffixLength)
	}
	return
}

// Removes the first suffix found that is in `word.RS[startPos:len(word.RS)]`
func (w *SnowballWord) RemoveFirstSuffixIn(startPos int, suffixes ...string) (suffix string) {
	suffix = w.FirstSuffixIn(startPos, len(w.RS), suffixes...)
	suffixLength := RuneLen(suffix)
	if suffix != "" {
		w.RemoveLastNRunes(suffixLength)
	}
	return
}

// Removes the first suffix found
func (w *SnowballWord) RemoveFirstSuffix(suffixes ...string) (suffix string) {
	return w.RemoveFirstSuffixIn(0, suffixes...)
}

// Return the first suffix found or the empty string.
func (w *SnowballWord) FirstSuffix(suffixes ...string) (suffix string) {
	return w.FirstSuffixIfIn(0, len(w.RS), suffixes...)
}
