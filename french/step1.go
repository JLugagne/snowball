package french

import (
	"github.com/JLugagne/snowball/snowballword"
)

// Step 1 is the removal of standard suffixes
func step1(word *snowballword.SnowballWord) bool {
	suffix := word.FirstSuffix(
		"issements", "issement", "atrices", "utions", "usions", "logies",
		"emment", "ements", "atrice", "ations", "ateurs", "amment", "ution",
		"usion", "ments", "logie", "istes", "ismes", "iqUes", "euses",
		"ences", "ement", "ation", "ateur", "ances", "ables", "ment",
		"ités", "iste", "isme", "iqUe", "euse", "ence", "eaux", "ance",
		"able", "ives", "ité", "eux", "aux", "ive", "ifs", "if",
	)

	if suffix == "" {
		return false
	}
	suffixLength := snowballword.RuneLen(suffix)

	isInR1 := (word.R1start <= len(word.RS)-suffixLength)
	isInR2 := (word.R2start <= len(word.RS)-suffixLength)
	isInRV := (word.RVstart <= len(word.RS)-suffixLength)

	// Handle simple replacements & deletions in R2 first
	if isInR2 {

		// Handle simple replacements in R2
		repl := ""
		switch suffix {
		case "logie", "logies":
			repl = "log"
		case "usion", "ution", "usions", "utions":
			repl = "u"
		case "ence", "ences":
			repl = "ent"
		}
		if repl != "" {
			word.ReplaceSuffixString(suffix, repl, true)
			return true
		}

		// Handle simple deletions in R2
		switch suffix {
		case "ance", "iqUe", "isme", "able", "iste", "eux", "ances", "iqUes", "ismes", "ables", "istes":
			word.RemoveLastNRunes(suffixLength)
			return true
		}
	}

	// Handle simple replacements in RV
	if isInRV {

		// NOTE: these are "special" suffixes in that
		// we must still do steps 2a and 2b of the
		// French stemmer even when these suffixes are
		// found in step1.  Therefore, we are returning
		// `false` here.

		repl := ""
		switch suffix {
		case "amment":
			repl = "ant"
		case "emment":
			repl = "ent"
		}
		if repl != "" {
			word.ReplaceSuffixString(suffix, repl, true)
			return false
		}

		// Delete if preceded by a vowel that is also in RV
		if suffix == "ment" || suffix == "ments" {
			idx := len(word.RS) - suffixLength - 1
			if idx >= word.RVstart && isLowerVowel(word.RS[idx]) {
				word.RemoveLastNRunes(suffixLength)
				return false
			}
			return false
		}
	}

	// Handle all the other "special" cases.  All of these
	// return true immediately after changing the word.
	//
	switch suffix {
	case "eaux":

		// Replace with eau
		word.ReplaceSuffixString(suffix, "eau", true)
		return true

	case "aux":

		// Replace with al if in R1
		if isInR1 {
			word.ReplaceSuffixString(suffix, "al", true)
			return true
		}

	case "euse", "euses":

		// Delete if in R2, else replace by eux if in R1
		if isInR2 {
			word.RemoveLastNRunes(suffixLength)
			return true
		} else if isInR1 {
			word.ReplaceSuffixString(suffix, "eux", true)
			return true
		}

	case "issement", "issements":

		// Delete if in R1 and preceded by a non-vowel
		if isInR1 {
			idx := len(word.RS) - suffixLength - 1
			if idx >= 0 && isLowerVowel(word.RS[idx]) == false {
				word.RemoveLastNRunes(suffixLength)
				return true
			}
		}
		return false

	case "atrice", "ateur", "ation", "atrices", "ateurs", "ations":

		// Delete if in R2
		if isInR2 {
			word.RemoveLastNRunes(suffixLength)

			// If preceded by "ic", delete if in R2, else replace by "iqU".
			newSuffix := word.FirstSuffix("ic")
			if newSuffix != "" {
				newSuffixLen := snowballword.RuneLen(newSuffix)
				if word.FitsInR2(newSuffixLen) {
					word.RemoveLastNRunes(newSuffixLen)
				} else {
					word.ReplaceSuffixString(newSuffix, "iqU", true)
				}
			}
			return true
		}

	case "ement", "ements":

		if isInRV {

			// Delete if in RV
			word.RemoveLastNRunes(suffixLength)

			// If preceded by "iv", delete if in R2
			// (and if further preceded by "at", delete if in R2)
			newSuffix := word.RemoveFirstSuffixIfIn(word.R2start, "iv")
			if newSuffix != "" {
				word.RemoveFirstSuffixIfIn(word.R2start, "at")
				return true
			}

			// If preceded by "eus", delete if in R2, else replace by "eux" if in R1
			newSuffix = word.FirstSuffix("eus")
			if newSuffix != "" {
				newSuffixLen := snowballword.RuneLen(newSuffix)
				if word.FitsInR2(newSuffixLen) {
					word.RemoveLastNRunes(newSuffixLen)
				} else if word.FitsInR1(newSuffixLen) {
					word.ReplaceSuffixString(newSuffix, "eux", true)
				}
				return true
			}

			// If preceded by abl or iqU, delete if in R2, otherwise,
			newSuffix = word.FirstSuffix("abl", "iqU")
			if newSuffix != "" {
				newSuffixLen := snowballword.RuneLen(newSuffix)
				if word.FitsInR2(newSuffixLen) {
					word.RemoveLastNRunes(newSuffixLen)
				}
				return true
			}

			// If preceded by ièr or Ièr, replace by i if in RV
			newSuffix = word.FirstSuffix("ièr", "Ièr")
			if newSuffix != "" {
				if word.FitsInRV(snowballword.RuneLen(newSuffix)) {
					word.ReplaceSuffixString(newSuffix, "i", true)
				}
				return true
			}

			return true
		}

	case "ité", "ités":

		if isInR2 {

			// Delete if in R2
			word.RemoveLastNRunes(suffixLength)

			// If preceded by "abil", delete if in R2, else replace by "abl"
			newSuffix := word.FirstSuffix("abil")
			if newSuffix != "" {
				newSuffixLen := snowballword.RuneLen(newSuffix)
				if word.FitsInR2(newSuffixLen) {
					word.RemoveLastNRunes(newSuffixLen)
				} else {
					word.ReplaceSuffixString(newSuffix, "abl", true)
				}
				return true
			}

			// If preceded by "ic", delete if in R2, else replace by "iqU"
			newSuffix = word.FirstSuffix("ic")
			if newSuffix != "" {
				newSuffixLen := snowballword.RuneLen(newSuffix)
				if word.FitsInR2(newSuffixLen) {
					word.RemoveLastNRunes(newSuffixLen)
				} else {
					word.ReplaceSuffixString(newSuffix, "iqU", true)
				}
				return true
			}

			// If preceded by "iv", delete if in R2
			newSuffix = word.RemoveFirstSuffixIfIn(word.R2start, "iv")
			return true
		}
	case "if", "ive", "ifs", "ives":

		if isInR2 {

			// Delete if in R2
			word.RemoveLastNRunes(suffixLength)

			// If preceded by at, delete if in R2
			newSuffix := word.RemoveFirstSuffixIfIn(word.R2start, "at")
			if newSuffix != "" {

				// And if further preceded by ic, delete if in R2, else replace by iqU
				newSuffix = word.FirstSuffix("ic")
				if newSuffix != "" {
					newSuffixLen := snowballword.RuneLen(newSuffix)
					if word.FitsInR2(newSuffixLen) {
						word.RemoveLastNRunes(newSuffixLen)
					} else {
						word.ReplaceSuffixString(newSuffix, "iqU", true)
					}
				}
			}
			return true

		}
	}
	return false
}
