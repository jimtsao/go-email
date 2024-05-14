package header

import (
	"fmt"
	"strings"

	"github.com/jimtsao/go-email/folder"
)

// MessageID represents the 'Message-ID' header field.
//
// Usage:
//
//	m := MessageID("<2024.12.31@localhost>")
//
// Syntax:
//
//	message-id      =   "Message-ID:" msg-id CRLF
//	msg-id          =   [CFWS] "<" id-left "@" id-right ">" [CFWS]
//	id-left         =   dot-atom-text
//	id-right        =   dot-atom-text / no-fold-literal
//	no-fold-literal =   "[" *dtext "]"
//	dot-atom-text   =   1*atext *("." 1*atext)
//	atext           =   ALPHA / DIGIT /
//						"!" / "#" / "$" / "%" / "&" / "'" / "*" /
//						"+" / "-" / "/" / "=" / "?" / "^" / "_" /
//						"`" / "{" /	"|" / "}" / "~"
//	dtext           =   %d33-90 / %d94-126
//
// dtext: printable ascii excluding "[", "]", or "\"
type MessageID string

func (m MessageID) Name() string {
	return "Message-ID"
}

func (m MessageID) Validate() error {
	id := msgid(m)
	if err := id.validate(); err != nil {
		return fmt.Errorf("%s: %w", m.Name(), err)
	}

	// chars
	nameValid := IsValidHeaderName(m.Name())
	valValid := IsValidHeaderValue(id.string())
	if !nameValid && !valValid {
		return fmt.Errorf("%s: invalid characters in header name and body", m.Name())
	} else if !nameValid {
		return fmt.Errorf("%s: invalid characters in header name", m.Name())
	} else if !valValid {
		return fmt.Errorf("%s: invalid characters in header body", m.Name())
	}

	return nil
}

func (m MessageID) String() string {
	id := msgid(m).string()
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(m.Name()+":", 1, " ", id)
	f.Close()
	return sb.String()
}
