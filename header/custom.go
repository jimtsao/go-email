package header

import (
	"fmt"
	"mime"

	"github.com/jimtsao/go-email/syntax"
)

// CustomHeader represents an optional field header
//
// Syntax:
//
//	optional-field  =   field-name ":" unstructured CRLF
//	field-name      =   1*ftext
//	ftext           =   %d33-57 / %d59-126
//	unstructured   =   (*([FWS] VCHAR) *WSP)
//
// ftext: printable ascii except colon
type CustomHeader struct {
	FieldName string
	Value     string
}

// Name returns canonical form of header name
func (u CustomHeader) Name() string {
	return CanonicalHeaderKey(u.FieldName)
}

func (u CustomHeader) Validate() error {
	if !syntax.IsWordEncodable(u.Value) {
		return fmt.Errorf("%s must contain only printable or white space characters", u.Name())
	}

	return nil
}

func (u CustomHeader) String() string {
	v := mime.QEncoding.Encode("utf-8", string(u.Value))
	return fmt.Sprintf("%s: %s\r\n", u.Name(), v)
}
