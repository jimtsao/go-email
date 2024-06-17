package goemail_test

import (
	"testing"

	goemail "github.com/jimtsao/go-email"
	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestManagedEmail(t *testing.T) {
	m := goemail.New()
	m.From = "a@a.com"
	m.To = "b@b.com, c@c.com"
	m.Cc = "\"David\" <d@d.com>, Evé <e@e.com>"
	m.Bcc = "f@f.com"
	m.Subject = "Eve is éavesdropping"
	m.AddHeader(header.MessageID("<local@host.com>"))
	m.Body = "<b>attack at dawn</b>"
	want := "MIME-Version: 1.0\r\n" +
		"From: <a@a.com>\r\n" +
		"To: <b@b.com>,<c@c.com>\r\n" +
		"Cc: \"David\" <d@d.com>,=?utf-8?q?Ev=C3=A9?= <e@e.com>\r\n" +
		"Bcc: <f@f.com>\r\n" +
		"Subject: =?utf-8?q?Eve_is_=C3=A9avesdropping?=\r\n" +
		"Message-ID: <local@host.com>\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		"<b>attack at dawn</b>"
	got := m.Raw()
	assert.Equal(t, want, got)
}
