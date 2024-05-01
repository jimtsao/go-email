package header

// Date represents the 'Date' header field
//
// Syntax:
//
import (
	"errors"
	"fmt"

	"github.com/jimtsao/go-email/syntax"
)

// MessageID represents the 'Message-ID' header field.
//
// Automatically inserts @ and wraps in angle brackets if required
//
// Usage:
//
//	m := MessageID("<2024.12.31@localhost>")
//	m := MessageID("unique") // becomes <unique@>
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
	// folding not permitted within actual content of msg-id
	maxContentLen := maxLineLen - len(m.Name()+": ")
	id := msgid(m).String()
	if len(id) > maxContentLen {
		return fmt.Errorf("message-id must not exceed %d octets, has %d octets", maxContentLen, len(id))
	}

	// chars
	nameValid := IsValidHeaderName(m.Name())
	valValid := IsValidHeaderValue(id)
	if !nameValid && !valValid {
		return fmt.Errorf("%s: invalid characters in header name and body", m.Name())
	} else if !nameValid {
		return fmt.Errorf("%s: invalid characters in header name", m.Name())
	} else if !valValid {
		return fmt.Errorf("%s: invalid characters in header body", m.Name())
	}

	// syntax
	if !syntax.IsMsgID(id) {
		return errors.New("message-id invalid syntax")
	}

	return nil
}

func (m MessageID) String() string {
	id := msgid(m).String()
	return fmt.Sprintf("Message-ID: %s\r\n", id)
}
