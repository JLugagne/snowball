Snowball
========


A [Go (golang)](http://golang.org) implementation of the
[Snowball stemmer](http://snowball.tartarus.org/)
for natural language processing.

Fork of [kljensen/snowball](https://github.com/kljensen/snowball) with
performance optimizations.


|                      |  Status                   |
| -------------------- | ------------------------- |
| Languages available  |  English, Spanish (español), French (le français), Russian (ру́сский язы́к), Swedish (svenska), Norwegian (norsk), Hungarian (magyar)|
| License              |  MIT                      |


## What changed in this fork

This fork adds a zero-allocation `Stemmer` API and several internal
optimizations. Compared to the original `Stem()` function:

| Metric | Original `Stem()` | `Stemmer.StemRunes()` | Change |
|--------|------------------:|----------------------:|-------:|
| English (26 words) | 23.25 µs | 11.01 µs | **-53%** |
| French (15 words) | 12.80 µs | 6.26 µs | **-51%** |
| Spanish (11 words) | 23.35 µs | 11.99 µs | **-49%** |
| Bytes/op | 688–1,856 | **0** | **-100%** |
| Allocs/op | 22–52 | **0** | **-100%** |

Changes:

* **`Stemmer` struct** -- reuses internal rune buffer across calls, eliminating
  the per-word `[]rune` allocation from `snowballword.New`.
* **`StemRunes`** -- returns the stemmed word as a `[]rune` slice (owned by the
  Stemmer), avoiding the `string()` return allocation. Achieves 0 allocs/op.
* **`Stemmer.Stem`** -- like the original but reuses buffers. 1 alloc/word (the
  returned string).
* **`WithoutToLower()` / `WithoutTrimSpace()`** -- options to skip normalization
  when the caller provides pre-normalized input.
* **`RuneLen` fast path** -- ASCII-optimized rune counting replaces
  `utf8.RuneCountInString` in all hot paths.
* **`SnowballWord.Reset`** -- reuses the backing `[]rune` array when capacity
  allows.
* **`StemWord` exported** from each language package for direct use with a
  caller-managed `SnowballWord`.
* Module renamed to `github.com/JLugagne/snowball`.
* Go version bumped to 1.25.


## Usage

### Drop-in replacement (same API as original)

```go
package main

import (
	"fmt"
	"github.com/JLugagne/snowball"
)

func main() {
	stemmed, err := snowball.Stem("Accumulations", "english", true)
	if err == nil {
		fmt.Println(stemmed) // Prints "accumul"
	}
}
```

### Zero-allocation API

Use `NewStemmer` with `StemRunes` to stem words in a tight loop with zero
heap allocations:

```go
package main

import (
	"fmt"
	"github.com/JLugagne/snowball"
)

func main() {
	// Create a reusable stemmer. Not safe for concurrent use.
	s, err := snowball.NewStemmer("english",
		snowball.WithoutToLower(),   // caller guarantees lowercase input
		snowball.WithoutTrimSpace(), // caller guarantees trimmed input
	)
	if err != nil {
		panic(err)
	}

	words := []string{"accumulations", "agreement", "skating"}
	for _, word := range words {
		// StemRunes returns a []rune owned by the Stemmer.
		// It is overwritten on the next call -- copy if needed.
		runes, stemmed := s.StemRunes(word, true)
		if stemmed {
			fmt.Println(string(runes))
		} else {
			// Word was too short or is a stop word; use as-is.
			fmt.Println(word)
		}
	}
}
```

If you need a `string` result but still want buffer reuse:

```go
s, _ := snowball.NewStemmer("english")
result := s.Stem("Accumulations", true) // 1 alloc (the returned string)
```


## Organization & Implementation

The code is organized as follows:

* The top-level `snowball` package has `snowball.Stem` (original API) and
  `snowball.NewStemmer` (zero-alloc API), defined in `snowball.go`.
* The stemmer for each language is defined in a "sub-package", e.g `english`.
* Each language exports `Stem` (allocating) and `StemWord` (in-place on a
  `*SnowballWord`).
* The `snowballword` package defines the `SnowballWord` struct with `New`,
  `Reset`, and suffix/prefix matching methods.
* Code common to multiple languages is in the `romance` package.


## Testing

To run the tests, do `go test ./...` in the top-level directory.


## Related work

Other stemmers available in Go:

* [stemmer](https://github.com/dchest/stemmer) by [Dmitry Chestnykh](https://github.com/dchest).
* [porter-stemmer](https://github.com/a2800276/porter-stemmer.go) - the original Porter algorithm.
* [go-stem](https://github.com/agonopol/go-stem) by [Alex Gonopolskiy](https://github.com/agonopol).
* [paicehusk](https://github.com/Rookii/paicehusk) by [Aaron Groves](https://github.com/rookii).
* [golibstemmer](https://github.com/rjohnsondev/golibstemmer) - Go bindings for libstemmer.


## License (MIT)

Copyright (c) the Contributors

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
