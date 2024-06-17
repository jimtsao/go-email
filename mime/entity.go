// Took me a while to wrap my head around so here is a bit of an explanation
// for future reference. Think of Entity as an expanded definition of an RFC 5322
// message. The basic syntax is: header + blank line + body.
//
// The header may or may not contain MIME related headers, i.e. Content-*
// A multipart entity is just one with a Content-Type: multipart/*; boundary=*
// header along with a special body syntax, but otherwise still retains the
// basic header + blank line + (multipart-body) syntax
package mime

import (
	"fmt"
	"strings"

	"github.com/jimtsao/go-email/header"
)

type String string

func (s String) String() string {
	return string(s)
}

// Entity refers to MIME-defined header fields and contents
// can be either message entity or multipart entity
//
// syntax:
//
//	entity-headers         :=   [ content CRLF ]
//	                            [ encoding CRLF ]
//	                            [ id CRLF ]
//	                            [ description CRLF ]
//	                            *( MIME-extension-field CRLF )
//	MIME-message-headers   :=   entity-headers
//	                            fields
//	                            version CRLF
//	MIME-part-headers      :=   entity-headers
//	                            [ fields ]
type Entity struct {
	Headers []fmt.Stringer
	Body    fmt.Stringer
}

// NewEntity provides a convenient interface for creating an Entity
//
// Usage:
//
//	NewEntity("text/plain", "hello world", "charset", "utf-8")
func NewEntity(ContentType string, Data string, Params ...string) *Entity {
	params := map[string]string{}
	var key string
	for idx, v := range Params {
		if idx%2 == 0 {
			key = v
			params[key] = ""
		} else {
			params[key] = v
		}
	}

	return &Entity{
		Headers: []fmt.Stringer{&header.MIMEContentType{
			ContentType: ContentType,
			Params:      params}},
		Body: String(Data),
	}
}

func (e *Entity) String() string {
	// header fields + blank line + body
	sb := strings.Builder{}
	for _, h := range e.Headers {
		sb.WriteString(h.String())
	}
	sb.WriteString("\r\n")
	sb.WriteString(e.Body.String())

	return sb.String()
}
