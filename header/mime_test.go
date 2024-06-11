package header_test

import (
	"fmt"
	"strings"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func headerContains(t *testing.T, header string, contains []string) {
	for _, h := range contains {
		before, after, found := strings.Cut(header, h)
		assert.Truef(t, found, "missing header: %s", h)
		header = before + after
	}
}

func TestMIMEContentType(t *testing.T) {
	// text/plain
	h := &header.MIMEContentType{
		ContentType: "text/plain",
		Params:      map[string]string{"charset": "utf-8"},
	}
	assert.NoError(t, h.Validate(), "text/plain")
	assert.Equal(t, "Content-Type: text/plain; charset=utf-8\r\n", h.String(), "text/plain")

	// detect type
	h.DetectFromContent([]byte("<html>foo</html>"))
	assert.NoError(t, h.Validate(), "detect type")
	assert.Equal(t, "Content-Type: text/html; charset=utf-8\r\n", h.String(), "detect type")

	// tspecial
	h = &header.MIMEContentType{
		ContentType: "multipart/mixed",
		Params: map[string]string{
			"boundary": "(foo)",
		},
	}
	assert.NoError(t, h.Validate(), "tspecials")
	assert.Equal(t, "Content-Type: multipart/mixed; boundary=\"(foo)\"\r\n", h.String(), "tspecials")

	// folding
	h = &header.MIMEContentType{
		ContentType: "text/html",
		Params: map[string]string{
			"charset": "utf-8",
			"param1":  strings.Repeat("i", 80),
			"param2":  strings.Repeat("i", 80),
		},
	}
	assert.NoError(t, h.Validate(), "folding")
	headerContains(t, h.String(), []string{
		"charset=utf-8",
		fmt.Sprintf(" param1=%s", strings.Repeat("i", 80)),
		fmt.Sprintf(" param2=%s", strings.Repeat("i", 80)),
	})
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
