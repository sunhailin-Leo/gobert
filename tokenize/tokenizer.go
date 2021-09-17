// Package tokenize supplies tokenization operations for BERT.
// Ports the tokenizer.py capability from the core BERT repo
//
// NOTE: All definitions are related to BERT and may vary from unicode definitions,
// for example, BERT considers '$' punctuation, but unicode does not.
package tokenize

import (
	"github.com/sunhailin-Leo/gobert/tokenize/vocab"
	"github.com/valyala/bytebufferpool"
)

// Tokenizer is an interface for chunking a string into it's tokens as per the BERT implementation
type Tokenizer interface {
	Tokenize(text string) []string
}

// VocabTokenizer comprises of a Tokenizer and VocabProvider
type VocabTokenizer interface {
	Tokenizer
	vocab.Provider
}

// NewTokenizer returns a new FullTokenizer
// Use Option array to modify default behavior
func NewTokenizer(voc vocab.Dict, buf *bytebufferpool.ByteBuffer, opts ...Option) VocabTokenizer {
	tkz := &Full{
		Basic:     NewBasic(),
		Wordpiece: NewWordpiece(voc, buf),
	}
	for _, opt := range opts {
		tkz = opt(tkz)
	}
	return tkz
}

// Option alter the behavior of the tokenizer
// TODO add tests for these behavior changes
type Option func(tkz *Full) *Full

// WithLower will lowercase all input if set to true, or skip lowering if false
// NOTE: kink from reference implementation is that lowering also strips accents
func WithLower(lower bool) Option {
	return func(tkz *Full) *Full {
		tkz.Basic.Lower = lower
		return tkz
	}
}

// WithUnknownToken will alter the unknown token from default [UNK]
func WithUnknownToken(unk string) Option {
	return func(tkz *Full) *Full {
		tkz.Wordpiece.unknownToken = unk
		return tkz
	}
}

// WithMaxChars sets the maximum len of a token to be tokenized, if longer will be labeled as unknown
func WithMaxChars(wc int) Option {
	return func(tkz *Full) *Full {
		tkz.Wordpiece.maxWordChars = wc
		return tkz
	}
}
