package folder

import (
	"encoding/base64"
	"fmt"
	"mime"
	"strings"
	"unicode/utf8"
)

// NewWordEncodable represents a managed optionally encodable string that handles
// folding at a customizable position. This is useful for folding an otherwise long
// string where a foldable white space may not be present.
// Non us-ascii will trigger encoding.
type WordEncodable struct {
	Decoded      string
	Enc          mime.WordEncoder
	MustEncode   bool
	FoldPriority int
}

func (w WordEncodable) Value() string {
	return w.encode(w.Decoded, w.MustEncode)
}

func (w WordEncodable) Priority() int {
	return w.FoldPriority
}

func (w WordEncodable) Fold(limit int) string {
	sb := strings.Builder{}
	remaining := w.Decoded

	// iterations of folding
ITERATE:
	for {
		// adjust limit
		var maxContentLen int
		if w.Enc == mime.QEncoding {
			// quoted-printable has max limit of 75 octets
			if limit > 75 {
				limit = 75
			}
			maxContentLen = limit - len("=?utf-8?q?") - len("?=")
		} else {
			if limit > maxLineLen {
				limit = maxLineLen
			}
			maxContentLen = limit - len("=?utf-8?b?") - len("?=")
			maxContentLen = base64.StdEncoding.DecodedLen(maxContentLen)
		}

		// quick foldable check
		if maxContentLen <= 0 {
			return ""
		}

		// go rune by rune
		var runeLen int
		for i := 0; i < len(remaining); i += runeLen {
			// figure out encoded length of rune
			var encLen int
			b := remaining[i]
			if w.Enc == mime.QEncoding {
				if b >= ' ' && b <= '~' && b != '=' && b != '?' && b != '_' {
					runeLen, encLen = 1, 1
				} else {
					_, runeLen = utf8.DecodeRuneInString(remaining[i:])
					encLen = 3 * runeLen
				}
			} else {
				_, runeLen = utf8.DecodeRuneInString(remaining[i:])
				encLen = runeLen
			}

			// fold now if this rune will exceed limit
			if encLen > maxContentLen {
				// unable to split even 1 rune and stay within limit
				if i == 0 {
					return ""
				}

				// write folded part
				split := w.encode(remaining[:i], true)
				sb.WriteString(split + fwsToken)

				// set remaining part
				remaining = remaining[i:]

				// reset limit then continue next folding iteration
				limit = maxLineLen
				continue ITERATE
			}

			// otherwise continue onto next rune
			maxContentLen -= encLen
		}

		// end of string reached, we are finished
		remaining = w.encode(remaining, true)
		sb.WriteString(remaining)
		break
	}

	return sb.String()
}

func (w WordEncodable) encode(s string, force bool) string {
	if !force {
		return w.Enc.Encode("utf-8", w.Decoded)
	}

	// force encode
	if w.Enc == mime.QEncoding {
		// extra fits exactly 1 full length 75 octet encoded word
		extra := "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"
		s = mime.QEncoding.Encode("utf-8", extra+s)
		return s[76:]
	}

	// base64
	s = base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("=?utf-8?b?%s?=", s)
}
