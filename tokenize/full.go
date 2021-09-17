package tokenize

import (
	"strings"
	"sync"

	"github.com/sunhailin-Leo/gobert/tokenize/vocab"
)

// Full is a FullTokenizer which comprises of a Basic & Wordpiece tokenizer
type Full struct {
	Basic     Basic
	Wordpiece Wordpiece

	tokenizeMutex sync.Mutex
}

// Tokenize will tokenize the input text
// First basic is applited, then wordpiece on the tokens from basic
func (f *Full) Tokenize(text string) []string {
	toks := make([]string, 0)
	f.tokenizeMutex.Lock()
	{
		for _, tok := range f.Basic.Tokenize(text) {
			if strings.Index(tok, " ") < 0 && tok != "" {
				toks = append(toks, strings.TrimSpace(f.Wordpiece.WordTokenize(tok)))
			}
		}
	}
	f.tokenizeMutex.Unlock()
	return toks
}

// Vocab returns the vocab used for this tokenizer
func (f *Full) Vocab() vocab.Dict {
	return f.Wordpiece.vocab
}
