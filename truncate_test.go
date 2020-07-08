package truncate_test

import (
	"testing"

	"golang.org/x/text/unicode/norm"

	"github.com/dolmen-go/truncate"
)

func TestString(t *testing.T) {
	for _, test := range []struct {
		s        string
		limit    int
		expected string
	}{
		{"", 0, ""},
		{"a", 0, ""},
		{"é", 0, ""},
		{"a", 1, "a"},
		{"ab", 1, "a"},
		{"é", 1, "é"},
		{"éé", 1, "é"},
		{"\u00C0", 1, "\u00C0"},
		{"\u00C0\u00C0", 1, "\u00C0"},
		{"A\u0300", 1, "\u00C0"},
		{"A\u0300x", 1, "\u00C0"},
		{"A\u0300A\u0300x", 2, "\u00C0\u00C0"},
		{"A\u0300Ax", 2, "\u00C0A"},
		{"AA\u0300", 2, "A\u00C0"},
		{"AA\u0300x", 2, "A\u00C0"},
		{"x²", 1, "x"},
		{"x²*2", 2, "x²"},
		{"x\u304B\u3099", 2, "x\u304C"},
		{"e\u0301\u0303", 1, "é\u0303"},
		{"e\u0303\u0301", 1, "\u1ebd\u0301"},
		{"e\u0301\u0301\u0301", 1, "é\u0301\u0301"},      // https://blog.golang.org/normalization#TOC_5.
		{"é́́", 1, "é\u0301\u0301"},                     // https://blog.golang.org/normalization#TOC_5.
		{"\u01B5\u0327\u0308", 1, "\u01B5\u0327\u0308"},  // https://www.unicode.org/faq/char_combmark.html#11
		{"\u01B5\u0327\u0308é", 1, "\u01B5\u0327\u0308"}, // https://www.unicode.org/faq/char_combmark.html#11
		{"\u0061\u0328\u0301é", 1, "\u0105\u0301"},       // https://www.unicode.org/faq/char_combmark.html#12
		{"\u0061\u0328\u0301é", 1, "\u0105\u0301"},       // https://www.unicode.org/faq/char_combmark.html#12
		/*
			// TODO
			{"c\u034fho", 2, "c\u034fh"},                     // https://www.unicode.org/faq/char_combmark.html#17
			// TODO: U+035D COMBINING DOUBLE BREVE
		*/
		// {"a", 0, "fail"}, // a failure, just for easy show in VSCode terminal
	} {
		out := truncate.String(test.s, test.limit)
		if out == test.expected {
			t.Logf("%q, %d => %q", test.s, test.limit, test.expected)
		} else {
			printf := t.Errorf
			/*
				if test.TODO {
					printf = func(fmt string, args ...interface{}) {
						t.Logf("TODO "+fmt, args...)
					}
				} else
			*/
			if norm.NFKD.String(out) == norm.NFKD.String(test.expected) {
				// no exact match, but NFKD match => just a warning
				printf = t.Logf
			}
			printf("%q, %d: got %q %04x, expecting %q %04x", test.s, test.limit, out, []rune(out), test.expected, []rune(test.expected))
		}
	}
}
