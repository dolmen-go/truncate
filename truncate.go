// Package truncate provides utilities to truncate strings.
package truncate

import (
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

// String truncates a string to the given limit in Unicode characters (not just bytes or codepoints).
func String(s string, limit int) string {
	// If the number of bytes is below the limit, this is identity
	if len(s) <= limit {
		return s
	}
	// Check if the start of the string below the limit is made of
	// single byte characters (pure ASCII)
	// Note: we also check the byte just after the limit as it may be a complement of
	// the last ASCII byte
	pureASCII := true
	endASCII := 0
	for i := 0; i <= limit; i++ {
		if s[i] >= utf8.RuneSelf {
			pureASCII = false
			if i > 0 {
				endASCII = i - 1
			}
			break
		}
	}
	if pureASCII {
		return s[:limit]
	}

	limit -= endASCII

	// Use big artillery: Unicode semantics
	var b strings.Builder

	if endASCII > 0 {
		b.WriteString(s[:endASCII])
	}
	g := uniseg.NewGraphemes(s[endASCII:])
	for g.Next() && limit > 0 {
		ru := g.Runes()
		for _, r := range ru {
			b.WriteRune(r)
		}
		limit--
	}

	return b.String()
}
