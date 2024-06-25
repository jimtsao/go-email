// folder formats header fields, folding lines when needed
package folder

import (
	"io"
	"strings"
)

const maxLineLen = 78 // octets, excluding CRLF

var fwsToken = "\r\n "

type Foldable interface {
	Value() string         // value before optional transformations
	Fold(limit int) string // folded output
	Priority() int         // folding priority, priority 0 is ignored
}

type Folder struct {
	Err     error // io.Writer error
	w       io.Writer
	written int           // current line length written
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
// Tokens can be int, string or Foldable
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
			if v != 0 {
				f.acc = append(f.acc, v)
			}
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

// fold as many times as needed, consumes tokens from
// accumulator and recalculates new written length
func (f *Folder) fold() {
	// find first token where line length is exceeded,
	// and highest priority delimiter up to that token
	currentLen := f.written
	var exceededAt, delim int

	for idx, tok := range f.acc {
		switch v := tok.(type) {
		case int:
			// register delimiter
			if delim == 0 {
				delim = v
			} else if v < delim {
				delim = v
			}
		case string:
			currentLen += len(v)
			if currentLen > maxLineLen {
				exceededAt = idx
				if f.foldAt(exceededAt, delim, currentLen) {
					return
				}
			}
		case Foldable:
			// register priority
			if v.Priority() != 0 {
				if delim == 0 {
					delim = v.Priority()
				} else if v.Priority() < delim {
					delim = v.Priority()
				}
			}

			// track length
			currentLen += len(v.Value())
			if currentLen > maxLineLen {
				exceededAt = idx
				if f.foldAt(exceededAt, delim, currentLen) {
					return
				}
			}
		}
	}

	// line length not exceeded, no need to fold
}

func (f *Folder) foldAt(pos int, delim int, currentLen int) bool {
	// iterate backwards from the token that exceeds line limit,
	// token by token to find first suitable place where we can fold
	for i := pos; i >= 0; i-- {
		switch v := f.acc[i].(type) {
		case int:
			if v != delim || !f.canFold(i, true, true) {
				continue
			}

			// write parts before the delim
			oldAcc := f.acc
			newAcc := f.acc[i+1:]
			f.acc = append(f.acc[:i], fwsToken)
			if f.flush(); f.Err != nil {
				f.acc = oldAcc
				break
			}

			// set new accumulator
			f.acc = newAcc

			// keep trying to fold
			f.written = len(" ")
			f.fold()
			return true
		case string:
			currentLen -= len(v)
		case Foldable:
			// keep track of current len (written + len of strings up to index)
			currentLen -= len(v.Value())

			// do we have folding priority
			if delim == 0 || v.Priority() != delim {
				continue
			}

			// attempt to fold
			output := v.Fold(maxLineLen - currentLen)

			// check for fwsToken to confirm fold occurred
			idxFirst := strings.Index(output, fwsToken)
			idxLast := strings.LastIndex(output, fwsToken)
			if idxLast == -1 {
				continue
			}

			// if pre or post fold string are all whitespace, we need
			// to run an additional whitespace check before proceeding
			lEmpty := strings.TrimLeft(output[:idxFirst], "\t ") == ""
			rEmpty := strings.TrimLeft(output[idxLast+len("\r\n"):], "\t ") == ""
			if lEmpty && !f.canFold(i, true, false) ||
				rEmpty && !f.canFold(i, false, true) {
				continue
			}

			// write the folded part
			oldAcc := f.acc
			newAcc := f.acc[i+1:]
			f.acc = append(f.acc[:i], output)
			if f.flush(); f.Err != nil {
				f.acc = oldAcc
				break
			}
			// set new written length
			lastToken := output[idxLast+len("\r\n"):]
			f.written = len(lastToken)

			// set remaining tokens as new accumulator
			f.acc = newAcc

			// keep folding
			f.written = len(" ")
			f.fold()
			return true
		}
	}

	return false
}

// CFWS MUST NOT be inserted in such a way that any line of a folded header
// field is made up entirely of WSP characters and nothing else.
// we check both left and right of index specified does not consist of only WSP
func (f *Folder) canFold(i int, checkLeft bool, checkRight bool) bool {
	// noop
	if !checkLeft && !checkRight {
		return true
	}

	// set check range
	lok, rok := !checkLeft, !checkRight
	from, to := 0, len(f.acc)
	if !checkLeft {
		from = i
	}
	if !checkRight {
		to = i
	}

	// run check
	for idx := from; idx < to; idx++ {
		if idx == i {
			continue
		}

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
