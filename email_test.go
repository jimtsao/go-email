package goemail_test

import (
	"testing"

	goemail "github.com/jimtsao/go-email"
	"github.com/stretchr/testify/assert"
)

func TestEmailEncoding(t *testing.T) {
	// test address, subject and mime header param encoding
	m := &goemail.Email{
		From:    "a@a.com",
		To:      "\"David\" <d@d.com>, Evé <e@e.com>",
		Subject: "Eve is éavesdropping",
	}
	m.Attachments = []*goemail.Attachment{
		{Filename: "cát.pdf", Data: []byte("%PDF-1.7\n\nmeow")},
	}

	want := "MIME-Version: 1.0\r\n" +
		"From: <a@a.com>\r\n" +
		"To: \"David\" <d@d.com>,=\\?utf-8\\?q\\?Ev=C3=A9\\?= <e@e.com>\r\n" +
		"Subject: =\\?utf-8\\?q\\?Eve_is_=C3=A9avesdropping\\?=\r\n" +
		"Content-Type: application/pdf\r\n" +
		"Content-Disposition: attachment; filename\\*=utf-8''c%C3%A1t.pdf\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		"JVBERi0xLjcKCm1l"
	got := m.Raw()
	assert.Regexp(t, want, got)
}

func TestEmailComposition(t *testing.T) {
	m := goemail.New()
	m.From = "a@a.com"
	m.To = "b@b.com, c@c.com"
	m.Subject = "foo bar"

	// 0 body, 0 inline, 0 attachment
	wantHeader := "From: <a@a.com>\r\n" +
		"To: <b@b.com>,<c@c.com>\r\n" +
		"Subject: foo bar\r\n"
	want := wantHeader + "\r\n"
	assert.Equal(t, want, m.Raw(), "0 body, 0 inline, 0 attachments")

	// 1 body
	wantHeader = "MIME-Version: 1.0\r\n" + wantHeader
	body := "<b>hello world</b>"
	m.Body = body
	wantBody := "Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		m.Body
	want = wantHeader + wantBody
	assert.Equal(t, want, m.Raw(), "1 body")

	// 1 attachment
	m.Body = ""
	attachment := &goemail.Attachment{
		Filename: "cat.png",
		Data:     []byte("\x89PNG\x0D\x0A\x1A\x0A")}
	m.Attachments = []*goemail.Attachment{attachment}
	wantAttachment := "Content-Type: image/png\r\n" +
		"Content-Disposition: attachment; filename=cat.png\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		"iVBORw0K"
	want = wantHeader + wantAttachment
	assert.Equal(t, want, m.Raw(), "1 attachment")

	// 1 inline
	inline := &goemail.Attachment{
		Inline:    true,
		Filename:  "cat.png",
		ContentID: "<cat@png>",
		Data:      []byte("\x89PNG\x0D\x0A\x1A\x0A"),
	}
	m.Attachments = []*goemail.Attachment{inline}
	wantInline := "Content-Type: image/png\r\n" +
		"Content-Disposition: inline; filename=cat.png\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"Content-ID: <cat@png>\r\n" +
		"\r\n" +
		"iVBORw0K"
	want = wantHeader + wantInline
	assert.Equal(t, want, m.Raw(), "1 inline")

	// delimiters
	start := "--.*?\r\n"
	mid := "\r\n--.*?\r\n"
	end := "\r\n--.*?--"

	// multipart related
	// 2 inline
	m.Body = ""
	m.Attachments = []*goemail.Attachment{inline, inline}
	want = wantHeader +
		"Content-Type: multipart/related; boundary=.*?\r\n" +
		"\r\n" +
		start + wantInline +
		mid + wantInline +
		end
	assert.Regexp(t, want, m.Raw(), "2 inline")

	// 1 body, 1 inline
	m.Body = body
	m.Attachments = []*goemail.Attachment{inline}
	want = wantHeader +
		"Content-Type: multipart/related; boundary=.*?\r\n" +
		"\r\n" +
		start + wantBody +
		mid + wantInline +
		end
	assert.Regexp(t, want, m.Raw(), "1 body, 1 inline")

	// 1 body, 2 inline
	m.Body = body
	m.Attachments = []*goemail.Attachment{inline, inline}
	want = wantHeader +
		"Content-Type: multipart/related; boundary=.*?\r\n" +
		"\r\n" +
		start + wantBody +
		mid + wantInline +
		mid + wantInline +
		end
	assert.Regexp(t, want, m.Raw(), "1 body, 2 inline")

	// multipart mixed
	// 2 attachments
	m.Body = ""
	m.Attachments = []*goemail.Attachment{attachment, attachment}
	want = wantHeader +
		"Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n" +
		start + wantAttachment +
		mid + wantAttachment +
		end
	assert.Regexp(t, want, m.Raw(), "2 attachment")

	// 1 body, 1 attachment
	m.Body = body
	m.Attachments = []*goemail.Attachment{attachment}
	want = wantHeader +
		"Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n" +
		start + wantBody +
		mid + wantAttachment +
		end
	assert.Regexp(t, want, m.Raw(), "1 body, 1 attachment")

	// 1 body, 2 attachment
	m.Body = body
	m.Attachments = []*goemail.Attachment{attachment, attachment}
	want = wantHeader +
		"Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n" +
		start + wantBody +
		mid + wantAttachment +
		mid + wantAttachment +
		end
	assert.Regexp(t, want, m.Raw(), "1 body, 2 attachment")

	// 1 inline, 1 attachment
	m.Body = ""
	m.Attachments = []*goemail.Attachment{inline, attachment}
	want = wantHeader +
		"Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n" +
		start + wantInline +
		mid + wantAttachment +
		end
	assert.Regexp(t, want, m.Raw(), "1 inline, 1 attachment")

	// 1 inline, 2 attachment
	m.Body = ""
	m.Attachments = []*goemail.Attachment{inline, attachment, attachment}
	want = wantHeader +
		"Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n" +
		start + wantInline +
		mid + wantAttachment +
		mid + wantAttachment +
		end
	assert.Regexp(t, want, m.Raw(), "1 inline, 2 attachment")

	// multipart mixed > multipart related
	// 1 body, 1 inline, 1 attachment
	m.Body = body
	m.Attachments = []*goemail.Attachment{inline, attachment}
	want = wantHeader +
		"Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n" +
		start +
		"Content-Type: multipart/related; boundary=.*?\r\n" +
		"\r\n" +
		start + wantBody +
		mid + wantInline +
		end +
		mid + wantAttachment +
		end
	assert.Regexp(t, want, m.Raw(), "1 body, 1 inline, 1 attachment")
}
