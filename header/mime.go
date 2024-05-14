package header

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jimtsao/go-email/folder"
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
	Params      map[string]string
}

func (m *MIMEContentType) DetectFromContent(data []byte) {
	ct := http.DetectContentType(data)
	ct, cs, _ := strings.Cut(ct, "; charset=")
	m.ContentType = ct
	if cs != "" {
		m.Params["charset"] = cs
	}
}

func (m *MIMEContentType) Name() string {
	return "Content-Type"
}

func (m *MIMEContentType) Validate() error {
	return nil
}

func (m *MIMEContentType) String() string {
	s := fmt.Sprintf("%s: %s", m.Name(), m.ContentType)
	for attr, val := range m.Params {
		s += fmt.Sprintf("; %s=\"%s\"", attr, val)
	}
	s += "\r\n"
	return s
}

// MIMEContentTransferEncoding represents the 'Content-Transfer-Encoding' header
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
type MIMEContentTransferEncoding string

func (m MIMEContentTransferEncoding) Name() string {
	return "Content-Transfer-Encoding"
}

func (m MIMEContentTransferEncoding) Validate() error {
	return nil
}

func (m MIMEContentTransferEncoding) String() string {
	return fmt.Sprintf("%s: %s\r\n", m.Name(), string(m))
}

// MIMEContentID represents the 'Content-ID' header
//
// Syntax:
//
//	content-id = "Content-ID" ":" msg-id
//	msg-id     = [CFWS] "<" id-left "@" id-right ">" [CFWS]
type MIMEContentID string

func (m MIMEContentID) Name() string {
	return "Content-ID"
}

func (m MIMEContentID) Validate() error {
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

func (m MIMEContentID) String() string {
	id := msgid(m).string()
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(m.Name()+":", 1, " ", id)
	f.Close()
	return sb.String()
}

type DispositionType string

const (
	Inline     DispositionType = "inline"
	Attachment DispositionType = "attachment"
)

// MIMEContentDisposition represents the 'Content-Disposition' header
//
// Syntax:
//
//	disposition      := "Content-Disposition" ":" disposition-type
//	                    *(";" disposition-parm)
//	disposition-type := "inline" / "attachment" / extension-token
//	                    ; values are not case-sensitive
//	disposition-parm := filename-parm /
//	                    creation-date-parm /
//	                    modification-date-parm /
//	                    read-date-parm /
//	                    size-parm /
//	                    parameter
//
//	filename-parm          := "filename" "=" value
//	creation-date-parm     := "creation-date" "=" quoted-date-time
//	modification-date-parm := "modification-date" "=" quoted-date-time
//	read-date-parm         := "read-date" "=" quoted-date-time
//	size-parm              := "size" "=" 1*DIGIT
//	quoted-date-time       := quoted-string
//	                          ; contents MUST be an RFC 822 `date-time'
//	                          ; numeric timezones (+HHMM or -HHMM) MUST be used
type MIMEContentDisposition struct {
	Type             DispositionType
	Filename         string
	CreationDate     time.Time
	Modificationdate time.Time
	ReadDate         time.Time
	Size             int // approximate size in octets
}

func (m MIMEContentDisposition) Name() string {
	return "Content-Disposition"
}

func (m MIMEContentDisposition) Validate() error {
	return nil
}

func (m MIMEContentDisposition) String() string {
	s := fmt.Sprintf("%s: %s", m.Name(), m.Type)
	if m.Filename != "" {
		s += fmt.Sprintf("; filename=\"%s\"", m.Filename)
	}
	if !m.CreationDate.IsZero() {
		t := datetime(m.CreationDate)
		s += fmt.Sprintf("; creation-date=\"%s\"", t)
	}
	if !m.Modificationdate.IsZero() {
		t := datetime(m.CreationDate)
		s += fmt.Sprintf("; modification-date=\"%s\"", t)
	}
	if !m.ReadDate.IsZero() {
		t := datetime(m.CreationDate)
		s += fmt.Sprintf("; read-date=\"%s\"", t)
	}
	if m.Size > 0 {
		s += fmt.Sprintf("; size=%d", m.Size)
	}
	s += "\r\n"

	return s
}
