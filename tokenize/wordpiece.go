package tokenize

import (
	"github.com/sunhailin-Leo/gobert/tokenize/vocab"
	"strings"
)

// DefaultMaxWordChars is the max length of a token for it to be tokenized, otherwise marked as unknown
const DefaultMaxWordChars = 200

// DefaultUnknownToken is the token used to signify an unknown token
const DefaultUnknownToken = "[UNK]"

// Wordpiece is a tokenizer that breaks tokens into sub-word units based on a supplied vocabulary
// https://arxiv.org/pdf/1609.08144.pdf Section 4.1 for details
type Wordpiece struct {
	vocab        vocab.Dict
	maxWordChars int
	unknownToken string
}

// NewWordpiece returns a WordpieceTokenizer with the default settings.
// Generally should be used in a FullTokenizer
func NewWordpiece(voc vocab.Dict) Wordpiece {
	return Wordpiece{
		vocab:        voc,
		maxWordChars: DefaultMaxWordChars,
		unknownToken: DefaultUnknownToken,
	}
}

// Tokenize will segment the text into sub-word tokens from the supplied vocabulary
// NOTE: This implementation does not EXACTLY match the ref-impl and behaves slightly differently
// See https://github.com/google-research/bert/issues/763
func (wp Wordpiece) Tokenize(text string) (toks []string) {
	// TODO: determine if utf8 conversion is necessary, per python impl
	// text = convert_to_unicode(text)
	// Decrease a for-loop
	if strings.Index(text, " ") > 0 {
		for _, tok := range tokenizeWhitespaceV1(text) {
			toks = append(toks, wp.subTokenize(tok)...)
		}
		return toks
	}
	toks = wp.subTokenize(text)
	return toks
}

// subTokenize for Tokenize (impl from python)
func (wp Wordpiece) subTokenize(text string) (toks []string) {
	if len(text) > wp.maxWordChars {
		toks = append(toks, wp.unknownToken)
		return toks
	}
	start := 0
	for start < len(text) {
		end := len(text)
		subStrIsValid := false
		var curSubstr string
		for start < end {
			curSubstr = text[start:end]
			if wp.vocab.IsInVocab(curSubstr) {
				subStrIsValid = true
				break
			}
			end -= 1
		}
		if start > 0 {
			curSubstr = "##" + curSubstr
		}
		if !subStrIsValid {
			toks = append(toks, wp.unknownToken)
			break
		}
		toks = append(toks, curSubstr)
		start = end
	}
	return toks
}

// SetMaxWordChars will set the max chars for a word to be tokenized,
// generally this should be configured through the FullTokenizer
func (wp Wordpiece) SetMaxWordChars(c int) {
	wp.maxWordChars = c
}

// SetUnknownToken will set the unknown token, generally this should be configured through the FullTokenizer
func (wp Wordpiece) SetUnknownToken(tok string) {
	wp.unknownToken = tok
}
