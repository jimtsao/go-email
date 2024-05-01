package header_test

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestMIMEContentType(t *testing.T) {
	// text/plain
	h := &header.MIMEContentType{
		ContentType: "text/plain",
		Params:      map[string]string{"charset": "utf-8"},
	}
	assert.NoError(t, h.Validate())
	assert.Equal(t, "Content-Type: text/plain; charset=\"utf-8\"\r\n", h.String())

	// detect type
	h.DetectFromContent([]byte("<html>foo</html>"))
	assert.NoError(t, h.Validate())
	assert.Equal(t, "Content-Type: text/html; charset=\"utf-8\"\r\n", h.String())

	// octet-stream
	h = &header.MIMEContentType{
		ContentType: "text/html",
		Params: map[string]string{
			"charset": "utf-8",
			"name":    "foo.exe",
		},
	}
	assert.NoError(t, h.Validate())
	want := "Content-Type: text/html; charset=\"utf-8\"; name=\"foo.exe\"\r\n"
	got := h.String()
	assert.Equal(t, len(want), len(got))
	assert.Contains(t, want, "; charset=\"utf-8\"")
	assert.Contains(t, want, "; name=\"foo.exe\"")
}

func TestMIMEContentTransferEncoding(t *testing.T) {
	h := header.MIMEContentTransferEncoding("7bit")
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-Transfer-Encoding: 7bit\r\n"
	assert.Equal(t, want, h.String())
}

func TestMIMEContentID(t *testing.T) {
	h := header.MIMEContentID("<hello@127.0.0.1>")
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-ID: <hello@127.0.0.1>\r\n"
	assert.Equal(t, want, h.String())
}

func TestMIMEContentDisposition(t *testing.T) {
	sydney, _ := time.LoadLocation("Australia/Sydney")
	dt := time.Date(1990, time.April, 3, 5, 30, 15, 20, sydney)
	twant := "Tue, 3 Apr 1990 05:30:15 +1000"

	h := header.MIMEContentDisposition{
		Type:             header.Inline,
		Filename:         "foo.txt",
		CreationDate:     dt,
		Modificationdate: dt,
		ReadDate:         dt,
		Size:             1024,
	}
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-Disposition: inline" +
		"; filename=\"foo.txt\"" +
		"; creation-date=\"" + twant + "\"" +
		"; modification-date=\"" + twant + "\"" +
		"; read-date=\"" + twant + "\"" +
		"; size=1024" +
		"\r\n"
	assert.Equal(t, want, h.String())
}
