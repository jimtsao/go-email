// folder formats header fields, folding lines when needed
package folder

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"mime"
	"strings"
	"unicode/utf8"
)

const maxLineLen = 78 // octets, excluding CRLF

const leastPriorityToken = math.MaxInt

var fwsToken = "\r\n "
var spaceLen = len(" ")

type Folder struct {
	Err     error // io.Writer error
	w       io.Writer
	written int           // length written to underlying io.Writer
	acc     []interface{} // accumulator
	closed  bool
}

// New returns folder that supports header folding.
// Max line length is set to 78 octets, excluding CRLF.
//
// Folder only has 1 line lookahead and respects priority
// of folding locations over minimising number of folded lines
func New(w io.Writer) *Folder {
	return &Folder{w: w, acc: []interface{}{}}
}

// Write expects a list of header string values and folding white space integer
// values, where lower int values takes precedence over higher values.
//
// The positioning of integer values signifies where folding may occur.
// If no int is specified, no folding will occur.
//
// (1) An ascii space preceded by an integer will be treated as an optionally foldable
// white space. For example::
//
//	e.Write("foo", 2, " ", "bar", 1, "baz") => foo bar\r\n baz
//	e.Write("foo", 1, "bar", "baz") => foobar\r\n baz
//
// (2) utf-8 encoded words i.e. =?utf-8?...?= may be decoded, split then re-encoded into multiple
// encoded word tokens where line limit is exceeded. Splitting only occurs as last resort.
//
// (3) use special token type WordEncodable to allow conditional word encoding. This is useful
// for strings that do not otherwise contain a foldable white space but can be word encoded
// to facilitate unlimited folding
func (f *Folder) Write(tokens ...interface{}) {
	// checks
	if f.Err != nil || f.closed {
		return
	}

	// push to accumulator
	dec := mime.WordDecoder{}
	for _, tok := range tokens {
		if f.Err != nil {
			return
		}

		switch v := tok.(type) {
		case int:
			f.acc = append(f.acc, v)
		case string:
			if dword, err := dec.DecodeHeader(v); err == nil && dword != v {
				// encoded words: split and append
				ewords := strings.Split(v, " ")
				for idx, eword := range ewords {
					if idx == 0 {
						f.acc = append(f.acc, eword)
					} else {
						f.acc = append(f.acc, leastPriorityToken, " ", eword)
					}
				}
			} else {
				f.acc = append(f.acc, v)
			}

			// fold if required
			f.fold()
		case WordEncodable:
			val := mime.QEncoding.Encode("utf-8", string(v))
			if val == string(v) {
				// us-ascii only chars
				f.acc = append(f.acc, v)
			} else {
				// requires encoding
				f.acc = append(f.acc, val)
			}
			f.fold()
		}
	}
}

// write all tokens in accumulator, then clear accumulator
func (f *Folder) flush() {
	var toWrite string
	for _, tok := range f.acc {
		switch v := tok.(type) {
		case string:
			toWrite += v
		case WordEncodable:
			toWrite += string(v)
		}
	}

	if _, f.Err = f.w.Write([]byte(toWrite)); f.Err != nil {
		return
	}

	f.written += len(toWrite)
	f.acc = []interface{}{}
}

// fold as many times as needed, consumes tokens from accumulator and tracks written length
func (f *Folder) fold() {
	// find first string where line length is exceeded, and
	// highest priority delimiter up to that token
	currentLen := f.written
	var exceededAt, delim int
	var needsFold, delimFound bool
	for idx, tok := range f.acc {
		switch v := tok.(type) {
		case int:
			if !delimFound {
				delimFound = true
				delim = v
			} else if v < delim {
				delim = v
			}
		case string, WordEncodable:
			if val, ok := v.(string); ok {
				currentLen += len(val)
			} else if val, ok := v.(WordEncodable); ok {
				currentLen += len(val)
			}
			if currentLen > maxLineLen {
				exceededAt = idx
				needsFold = true
				break
			}
		}
	}

	// line length not exceeded, no need to fold
	if !needsFold {
		return
	}

	// iterate backwards from the token that exceeds line limit, token by token
	// to find first suitable place where we can fold
	for i := len(f.acc[:exceededAt+1]) - 1; i >= 0; i-- {
		switch v := f.acc[i].(type) {
		case int:
			if v != delim || !f.canFold(i) {
				continue
			}

			// write parts before the delim
			oldAcc := f.acc
			newAcc := f.acc[i+1:]
			f.acc = append(f.acc[:i], fwsToken)

			if f.flush(); f.Err != nil {
				f.acc = oldAcc
				return
			}

			// set new accumulator
			f.acc = newAcc

			// if folded is white space, ignore the white space
			if len(f.acc) > 0 && f.acc[0] == " " {
				f.acc = f.acc[1:]
			}

			goto FOLDED
		case string, WordEncodable:
			// keep track of current len (written + len of strings up to index)
			if val, ok := v.(string); ok {
				currentLen -= len(val)
			} else if val, ok := v.(WordEncodable); ok {
				currentLen -= len(val)
			}

			// continue to next token until we find delimiter
			if delimFound {
				continue
			}

			// no delimiters to fold at, check if we can split the encoded word
			parts, didSplit := splitEncodedWord(v, maxLineLen-currentLen)
			if !didSplit {
				continue
			}

			// write all split parts except last one
			oldAcc := f.acc
			remAcc := f.acc[exceededAt+1:]
			for i := 0; i < len(parts); i++ {
				if i == 0 {
					// first part
					f.acc = f.acc[:exceededAt]
				} else if i == len(parts)-1 {
					// last part
					f.acc = append([]interface{}{parts[i]}, remAcc...)
					continue
				}

				f.acc = append(f.acc, parts[i], fwsToken)
				if f.flush(); f.Err != nil {
					f.acc = oldAcc
					return
				}
			}

			goto FOLDED
		}

		continue

	FOLDED:
		f.written = spaceLen

		// keep folding, tokens remaining in accumulator may still exceed max line length
		f.fold()
		return
	}
}

// CFWS MUST NOT be inserted in such a way that any line of a folded header
// field is made up entirely of WSP characters and nothing else.
// we check both left and right of index specified does not consist of only WSP
func (f *Folder) canFold(i int) bool {
	var lok, rok bool
	for idx := 0; idx < len(f.acc); idx++ {
		var val string
		if v, ok := f.acc[idx].(string); ok {
			val = v
		} else if v, ok := f.acc[idx].(WordEncodable); ok {
			val = string(v)
		}

		if strings.TrimLeft(val, "\t ") != "" {
			if idx < i {
				lok = true
				idx = i
			} else {
				rok = true
				break
			}
		}
	}
	return lok && rok
}

// Close flushes rest of buffered content and closes header
func (f *Folder) Close() {
	if f.closed || f.Err != nil {
		return
	}

	if f.flush(); f.Err != nil {
		return
	}

	if _, f.Err = f.w.Write([]byte("\r\n")); f.Err != nil {
		return
	}

	f.closed = true
}

// splits encoded word so first part is within specified octet limit
//
// max limit is determined by encoding where quoted-string is 75 octets
// and base64 by maxLineLen
func splitEncodedWord(word interface{}, limit int) ([]string, bool) {
	var dword string
	enc := mime.QEncoding
	if v, ok := word.(string); ok {
		// check if encoded word
		dec := &mime.WordDecoder{}
		dw, err := dec.Decode(v)
		if err != nil {
			return nil, false
		}
		dword = dw

		// check which encoding
		switch v[len("=?utf-8") : len("=?utf-8")+3] {
		case "?b?", "?B?":
			enc = mime.BEncoding
		}
	} else if v, ok := word.(WordEncodable); ok {
		dword = string(v)
	}

	// adjust limit if needed
	var maxContentLen int
	if enc == mime.QEncoding {
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
		return nil, false
	}

	// go rune by rune
	var runeLen int
	for i := 0; i < len(dword); i += runeLen {
		// figure out encoded length of rune
		var encLen int
		b := dword[i]
		if enc == mime.QEncoding {
			if b >= ' ' && b <= '~' && b != '=' && b != '?' && b != '_' {
				runeLen, encLen = 1, 1
			} else {
				_, runeLen = utf8.DecodeRuneInString(dword[i:])
				encLen = 3 * runeLen
			}
		} else {
			_, runeLen = utf8.DecodeRuneInString(dword[i:])
			encLen = runeLen
		}

		// split if this rune will exceed limit
		if encLen > maxContentLen {
			// unable to split even 1 rune and stay within limit
			if i == 0 {
				return nil, false
			}
			split := forceEncode(enc, dword[:i])
			rem := forceEncode(enc, dword[i:])
			parts := strings.Split(rem, " ")
			return append([]string{split}, parts...), true
		}

		// otherwise continue onto next rune
		maxContentLen -= encLen
	}

	// limit not reached, no need to split
	return nil, false
}

func forceEncode(enc mime.WordEncoder, val string) string {
	if enc == mime.QEncoding {
		val = enc.Encode("utf-8", "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"+val)
		return val[76:]
	}

	return fmt.Sprintf("=?utf-8?b?%s?=", base64.StdEncoding.EncodeToString([]byte(val)))
}
