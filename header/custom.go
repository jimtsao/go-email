package header

import (
	"fmt"
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
	return nil
}

func (u CustomHeader) String() string {
	return fmt.Sprintf("%s: %s\r\n", u.Name(), u.Value)
}
