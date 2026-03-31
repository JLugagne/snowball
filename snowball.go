package snowball

import (
	"fmt"
	"strings"

	"github.com/JLugagne/snowball/english"
	"github.com/JLugagne/snowball/french"
	"github.com/JLugagne/snowball/hungarian"
	"github.com/JLugagne/snowball/norwegian"
	"github.com/JLugagne/snowball/russian"
	"github.com/JLugagne/snowball/snowballword"
	"github.com/JLugagne/snowball/spanish"
	"github.com/JLugagne/snowball/swedish"
)

const (
	VERSION string = "v0.7.0"
)

// Stem a word in the specified language.
func Stem(word, language string, stemStopWords bool) (stemmed string, err error) {

	var f func(string, bool) string
	switch language {
	case "english":
		f = english.Stem
	case "spanish":
		f = spanish.Stem
	case "french":
		f = french.Stem
	case "russian":
		f = russian.Stem
	case "swedish":
		f = swedish.Stem
	case "norwegian":
		f = norwegian.Stem
	case "hungarian":
		f = hungarian.Stem
	default:
		err = fmt.Errorf("Unknown language: %s", language)
		return
	}
	stemmed = f(word, stemStopWords)
	return

}

// StemmerOption configures a Stemmer.
type StemmerOption func(*Stemmer)

// WithoutTrimSpace disables automatic whitespace trimming.
// Use this when the caller guarantees input has no leading/trailing whitespace.
func WithoutTrimSpace() StemmerOption {
	return func(s *Stemmer) { s.noTrimSpace = true }
}

// WithoutToLower disables automatic lowercasing.
// Use this when the caller guarantees input is already lowercase.
func WithoutToLower() StemmerOption {
	return func(s *Stemmer) { s.noToLower = true }
}

// Stemmer reuses internal buffers across calls, achieving zero heap
// allocations per stem when the output is consumed via StemRunes.
// A Stemmer is not safe for concurrent use.
type Stemmer struct {
	w           snowballword.SnowballWord
	language    string
	stemWord    func(*snowballword.SnowballWord)
	isStop      func(string) bool
	special     func(string) string // english only
	noTrimSpace bool
	noToLower   bool
}

// NewStemmer creates a reusable Stemmer for the given language.
// Options can be passed to skip ToLower/TrimSpace when the caller
// guarantees clean input, avoiding those allocations.
func NewStemmer(language string, opts ...StemmerOption) (*Stemmer, error) {
	s := &Stemmer{language: language}
	for _, opt := range opts {
		opt(s)
	}
	switch language {
	case "english":
		s.stemWord = english.StemWord
		s.isStop = english.IsStopWord
		s.special = english.StemSpecialWord
	case "spanish":
		s.stemWord = spanish.StemWord
		s.isStop = spanish.IsStopWord
	case "french":
		s.stemWord = french.StemWord
		s.isStop = french.IsStopWord
	case "russian":
		s.stemWord = russian.StemWord
		s.isStop = russian.IsStopWord
	case "swedish":
		s.stemWord = swedish.StemWord
		s.isStop = swedish.IsStopWord
	case "norwegian":
		s.stemWord = norwegian.StemWord
		s.isStop = norwegian.IsStopWord
	case "hungarian":
		s.stemWord = hungarian.StemWord
		s.isStop = hungarian.IsStopWord
	default:
		return nil, fmt.Errorf("Unknown language: %s", language)
	}
	return s, nil
}

func (s *Stemmer) prepare(word string) string {
	if !s.noTrimSpace {
		word = strings.TrimSpace(word)
	}
	if !s.noToLower {
		word = strings.ToLower(word)
	}
	return word
}

// Stem returns the stemmed word as a string. This allocates a string
// for the return value but reuses the internal rune buffer.
func (s *Stemmer) Stem(word string, stemStopWords bool) string {
	word = s.prepare(word)

	if len(word) <= 2 || (!stemStopWords && s.isStop(word)) {
		return word
	}

	if s.special != nil {
		if sv := s.special(word); sv != "" {
			return sv
		}
	}

	s.w.Reset(word)
	s.stemWord(&s.w)
	return s.w.String()
}

// StemRunes stems the word and returns the result as a rune slice.
// The returned slice is owned by the Stemmer and will be overwritten
// on the next call — copy it if you need to keep it.
// Returns (result, true) if stemming was performed, or (nil, false)
// if the word was short/stop and should be used as-is.
func (s *Stemmer) StemRunes(word string, stemStopWords bool) ([]rune, bool) {
	word = s.prepare(word)

	if len(word) <= 2 || (!stemStopWords && s.isStop(word)) {
		return nil, false
	}

	if s.special != nil {
		if sv := s.special(word); sv != "" {
			s.w.Reset(sv)
			return s.w.RS, true
		}
	}

	s.w.Reset(word)
	s.stemWord(&s.w)
	return s.w.RS, true
}
