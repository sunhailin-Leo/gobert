package tokenize

import (
	"unicode"

	"strings"

	"golang.org/x/text/unicode/norm"
)

// Basic is a BasicTokenizer that run basic tokenize (punctuation splitting, lower casing, etc.).
type Basic struct {
	// Lower will apply a lower case filter to input
	Lower bool
}

// NewBasic returns a basic tokenizer. Method is supplied to match constructor of other tokenizers
func NewBasic() Basic {
	return Basic{Lower: true}
}

// Tokenize will segment a text into individual tokens. Follows algorithm from ref-imp
// Clean, PadChinese, Whitespace Split, Lower?, SplitPunc, Whitespace Split
func (bt Basic) Tokenize(text string) (toks []string) {
	// TODO assert text is unicode. text = unicode(text), from python impl
	//text = clean(text)
	//text = padChinese(text)
	for _, tok := range cleanAndPadChineseWithWhiteSpace(text) {
		if bt.Lower {
			tok = stripAccentsAndLower(tok)
		}
		toks = append(toks, splitPunctuation(tok)...)
	}
	// if white space is not in toks, it should return immediately
	//if isInStringArray(" ", toks) {
	//	toks = tokenizeWhitespace(strings.Join(toks, " "))
	//}
	return toks
}

// isInStringArray check a string data in string array
func isInStringArray(data string, array []string) bool {
	for _, item := range array {
		if item == data {
			return true
		}
	}
	return false
}

// clean function will clear some characters
func clean(text string) string {
	var b strings.Builder
	for _, c := range text {
		if c == 0 || c == 0xfffd || isControl(c) {
			continue
		} else if isWhitespace(c) {
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

// padChinese will add space padding around all CJK chars
// This implementation matches BasicTokenizer._tokenize_chinese_chars
func padChinese(text string) string {
	var b strings.Builder
	for _, c := range text {
		if isChinese(c) {
			b.WriteRune(' ')
			b.WriteRune(c)
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

// cleanAndPadChineseWithWhiteSpace combine three function clean, padChinese, tokenizeWhitespaceV1
func cleanAndPadChineseWithWhiteSpace(text string) []string {
	var b strings.Builder
	for _, c := range text {
		if c == 0 || c == 0xfffd || isControl(c) {
			continue
		} else if isChinese(c) {
			b.WriteRune(' ')
			b.WriteRune(c)
			b.WriteRune(' ')
		} else if isWhitespace(c) {
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return strings.Fields(strings.TrimSpace(b.String()))
}

// stripAccentsAndLower
func stripAccentsAndLower(text string) string {
	var b strings.Builder
	for _, c := range norm.NFD.String(text) {
		if !unicode.Is(unicode.Mn, c) {
			b.WriteRune(unicode.ToLower(c))
		}
	}
	return b.String()
}

// splitPunctuation
func splitPunctuation(text string) (toks []string) {
	// TODO test
	var b strings.Builder
	for _, c := range text {
		if isPunctuation(c) {
			toks = append(toks, b.String())
			toks = append(toks, string(c))
			b.Reset()
		} else {
			b.WriteRune(c)
		}
	}
	if b.Len() > 0 {
		toks = append(toks, b.String())
	}
	return
}

// tokenizeWhitespace splits text into tokens by whitespace, per python semantics empty strings are not included
func tokenizeWhitespace(text string) (toks []string) {
	for _, tok := range strings.Split(text, " ") {
		if tok != "" {
			toks = append(toks, tok)
		}
	}
	return toks
}

// tokenizeWhitespaceV1 splits text into tokens by whitespace, per python semantics empty strings are not included
func tokenizeWhitespaceV1(text string) []string {
	return strings.Fields(strings.TrimSpace(text))
}
