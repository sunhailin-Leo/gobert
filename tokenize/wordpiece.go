package tokenize

import (
	"strings"

	"github.com/sunhailin-Leo/gobert/tokenize/vocab"
	"github.com/valyala/bytebufferpool"
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
	bufferPool   *bytebufferpool.ByteBuffer
}

// NewWordpiece returns a WordpieceTokenizer with the default settings.
// Generally should be used in a FullTokenizer
func NewWordpiece(voc vocab.Dict, buffer *bytebufferpool.ByteBuffer) Wordpiece {
	return Wordpiece{
		vocab:        voc,
		maxWordChars: DefaultMaxWordChars,
		unknownToken: DefaultUnknownToken,
		bufferPool:   buffer,
	}
}

// Tokenize will segment the text into sub-word tokens from the supplied vocabulary
// NOTE: This implementation does not EXACTLY match the ref-impl and behaves slightly differently
// See https://github.com/google-research/bert/issues/763
func (wp Wordpiece) Tokenize(text string) []string {
	defer func() { wp.bufferPool.Reset() }()
	if strings.Index(text, " ") > 0 {
		toks := make([]string, 0)
		for _, tok := range tokenizeWhitespaceV1(text) {
			if !wp.SubTokenize(tok) {
				toks = append(toks, strings.Fields(wp.bufferPool.String())...)
				wp.bufferPool.Reset()
				break
			}
			toks = append(toks, strings.Fields(wp.bufferPool.String())...)
			wp.bufferPool.Reset()
		}
		return toks
	}
	wp.SubTokenize(text)
	toks := strings.Fields(wp.bufferPool.String())
	return toks
}

// SubTokenize impl for old method
func (wp Wordpiece) SubTokenize(text string) bool {
	if wp.CheckIsLargeThanMaxWordChars(text) {
		return false
	}
	wp.CharLoop(text)
	return true
}

// CheckIsLargeThanMaxWordChars check text is larger than wp.maxWordChars
func (wp Wordpiece) CheckIsLargeThanMaxWordChars(text string) bool {
	if len(text) > wp.maxWordChars {
		wp.storeResult(wp.unknownToken)
		return true
	}
	return false
}

// CharLoop simplify logic and avoid slice memory leak
func (wp Wordpiece) CharLoop(text string) {
	for len(text) > 0 && text != "##" {
		sub := wp.vocab.LongestSubstring(text)
		if sub == "" {
			wp.storeResult(wp.unknownToken)
			break
		}
		wp.storeResult(sub)
		if len(text) == len(sub) {
			break
		}
		if text[len(sub):] == "" {
			break
		} else {
			text = "##" + text[len(sub):]
		}
	}
}

// storeResult store tokenize result
func (wp Wordpiece) storeResult(result string) {
	_, writeErr := wp.bufferPool.WriteString(result + "\n")
	if writeErr != nil {
		panic(writeErr)
	}
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
