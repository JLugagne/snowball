package snowball

import "testing"

var words = []string{
	"aberration", "abruptness", "absolute", "abated", "acclivity",
	"accumulations", "agreement", "breed", "skating", "fluently",
	"generously", "documentation", "internationally", "responsibilities",
	"accomplishments", "environmental", "unfortunately", "communication",
	"understanding", "comprehensive", "representative", "approximately",
	"consciousness", "extraordinary", "discrimination", "organizational",
}

var sink string
var sinkRunes []rune

func BenchmarkStemEnglish(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			sink, _ = Stem(w, "english", true)
		}
	}
}

func BenchmarkStemmerEnglish(b *testing.B) {
	s, _ := NewStemmer("english")
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			sink = s.Stem(w, true)
		}
	}
}

func BenchmarkStemmerRunesEnglish(b *testing.B) {
	s, _ := NewStemmer("english", WithoutToLower(), WithoutTrimSpace())
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			sinkRunes, _ = s.StemRunes(w, true)
		}
	}
}

func BenchmarkStemFrench(b *testing.B) {
	frenchWords := []string{
		"abandonnant", "abondamment", "aboutissement", "absolument",
		"accompagnement", "accroissement", "administration", "affaiblissement",
		"agrandissement", "alourdissement", "ambitieusement", "appartenance",
		"approvisionnement", "authentification", "bourgeoisement",
	}
	for i := 0; i < b.N; i++ {
		for _, w := range frenchWords {
			sink, _ = Stem(w, "french", true)
		}
	}
}

func BenchmarkStemmerRunesFrench(b *testing.B) {
	frenchWords := []string{
		"abandonnant", "abondamment", "aboutissement", "absolument",
		"accompagnement", "accroissement", "administration", "affaiblissement",
		"agrandissement", "alourdissement", "ambitieusement", "appartenance",
		"approvisionnement", "authentification", "bourgeoisement",
	}
	s, _ := NewStemmer("french", WithoutToLower(), WithoutTrimSpace())
	for i := 0; i < b.N; i++ {
		for _, w := range frenchWords {
			sinkRunes, _ = s.StemRunes(w, true)
		}
	}
}

func BenchmarkStemSpanish(b *testing.B) {
	spanishWords := []string{
		"lejana", "preocuparse", "oposición", "prisionero",
		"ridiculización", "cotidianeidad", "portezuela",
		"enriquecerse", "campesinos", "desalojó", "anticipadas",
	}
	for i := 0; i < b.N; i++ {
		for _, w := range spanishWords {
			sink, _ = Stem(w, "spanish", true)
		}
	}
}

func BenchmarkStemmerRunesSpanish(b *testing.B) {
	spanishWords := []string{
		"lejana", "preocuparse", "oposición", "prisionero",
		"ridiculización", "cotidianeidad", "portezuela",
		"enriquecerse", "campesinos", "desalojó", "anticipadas",
	}
	s, _ := NewStemmer("spanish", WithoutToLower(), WithoutTrimSpace())
	for i := 0; i < b.N; i++ {
		for _, w := range spanishWords {
			sinkRunes, _ = s.StemRunes(w, true)
		}
	}
}
