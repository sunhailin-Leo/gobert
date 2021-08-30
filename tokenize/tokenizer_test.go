package tokenize_test

import (
	"fmt"
	"github.com/sunhailin-Leo/gobert/tokenize"
	"github.com/sunhailin-Leo/gobert/tokenize/vocab"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	for _, test := range []struct {
		name   string
		lower  bool
		text   string
		tokens []string
	}{
		{"chinese", false, "ah\u535A\u63A8zz", []string{"ah", "\u535A", "\u63A8", "zz"}},
		{"lower multi", true, " \tHeLLo!how  \n Are yoU?  ", []string{"hello", "!", "how", "are", "you", "?"}},
		{"lower single", true, "H\u00E9llo", []string{"hello"}},
		{"no lower multi", false, " \tHeLLo!how  \n Are yoU?  ", []string{"HeLLo", "!", "how", "Are", "yoU", "?"}},
		{"no lower single", false, "H\u00E9llo", []string{"H\u00E9llo"}},
	} {
		tkz := tokenize.Basic{Lower: test.lower}
		toks := tkz.Tokenize(test.text)
		if !reflect.DeepEqual(toks, test.tokens) {
			t.Errorf("Test %s - Invalid Tokenization - Want: %v, Got: %v", test.name, test.tokens, toks)
		}
	}
}

func TestWordpiece(t *testing.T) {
	voc := vocab.New([]string{"[UNK]", "[CLS]", "[SEP]", "want", "##want", "##ed", "wa", "un", "runn", "##ing"})
	for i, test := range []struct {
		text   string
		tokens []string
	}{
		{"", nil},
		{"unwanted", []string{"un", "##want", "##ed"}},
		{"unwanted running", []string{"un", "##want", "##ed", "runn", "##ing"}},
		// TODO determine if these tests are correct
		//	{"unwantedX", []string{"[UNK]"}},
		//{"unwantedX running", []string{"[UNK]", "runn", "##ing"}},
	} {
		tkz := tokenize.NewWordpiece(voc)
		toks := tkz.Tokenize(test.text)
		if !reflect.DeepEqual(toks, test.tokens) {
			t.Errorf("Test %d - Invalid Tokenization - Want: %v, Got: %v", i, test.tokens, toks)
		}
	}
}

func TestChineseTokenizer(t *testing.T) {
	voc, _ := vocab.FromFile("../export/vocab.txt")
	tkz := tokenize.NewTokenizer(voc)
	toks := tkz.Tokenize("广东省深圳市南山区人民政府")
	fmt.Println(toks)
	if !reflect.DeepEqual(toks, []string{"广", "东", "省", "深", "圳", "市", "南", "山", "区", "人", "民", "政", "府"}) {
		t.Errorf("Result is not equal")
	}
}

func BenchmarkChineseTokenizer(b *testing.B) {
	voc, _ := vocab.FromFile("../export/vocab.txt")
	tkz := tokenize.NewTokenizer(voc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tkz.Tokenize("广东省深圳市南山区人民政府")
	}
	b.ReportAllocs()
}

func TestChineseWordpieceTokenizer(t *testing.T) {
	voc, _ := vocab.FromFile("../export/vocab.txt")
	tkz := tokenize.NewWordpiece(voc)
	toks := tkz.Tokenize("广东省深圳市南山区人民政府")
	fmt.Println(toks)
	if !reflect.DeepEqual(toks, []string{"广", "##东", "##省", "##深", "##圳", "##市", "##南", "##山", "##区", "##人", "##民", "##政", "##府"}) {
		t.Errorf("Result is not equal")
	}
}

// Result1: BenchmarkChineseWordpieceTokenizer-12    	  207464	      5664 ns/op	     864 B/op	      19 allocs/op
// Result2: BenchmarkChineseWordpieceTokenizer-12    	  219673	      5336 ns/op	     832 B/op	      17 allocs/op
func BenchmarkChineseWordpieceTokenizer(b *testing.B) {
	voc, _ := vocab.FromFile("../export/vocab.txt")
	tkz := tokenize.NewWordpiece(voc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tkz.Tokenize("广东省深圳市南山区人民政府")
	}
	b.ReportAllocs()
}
