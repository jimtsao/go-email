package header

import (
	"fmt"
	"strings"

	"github.com/jimtsao/go-email/folder"
	"github.com/jimtsao/go-email/syntax"
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
	FieldName     string
	Value         string
	WordEncodable bool
}

// Name returns header name
func (u CustomHeader) Name() string {
	return u.FieldName
}

func (u CustomHeader) Validate() error {
	if strings.Contains(u.Value, ":") {
		return fmt.Errorf("%s must not contain a colon", u.FieldName)
	}

	if u.WordEncodable && !syntax.IsWordEncodable(u.Value) {
		return fmt.Errorf("%s must contain only printable or white space characters", u.FieldName)
	}

	if !syntax.IsFtext(u.FieldName) {
		return fmt.Errorf("%s invalid syntax", u.FieldName)
	}
	return nil
}

func (u CustomHeader) String() string {
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(u.Name()+":", 1, " ")
	if u.WordEncodable {
		we := folder.WordEncodable(u.Value)
		f.Write(we)
	} else {
		f.Write(u.Value)
	}
	f.Close()
	return sb.String()
}
