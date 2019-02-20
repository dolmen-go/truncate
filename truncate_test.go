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
		// {"a", 0, "a"}, // just for easy show in VSCode terminal
	} {
		out := truncate.String(test.s, test.limit)
		if len(out) != len(test.expected) || norm.NFKD.String(out) != norm.NFKD.String(test.expected) {
			t.Errorf("%q, %d: got %q %04x, expecting %q %04x", test.s, test.limit, out, []rune(out), test.expected, []rune(test.expected))
		} else if out != test.expected {
			// no exact match, but NFKD match => just a warning
			t.Logf("%q, %d: got %q %04x, expecting %q %04x", test.s, test.limit, out, []rune(out), test.expected, []rune(test.expected))
		} else {
			t.Logf("%q, %d => %q", test.s, test.limit, test.expected)
		}
	}
}