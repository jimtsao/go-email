package folder

import (
	"encoding/base64"
	"fmt"
	"mime"
	"unicode/utf8"
)

type wordEncodable struct {
	decoded    string
	enc        mime.WordEncoder
	mustEncode bool
}

// NewWordEncodable represents a managed optionally encodable string that handles
// folding at a customizable position. This is useful for folding an otherwise long
// string where a foldable white space may not be present.
// Non us-ascii will trigger encoding.
func NewWordEncodable(decoded string, encoder mime.WordEncoder, mustEncode bool) wordEncodable {
	return wordEncodable{decoded, encoder, mustEncode}
}

func (w wordEncodable) Value() string {
	return w.encode(w.decoded, w.mustEncode)
}

func (w wordEncodable) Length() int {
	return len(w.Value())
}

func (w wordEncodable) Fold(limit int) (split string, rest Foldable, didFold bool) {
	// length within limit
	if w.Length() <= limit {
		return w.Value(), nil, false
	}

	// adjust limit
	var maxContentLen int
	if w.enc == mime.QEncoding {
		// quoted-printable has max limit of 75, we cap limit to this
		// eg on a newly folded line, the limit can be 78 - 1 (whitespace)
		// which results in 2 encoded words where we only want 1
		if limit > 75 {
			limit = 75
		}
		maxContentLen = limit - len("=?utf-8?q?") - len("?=")
	} else {
		if limit > maxLineLen {
			limit = maxLineLen
		}
		maxContentLen = limit - len("=?utf-8?q?") - len("?=")
		maxContentLen = base64.StdEncoding.DecodedLen(maxContentLen)
	}

	// quick splittable check
	if maxContentLen <= 0 {
		return "", nil, false
	}

	// go rune by rune
	var runeLen int
	for i := 0; i < len(w.decoded); i += runeLen {
		// figure out encoded length of rune
		var encLen int
		b := w.decoded[i]
		if w.enc == mime.QEncoding {
			if b >= ' ' && b <= '~' && b != '=' && b != '?' && b != '_' {
				runeLen, encLen = 1, 1
			} else {
				_, runeLen = utf8.DecodeRuneInString(w.decoded[i:])
				encLen = 3 * runeLen
			}
		} else {
			_, runeLen = utf8.DecodeRuneInString(w.decoded[i:])
			encLen = runeLen
		}

		// split if this rune will exceed limit
		if encLen > maxContentLen {
			// unable to split even 1 rune and stay within limit
			if i == 0 {
				return "", nil, false
			}
			split := w.encode(w.decoded[:i], true)
			rest := wordEncodable{
				decoded:    w.decoded[i:],
				enc:        w.enc,
				mustEncode: true,
			}
			return split, rest, true
		}

		// otherwise continue onto next rune
		maxContentLen -= encLen
	}

	// limit not reached, no need to split
	return w.Value(), nil, false
}

func (w wordEncodable) encode(s string, force bool) string {
	if !force {
		return w.enc.Encode("utf-8", w.decoded)
	}

	// force encode
	if w.enc == mime.QEncoding {
		// extra fits exactly 1 full length 75 octet encoded word
		extra := "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"
		s = mime.QEncoding.Encode("utf-8", extra+s)
		return s[76:]
	}

	// base64
	s = base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("=?utf-8?b?%s?=", s)
}
