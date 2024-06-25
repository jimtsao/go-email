// package base64 produces RFC 2045 compliant base64 encoding
package base64

import (
	"encoding/base64"
	"io"
)

// The encoded output stream must be represented
// in lines of no more than 76 characters each
const maxLineLen = 76

type b64Writer struct {
	currentLen int
	w          io.Writer
}

func NewEncoder(w io.Writer) io.WriteCloser {
	return base64.NewEncoder(base64.StdEncoding, &b64Writer{w: w})
}

func (b *b64Writer) Write(p []byte) (int, error) {
	rem := p
	written := 0
	var err error
	for {
		writableLen := maxLineLen - b.currentLen
		// within limit, write all of remaining
		if len(rem) <= writableLen {
			n, e := b.w.Write(rem)
			written += n
			err = e
			b.currentLen += len(rem)
			break
		}

		// exceeds limit, write what we can
		n, e := b.w.Write(rem[:writableLen])
		written += n
		err = e
		if err != nil {
			break
		}
		_, err = b.w.Write([]byte("\n"))
		if err != nil {
			break
		}

		b.currentLen = 0
		rem = rem[writableLen:]
	}

	return written, err
}
