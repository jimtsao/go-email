package header

import (
	"fmt"
)

// CustomHeader represents an optional field header.
// WordEncoding can be optionally enabled for an
// Extension or user defined X-* header field
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
	FieldName    string
	Value        string
	WordEncoding bool
}

// Name returns header name
func (u CustomHeader) Name() string {
	return u.FieldName
}

func (u CustomHeader) Validate() error {
	if err := unstructured(u.Value).validate(u.WordEncoding); err != nil {
		return fmt.Errorf("%s: %w", u.Name(), err)
	}
	return nil
}

func (u CustomHeader) String() string {
	v := unstructured(u.Value).string(u.WordEncoding)
	return fmt.Sprintf("%s: %s\r\n", u.Name(), v)
}
