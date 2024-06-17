package header

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jimtsao/go-email/folder"
	"github.com/jimtsao/go-email/syntax"
)

// MIMEVersion represents the 'MIME-Version' header
//
// default output is MIME-Version: 1.0
type MIMEVersion struct{}

func (m MIMEVersion) Name() string {
	return "MIME-Version"
}

func (m MIMEVersion) Validate() error {
	return nil
}

func (m MIMEVersion) String() string {
	return "MIME-Version: 1.0\r\n"
}

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
		if m.Params == nil {
			m.Params = map[string]string{}
		}
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
	// fold using syntax: Content-Type:[1] type/subtype;[1] param
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(m.Name()+":", 1, " ", m.ContentType)

	// params
	for attr, val := range m.Params {
		f.Write(";", 1, " ")
		mp := folder.MIMEParam{Name: attr, Val: val}
		if syntax.ContainsTSpecials(val) {
			mp.Val = `"` + mp.Val + `"`
			f.Write(mp)
		} else {
			f.Write(mp)
		}
	}

	f.Close()
	return sb.String()
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
	Inline           bool // inline vs attachment
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
	// fold using syntax: Content-Disposition:[1] type/subtype;[1] param
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(m.Name()+":", 1, " ")
	if m.Inline {
		f.Write("inline")
	} else {
		f.Write("attachment")
	}

	// params
	if m.Filename != "" {
		f.Write(";", 1, " ", fmt.Sprintf("filename=\"%s\"", m.Filename))
	}
	if !m.CreationDate.IsZero() {
		t := datetime(m.CreationDate)
		f.Write(";", 1, " ", fmt.Sprintf("creation-date=\"%s\"", t))
	}
	if !m.Modificationdate.IsZero() {
		t := datetime(m.Modificationdate)
		f.Write(";", 1, " ", fmt.Sprintf("modification-date=\"%s\"", t))
	}
	if !m.ReadDate.IsZero() {
		t := datetime(m.ReadDate)
		f.Write(";", 1, " ", fmt.Sprintf("read-date=\"%s\"", t))
	}
	if m.Size > 0 {
		f.Write(";", 1, " ", fmt.Sprintf("size=%d", m.Size))
	}

	f.Close()
	return sb.String()
}
