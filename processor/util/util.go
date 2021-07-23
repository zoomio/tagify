package util

import (
	"bytes"
	"regexp"
	"strings"
	"sync"

	"github.com/zoomio/stopwords"
)

var (
	sanitizeRegex              = regexp.MustCompile(`([^\p{L}-']*)([\p{L}-']+)([^\p{L}-']*)`)
	notAWordRegex              = regexp.MustCompile(`([^\p{L}'-]+)`)
	noLetterWordRegex          = regexp.MustCompile(`[^\p{L}]`)
	doubleNotWordySymbolsRegex = regexp.MustCompile(`[^\p{L}]{2}`)
	punctuationRegex           = regexp.MustCompile(`[.,!;:]+`)

	newLine = []byte("\n")

	// stop words
	once              sync.Once
	stopWordsLang     string
	stopWordsRegister *stopwords.Register
	allStopWords      = map[string]stopwords.Option{
		"en": stopwords.Words(stopwords.StopWords),
		"ru": stopwords.Words(stopwords.StopWordsRu),
		"zh": stopwords.Words(stopwords.StopWordsZh),
		"ja": stopwords.Words(stopwords.StopWordsJa),
		"ko": stopwords.Words(stopwords.StopWordsKo),
		"hi": stopwords.Words(stopwords.StopWordsHi),
		"he": stopwords.Words(stopwords.StopWordsHe),
		"ar": stopwords.Words(stopwords.StopWordsAr),
		"de": stopwords.Words(stopwords.StopWordsDe),
		"es": stopwords.Words(stopwords.StopWordsEs),
		"fr": stopwords.Words(stopwords.StopWordsFr),
	}
)

func SetStopWords(lang string) *stopwords.Register {
	once.Do(func() {
		stopWordsLang = lang
		stopWordsRegister = stopwords.Setup(allStopWords[lang])
	})
	return stopWordsRegister
}

func StopWords() *stopwords.Register {
	return stopWordsRegister
}

func StopWordsLang() string {
	return stopWordsLang
}

// SplitToSentences splits given text into slice of sentences.
func SplitToSentences(text []byte) [][]byte {
	split := punctuationRegex.ReplaceAll(bytes.TrimSpace(text), newLine)
	return bytes.Split(split, newLine)
}

// Sanitize ...
func Sanitize(strs [][]byte, reg *stopwords.Register) []string {
	result := make([]string, 0)
	for _, s := range strs {
		// all letters to lower and with proper quote
		s = bytes.ToLower(bytes.Replace(s, []byte("’"), []byte("'"), -1))
		parts := notAWordRegex.Split(string(s), -1)
		for _, p := range parts {
			normilized, ok := Normalize(p, reg)
			if !ok {
				continue
			}
			result = append(result, normilized)
		}
	}
	return result
}

// Normalize sanitizes word and tells whether it is allowed token or not.
func Normalize(word string, reg *stopwords.Register) (string, bool) {
	// False if doesn't match allowed regex
	if !sanitizeRegex.MatchString(word) {
		return word, false
	}

	// Remove not allowed symbols (sanitize)
	word = sanitizeRegex.ReplaceAllString(word, "${2}")

	// False if it is a stop word
	if reg != nil && reg.IsStopWord(word) {
		return word, false
	}

	// Defensive check if sanitized result is still not a word
	if notAWordRegex.MatchString(word) || doubleNotWordySymbolsRegex.MatchString(word) {
		return word, false
	}

	// Defensive check if word starts with hyphen
	if strings.HasPrefix(word, "-") {
		return word, false
	}

	if len(word) == 1 && noLetterWordRegex.MatchString(word) {
		return word, false
	}

	// Allowed word
	return word, true
}
