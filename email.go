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
	Attachments []*Attachment
	headers     []header.Header
}

func New() *Email {
	return &Email{}
}

// AddHeader appends to generated headers. Headers
// are not guaranteed to be in any particular order
func (e *Email) AddHeader(h header.Header) {
	e.headers = append(e.headers, h)
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
	// create body, inline and attachment entities
	var body *mime.Entity
	if e.Body != "" {
		var ctHeader header.Header
		ct, cs := mime.DetectContentType([]byte(e.Body))
		if cs == "" {
			ctHeader = header.NewContentType(ct, nil)
		} else {
			ctHeader = header.NewContentType(ct, header.NewMIMEParams("charset", cs))
		}

		body = mime.NewEntity([]header.Header{ctHeader}, e.Body)
	}

	var inline, attachments []*mime.Entity
	for _, att := range e.Attachments {
		if att.Inline {
			inline = append(inline, att.Entity())
		} else {
			attachments = append(attachments, att.Entity())
		}
	}

	// everything empty
	headers := e.getHeaders()
	if body == nil && inline == nil && attachments == nil {
		return mime.NewEntity(headers, "").String()
	}

	// single mime entity
	headers = append([]header.Header{header.MIMEVersion{}}, headers...)
	if body != nil && inline == nil && attachments == nil {
		body.Headers = append(headers, body.Headers...)
		return body.String()
	} else if body == nil && len(inline) == 1 && attachments == nil {
		inline[0].Headers = append(headers, inline[0].Headers...)
		return inline[0].String()
	} else if body == nil && inline == nil && len(attachments) == 1 {
		attachments[0].Headers = append(headers, attachments[0].Headers...)
		return attachments[0].String()
	}

	// multipart related
	var parts []*mime.Entity
	if attachments == nil && inline != nil {
		parts = inline
		if body != nil {
			parts = append([]*mime.Entity{body}, inline...)
		}
		related := mime.NewMultipartRelated(headers, parts)
		return related.String()
	}

	// multipart mixed
	if attachments != nil && len(inline) <= 1 {
		if body != nil && inline == nil {
			parts = append([]*mime.Entity{body}, attachments...)
		} else if body == nil && len(inline) == 1 {
			parts = append(inline, attachments...)
		} else {
			parts = attachments
		}
		mixed := mime.NewMultipartMixed(headers, parts)
		return mixed.String()
	}

	// multipart mixed > multipart related
	parts = inline
	if body != nil {
		parts = append([]*mime.Entity{body}, inline...)
	}
	related := mime.NewMultipartRelated(nil, parts)
	parts = append([]*mime.Entity{related}, attachments...)
	mixed := mime.NewMultipartMixed(headers, parts)
	return mixed.String()
}

func (e *Email) getHeaders() []header.Header {
	var hh []header.Header
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
