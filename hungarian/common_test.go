package hungarian

import (
	"testing"

	"github.com/JLugagne/snowball/snowballword"
)

func TestFindRegions(t *testing.T) {
	for k, want := range map[string]int{
		"tóban":   2, //          consonant-vowel
		"ablakan": 2, //       vowel-consonant
		"acsony":  3, //         vowel-digraph
		"cvs":     3, //          null R1 region
	} {
		w := snowballword.New(k)
		got := findRegions(&w)
		if got != want {
			t.Errorf("%q: got %d, wanted %d", k, got, want)
		}
	}
}
