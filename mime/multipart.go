// sytax:
//
//	transport-padding := *LWSP-char
//	encapsulation     := delimiter transport-padding
//	                     CRLF body-part
//	delimiter         := CRLF dash-boundary
//	close-delimiter   := delimiter "--"
//	preamble          := discard-text
//	epilogue          := discard-text
//	discard-text      := *(*text CRLF) *text ; May be ignored or discarded.
//	body-part         := MIME-part-headers [CRLF *OCTET]
//	OCTET             := <any 0-255 octet value>
//
// note: composers MUST NOT generate non-zero length transport padding

package mime

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/jimtsao/go-email/header"
)

const (
	maxLineLen     = 78 // in octets
	maxBoundaryLen = 70
)

// boundary charset
var bcharnospace = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'()+_,-./:=?"
var bchars = bcharnospace + " "

// syntax:
//
//	boundary      := 0*69<bchars> bcharsnospace
//	bchars        := bcharsnospace / " "
//	bcharsnospace := DIGIT / ALPHA / "'" / "(" / ")" /
//	                 "+" / "_" / "," / "-" / "." /
//	                 "/" / ":" / "=" / "?"
func randomBoundary(n int) string {
	// length checks
	if n <= 0 {
		return ""
	} else if n > maxBoundaryLen {
		n = maxBoundaryLen
	}

	// random generation
	b := make([]byte, n-2)
	for i := 0; i < len(b)-1; i++ {
		b[i] = bchars[rand.Intn(len(bchars))]
	}
	b[len(b)-1] = bcharnospace[rand.Intn(len(bcharnospace))]

	return string(b)
}

// DetectContentType returns content type and charset if applicable
// It splits the two unlike http.DetectContentType
func DetectContentType(data []byte) (ctype string, charset string) {
	ct := http.DetectContentType(data)
	ct, cs, _ := strings.Cut(ct, "; charset=")
	return ct, cs
}

// syntax:
//
//	multipart-body := [preamble CRLF]
//	                  dash-boundary CRLF
//	                  body-part *encapsulation
//	                  close-delimiter
//	                  [CRLF epilogue]
//
// can think of syntax as basically:
//
//	multipart-body := [preamble CRLF]
//	                  dash-boundary CRLF body-part
//	                  *(CRLF dash-boundary CRLF body-part)
//	                  close-delimiter
//	                  [CRLF epilogue]
//
// the reason the first part is preceeded by dash-boundary rather than
// delimiter is because it is not preceeded by a CRLF. This is because
// with Entity syntax of header + blank line + body, so we don't end up
// with header + blank line + blank line + dash-boundary...
type multipartBody struct {
	boundary string
	parts    []*Entity
}

func (m *multipartBody) String() string {
	// 0:  dash-boundary CRLF body-part
	// 1+: delimiter CRLF body-part
	sb := strings.Builder{}
	for idx, body := range m.parts {
		content := "--" + m.boundary + "\r\n" + body.String()
		if idx == 0 {
			sb.WriteString(content)
		} else {
			sb.WriteString("\r\n" + content)
		}
	}

	// close-delimiter
	sb.WriteString("\r\n--" + m.boundary + "--")
	return sb.String()
}

func NewMultipartMixed(headers []header.Header, parts []*Entity) *Entity {
	return NewMultipart("mixed", headers, parts)
}

// NewMultipartAlternative returns multipart/alternative entity
// Parts should be in ascending (least to greatest) order of preference preference
func NewMultipartAlternative(headers []header.Header, parts []*Entity) *Entity {
	return NewMultipart("alternative", headers, parts)
}

func NewMultipartRelated(headers []header.Header, parts []*Entity) *Entity {
	return NewMultipart("related", headers, parts)
}

// NewMultipart returns an entity with content-type set as multipart/subtype
func NewMultipart(subtype string, headers []header.Header, parts []*Entity) *Entity {
	pre := fmt.Sprintf("Content-Type: multipart/%s; boundary=", subtype)
	boundary := randomBoundary(maxLineLen - len(pre))
	return &Entity{
		Headers: append(headers, header.NewContentType(
			"multipart/"+subtype,
			header.NewMIMEParams("boundary", boundary))),
		Body: &multipartBody{boundary: boundary, parts: parts}}
}
