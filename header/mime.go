package header

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jimtsao/go-email/folder"
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

type MIMEParam struct {
	Attribute string
	Value     string
}

// NewMIMEParams returns mime param array. Input format:
//
//	attribute, value, attribute, value...
func NewMIMEParams(params ...string) []MIMEParam {
	mp := []MIMEParam{}
	for idx, v := range params {
		if idx%2 == 1 {
			mp = append(mp, MIMEParam{params[idx-1], v})
		}
	}
	return mp
}

// MIMEHeader represents a Content-[Name] header:
//
//	parameter        :=   attribute "=" value
//	attribute        :=   token
//	value            :=   token / quoted-string
//	token            :=   1*<any (US-ASCII) CHAR except SPACE, CTLs, or tspecials>
//	tspecials        :=   "(" / ")" / "<" / ">" / "@" /
//	                      "," / ";" / ":" / "\" / <"> /
//	                      "/" / "[" / "]" / "?" / "="
//	                      ; Must be in quoted-string to use within parameter values
type MIMEHeader struct {
	name     string
	val      string
	params   []MIMEParam
	validate func() error
}

// NewMIMEHeader returns a Content-[name] MIME header
func NewMIMEHeader(name string, val string, params []MIMEParam) *MIMEHeader {
	return &MIMEHeader{name: name, val: val, params: params}
}

func (m *MIMEHeader) Name() string {
	return fmt.Sprintf("Content-%s", m.name)
}

func (m *MIMEHeader) Validate() error {
	if m.validate == nil {
		return nil
	}
	return m.validate()
}

func (m *MIMEHeader) String() string {
	// fold using syntax: Content-name:[2][space]val param
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(m.Name()+":", 2, " ", m.val)

	// params
	for _, p := range m.params {
		mp := folder.MIMEParam{Name: p.Attribute, Val: p.Value}
		f.Write(";", 1, " ", mp)
	}

	f.Close()
	return sb.String()
}

// NewContentType returns Content-Type header:
//
//	content          :=   "Content-Type" ":" type "/" subtype
//	                      *(";" parameter)
//	type             :=   discrete-type / composite-type
//	discrete-type    :=   "text" / "image" / "audio" / "video" /
//	                      "application" / extension-token
//	composite-type   :=   "message" / "multipart" / extension-token
//	subtype          :=   extension-token / iana-token
func NewContentType(val string, params []MIMEParam) *MIMEHeader {
	return &MIMEHeader{name: "Type", val: val, params: params}
}

func NewContentTypeFrom(data []byte) *MIMEHeader {
	ct := http.DetectContentType(data)
	ct, cs, _ := strings.Cut(ct, "; charset=")
	return &MIMEHeader{
		name:   "Type",
		val:    ct,
		params: NewMIMEParams("charset", cs)}
}

// NewContentTransferEncoding returns 'Content-Transfer-Encoding' header:
//
//	encoding   :=  "Content-Transfer-Encoding" ":" mechanism
//	mechanism  :=  "7bit" / "8bit" / "binary" / "quoted-printable" /
//	               "base64" / ietf-token / x-token
func NewContentTransferEncoding(val string) *MIMEHeader {
	return &MIMEHeader{name: "Transfer-Encoding", val: val}
}

// NewContentID returns 'Content-ID' header:
//
//	content-id = "Content-ID" ":" msg-id
//	msg-id     = [CFWS] "<" id-left "@" id-right ">" [CFWS]
func NewContentID(val string) *MIMEHeader {
	validate := func() error {
		if err := msgid(val).validate(); err != nil {
			return fmt.Errorf("%s: %w", "Content-ID", err)
		}
		return nil
	}
	return &MIMEHeader{name: "ID", val: val, validate: validate}
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
func NewContentDisposition(inline bool, filename string, params []MIMEParam) *MIMEHeader {
	val := "attachment"
	if inline {
		val = "inline"
	}
	p := params
	if filename != "" {
		p = append(NewMIMEParams("filename", filename), params...)
	}
	return &MIMEHeader{name: "Disposition", val: val, params: p}
}
