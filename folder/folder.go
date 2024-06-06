// folder formats header fields, folding lines when needed
package folder

import (
	"io"
	"strings"
)

const maxLineLen = 78 // octets, excluding CRLF

var fwsToken = "\r\n "
var spaceLen = len(" ")

type Foldable interface {
	// value before optional transformations that may take place
	// upon folding
	Value() string // value before optional transformations
	Length() int   // length of Value
	Fold(limit int) (split string, rest Foldable, didFold bool)
}

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

// Write expects a list of header string values and folding white space
// integer values, where lower int values takes precedence over higher values.
//
// The positioning of integer values signifies where folding may occur.
// If no int is specified, no folding will occur.
//
// (1) An ascii space preceded by an integer will be treated as an optionally
// foldable white space. For example::
//
//	e.Write("foo", 2, " ", "bar", 1, "baz") => foo bar\r\n baz
//	e.Write("foo", 1, "bar", "baz") => foobar\r\n baz
//
// (2) accepts any token that conforms to Foldable interface
func (f *Folder) Write(tokens ...interface{}) {
	// checks
	if f.Err != nil || f.closed {
		return
	}

	// push to accumulator
	for _, tok := range tokens {
		if f.Err != nil {
			return
		}

		switch v := tok.(type) {
		case int:
			f.acc = append(f.acc, v)
		case string, Foldable:
			f.acc = append(f.acc, v)
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
		case Foldable:
			toWrite += v.Value()
		}
	}

	if _, f.Err = f.w.Write([]byte(toWrite)); f.Err != nil {
		return
	}

	f.written += len(toWrite)
	f.acc = []interface{}{}
}

// fold as many times as needed, consumes tokens from accumulator
// and recalculates new written length
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
		case string:
			currentLen += len(v)
			if currentLen > maxLineLen {
				exceededAt = idx
				needsFold = true
				break
			}
		case Foldable:
			currentLen += v.Length()
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

			// keep trying to fold
			f.written = spaceLen
			f.fold()
			return
		case string:
			currentLen -= len(v)
		case Foldable:
			// keep track of current len (written + len of strings up to index)
			currentLen -= v.Length()

			// continue to next token until we find delimiter
			if delimFound {
				continue
			}

			// no delimiters to fold at, try to fold current token
			split, rest, didFold := v.Fold(maxLineLen - currentLen)
			if !didFold {
				continue
			}

			// write the split part
			oldAcc := f.acc
			remAcc := f.acc[exceededAt+1:]
			f.acc = append(f.acc[:exceededAt], split)
			if rest != nil {
				f.acc = append(f.acc, fwsToken)
			}
			if f.flush(); f.Err != nil {
				f.acc = oldAcc
				return
			}

			// set remaining as new accumulator
			f.acc = append([]interface{}{rest}, remAcc...)

			// keep folding
			f.written = spaceLen
			f.fold()
			return
		}

	}
}

// CFWS MUST NOT be inserted in such a way that any line of a folded header
// field is made up entirely of WSP characters and nothing else.
// we check both left and right of index specified does not consist of only WSP
func (f *Folder) canFold(i int) bool {
	var lok, rok bool
	for idx := 0; idx < len(f.acc); idx++ {
		var val string
		switch v := f.acc[idx].(type) {
		case string:
			val = v
		case Foldable:
			val = v.Value()
		default:
			continue
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
