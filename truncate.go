// Package truncate provides utilities to truncate strings.
package truncate

import (
	"bytes"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
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
	for i := 0; i <= limit; i++ {
		if s[i] >= utf8.RuneSelf {
			pureASCII = false
			break
		}
	}
	if pureASCII {
		return s[:limit]
	}

	// Use big artillery: Unicode semantics
	var it norm.Iter
	it.InitString(norm.NFC, s)
	var b bytes.Buffer
	b.Grow(limit)
	for limit > 0 && !it.Done() {
		b.Write(it.Next())
		limit--
	}
	return b.String()
}
