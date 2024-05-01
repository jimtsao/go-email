package header

import (
	"fmt"
	"mime"

	"github.com/jimtsao/go-email/syntax"
)

// Subject represents the 'Date' header field
//
// Syntax:
//
//	subject         =   "Subject:" unstructured CRLF
//	unstructured    =   (*([FWS] VCHAR) *WSP)
type Subject string

func (s Subject) Name() string {
	return "Subject"
}

// Validate
//
// since any word encoding will produce printable ascii
// and satisfy 'unstructured' definition, we check that
// it can be word encoded instead
func (s Subject) Validate() error {
	if !syntax.IsWordEncodable(string(s)) {
		return fmt.Errorf("%s must contain only printable or white space characters", s.Name())
	}

	return nil
}

func (s Subject) String() string {
	sj := mime.QEncoding.Encode("utf-8", string(s))
	return fmt.Sprintf("%s: %s\r\n", s.Name(), sj)
}
