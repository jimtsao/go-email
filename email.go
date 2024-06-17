package goemail

import (
	"github.com/jimtsao/go-email/header"
	"github.com/jimtsao/go-email/mime"
)

// Email is a wrapper around mime.Entity
//
// For full control use mime.Entity directly
type Email struct {
	From        string
	To          string // accepts comma-separated list
	Cc          string // accepts comma-separated list
	Bcc         string // accepts comma-separated list
	Subject     string // can contain any printable unicode characters
	Body        string
	Attachments []Attachment
	headers     []header.Header
}

func New() *Email {
	return &Email{}
}

// Validate checks syntax of headers for potential errors
// returns nil if no errors detected
func (e *Email) Validate() []error {
	hh := e.getHeaders()
	var errs []error
	for _, h := range hh {
		if err := h.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Raw produces RFC 5322 and MIME compliant email
func (e *Email) Raw() string {
	// body
	body := e.Body

	// headers
	hh := e.getHeaders()
	ct, cs := mime.DetectContentType([]byte(body))
	if cs == "" {
		hh = append(hh, header.NewContentType(ct, nil))
	} else {
		hh = append(hh, header.NewContentType(ct, header.NewMIMEParams("charset", cs)))
	}

	// inline attachments go into multipart/related
	// non-inline attachments go into multipart/mixed
	var inline, attached []*mime.Entity
	for _, att := range e.Attachments {
		if att.Inline {
			inline = append(inline, mime.NewEntity(nil, ""))
		} else {
			attached = append(attached, mime.NewEntity(nil, ""))
		}
	}

	ent := mime.NewEntity(hh, body)
	return ent.String()
}

func (e *Email) getHeaders() []header.Header {
	var hh []header.Header
	hh = append(hh, header.MIMEVersion{})
	if e.From != "" {
		hh = append(hh, header.Address{Field: header.AddressFrom, Value: e.From})
	}
	if e.To != "" {
		hh = append(hh, header.Address{Field: header.AddressTo, Value: e.To})
	}
	if e.Cc != "" {
		hh = append(hh, header.Address{Field: header.AddressCc, Value: e.Cc})
	}
	if e.Bcc != "" {
		hh = append(hh, header.Address{Field: header.AddressBcc, Value: e.Bcc})
	}
	if e.Subject != "" {
		hh = append(hh, header.Subject(e.Subject))
	}
	hh = append(hh, e.headers...)

	return hh
}

// AddHeader appends to generated headers. Headers
// are not guaranteed to be in any particular order
func (e *Email) AddHeader(h header.Header) {
	e.headers = append(e.headers, h)
}
