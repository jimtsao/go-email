package header

import (
	"fmt"
	"net/http"
	"strings"
)

// MIME Entity Headers
//
// Syntax:
//
//		entity-headers   :=   [ content CRLF ]
//		                      [ encoding CRLF ]
//		                      [ id CRLF ]
//		                      [ description CRLF ]
//		                      *( MIME-extension-field CRLF )
//		content          :=   "Content-Type" ":" type "/" subtype
//		                      *(";" parameter)
//		discrete-type    :=   "text" / "image" / "audio" / "video" / "application"
//		composite-type   :=   "multipart" / "message"
//		subtype          :=   extension-token / iana-token
//		iana-token       :=   <A publicly-defined extension token. Tokens
//		                      of this form must be registered with IANA
//		                      as specified in RFC 2048.>

//		encoding         :=   "Content-Transfer-Encoding" ":" mechanism
//		mechanism        :=   "7bit" / "8bit" / "binary" /
//		                      "quoted-printable" / "base64" /
//		                      ietf-token / x-token
//		MIME-extension-field :=  <Any RFC 822 header field which begins with
//	                             the string "Content-">
//
// tspecials: must be in quoted-string to use within parameter values

// MIMEContentType represents the 'Content-Type' header
//
// Usage:
//
//	ct := MIMEContentType{ContentType: "text/plain", Charset: "utf-8"}
//	ct := MIMEContentType{}.DetectFromContent([]byte("<html>foo</html>"))
//
// Syntax:
//
//	content          :=   "Content-Type" ":" type "/" subtype
//	                      *(";" parameter)
//	type             :=   discrete-type / composite-type
//	discrete-type    :=   "text" / "image" / "audio" / "video" /
//	                      "application" / extension-token
//	composite-type   :=   "message" / "multipart" / extension-token
//	subtype          :=   extension-token / iana-token
//	parameter        :=   attribute "=" value
//	attribute        :=   token
//	value            :=   token / quoted-string
//	token            :=   1*<any (US-ASCII) CHAR except SPACE, CTLs, or tspecials>
//	tspecials        :=   "(" / ")" / "<" / ">" / "@" /
//	                      "," / ";" / ":" / "\" / <"> /
//	                      "/" / "[" / "]" / "?" / "="
type MIMEContentType struct {
	ContentType string
	Charset     string
}

func (m *MIMEContentType) DetectFromContent(data []byte) {
	ct := http.DetectContentType(data)
	m.ContentType, m.Charset, _ = strings.Cut(ct, "; charset=")
}

func (m *MIMEContentType) Name() string {
	return "Content-Type"
}

func (m *MIMEContentType) Validate() error {
	return nil
}

func (m *MIMEContentType) String() string {
	s := fmt.Sprintf("%s: %s", m.Name(), m.ContentType)
	if m.Charset != "" {
		s += fmt.Sprintf("; charset=%s", m.Charset)
	}
	s += "\r\n"
	return s
}

// MIMEEncoding represents the 'Content-Transfer-Encoding' header
//
// Usage:
//
//	m := MIMEEncoding("7bit")
//
// Syntax:
//
//	encoding   :=  "Content-Transfer-Encoding" ":" mechanism
//	mechanism  :=  "7bit" / "8bit" / "binary" /
//	               "quoted-printable" / "base64" /
//	               ietf-token / x-token
type MIMEEncoding string

func (m MIMEEncoding) Name() string {
	return "Content-Transfer-Encoding"
}

func (m MIMEEncoding) Validate() error {
	return nil
}

func (m MIMEEncoding) String() string {
	return fmt.Sprintf("%s: %s\r\n", m.Name(), string(m))
}
