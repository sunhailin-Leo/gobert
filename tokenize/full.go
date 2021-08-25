package tokenize

import "github.com/sunhailin-Leo/gobert/tokenize/vocab"

// Full is a FullTokenizer which comprises of a Basic & Wordpiece tokenizer
type Full struct {
	Basic     Basic
	Wordpiece Wordpiece
}

// Tokenize will tokenize the input text
// First basic is applited, then wordpiece on the tokens from basic
func (f Full) Tokenize(text string) (toks []string) {
	for _, tok := range f.Basic.Tokenize(text) {
		toks = append(toks, f.Wordpiece.Tokenize(tok)...)
	}
	return
}

// Vocab returns the vocab used for this tokenizer
func (f Full) Vocab() vocab.Dict {
	return f.Wordpiece.vocab
}
