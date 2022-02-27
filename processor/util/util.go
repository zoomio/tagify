package util

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/zoomio/stopwords"
)

var (
	sanitizeRegex              = regexp.MustCompile(`([^\p{L}-']*)([\p{L}-']+)([^\p{L}-']*)`)
	notAWordRegex              = regexp.MustCompile(`([^\p{L}'-]+)`)
	noLetterWordRegex          = regexp.MustCompile(`[^\p{L}]`)
	doubleNotWordySymbolsRegex = regexp.MustCompile(`[^\p{L}]{2}`)
	punctuationRegex           = regexp.MustCompile(`[.,!;:]+`)

	newLine = []byte("\n")
)

// SplitToSentences splits given text into slice of sentences.
func SplitToSentences(text []byte) [][]byte {
	split := punctuationRegex.ReplaceAll(bytes.TrimSpace(text), newLine)
	return bytes.Split(split, newLine)
}

// Sanitize ...
func Sanitize(strs [][]byte, reg *stopwords.Register) []string {
	result := make([]string, 0)
	for _, s := range strs {
		str := string(s)
		// check if it is an URL
		if u, ok := isURL(str); ok && len(u.Hostname()) > 0 {
			str = strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
		} else {
			str = strings.ToLower(str)
		}
		// all letters to lower and with proper quote
		str = strings.Replace(str, "’", "'", -1)
		parts := notAWordRegex.Split(str, -1)
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

func isURL(s string) (*url.URL, bool) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, false
	}
	return u, true
}
