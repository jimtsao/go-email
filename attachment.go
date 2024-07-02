package goemail

import (
	"github.com/jimtsao/go-email/base64"
	"github.com/jimtsao/go-email/header"
	"github.com/jimtsao/go-email/mime"
)

type Attachment struct {
	Inline    bool // inline vs attachment
	Filename  string
	ContentID string // for inline ref, eg <img src="cid:[ContentID]" />
	Data      []byte
}

// Entity converts to mime.Entity form
func (a *Attachment) Entity() *mime.Entity {
	var ctype header.Header
	ct, cs := mime.DetectContentType(a.Data)
	if cs == "" {
		ctype = header.NewContentType(ct, nil)
	} else {
		ctype = header.NewContentType(ct, header.NewMIMEParams("charset", cs))
	}
	hh := []header.Header{
		ctype,
		header.NewContentDisposition(a.Inline, a.Filename, nil),
		header.NewContentTransferEncoding("base64"),
	}
	if a.ContentID != "" {
		hh = append(hh, header.NewContentID(a.ContentID))
	}

	b64Data := base64.EncodeToString(a.Data)
	return mime.NewEntity(hh, b64Data)
}
