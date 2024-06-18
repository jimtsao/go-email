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

func TestMIMEParam(t *testing.T) {
	for _, c := range []struct {
		desc string
		h    header.Header
		want string
	}{
		{"empty val", header.NewContentDisposition(false, "", nil), ""},
		{"empty quoted val", header.NewContentDisposition(false, `""`, nil), `; filename=""`},
		{"convert to quoted string", header.NewContentDisposition(false, `foo bar.txt`, nil), `; filename="foo bar.txt"`},
		{"convert to extended format", header.NewContentDisposition(false, "méow.txt", nil), "; filename*=utf-8''m%C3%A9ow.txt"},
	} {
		assert.NoError(t, c.h.Validate(), c.desc)
		want := fmt.Sprintf("Content-Disposition: attachment%s\r\n", c.want)
		assert.Equal(t, want, c.h.String(), c.desc)
	}

	// folding
	h := header.NewContentType("text/html",
		header.NewMIMEParams(
			"charset", "utf-8",
			"param_one", strings.Repeat("i", 80),
			"param_two", strings.Repeat("i", 80)))
	assert.NoError(t, h.Validate(), "folding")
	want := "Content-Type: text/html; charset=utf-8;" +
		"\r\n param_one*0*=utf-8''iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii" +
		"\r\n param_one*1*=iiiiiiiiiiiiiiiiiiiiiii;" +
		"\r\n param_two*0*=utf-8''iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii" +
		"\r\n param_two*1*=iiiiiiiiiiiiiiiiiiiiiii\r\n"
	assert.Equal(t, want, h.String())
}

func TestMIMEContentType(t *testing.T) {
	// text/plain
	h := header.NewContentType("text/plain", header.NewMIMEParams("charset", "us-ascii"))
	assert.NoError(t, h.Validate(), "text/plain")
	assert.Equal(t, "Content-Type: text/plain; charset=us-ascii\r\n", h.String(), "text/plain")

	// detect type
	h = header.NewContentTypeFrom([]byte("<html>foo</html>"))
	assert.NoError(t, h.Validate(), "detect type")
	assert.Equal(t, "Content-Type: text/html; charset=utf-8\r\n", h.String(), "detect type")

	// tspecial
	h = header.NewContentType("multipart/mixed", header.NewMIMEParams("boundary", "(foo)"))
	assert.NoError(t, h.Validate(), "tspecials")
	assert.Equal(t, "Content-Type: multipart/mixed; boundary=\"(foo)\"\r\n", h.String(), "tspecials")
}

func TestMIMEContentTransferEncoding(t *testing.T) {
	h := header.NewContentTransferEncoding("7bit")
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-Transfer-Encoding: 7bit\r\n"
	assert.Equal(t, want, h.String())
}

func TestMIMEContentID(t *testing.T) {
	h := header.NewContentID("<hello@127.0.0.1>")
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-ID: <hello@127.0.0.1>\r\n"
	assert.Equal(t, want, h.String())
}

func TestMIMEContentDisposition(t *testing.T) {
	sydney, _ := time.LoadLocation("Australia/Sydney")
	ctime := time.Date(1990, time.April, 1, 5, 30, 15, 20, sydney)
	cdate := "Sun, 1 Apr 1990 05:30:15 +1000"
	mtime := time.Date(1990, time.April, 2, 5, 30, 15, 20, sydney)
	mdate := "Mon, 2 Apr 1990 05:30:15 +1000"
	rtime := time.Date(1990, time.April, 3, 5, 30, 15, 20, sydney)
	rdate := "Tue, 3 Apr 1990 05:30:15 +1000"

	h := header.NewContentDisposition(false, "foo.txt",
		header.NewMIMEParams(
			"creation-date", ctime.Format(header.TimeRFC5322),
			"modification-date", mtime.Format(header.TimeRFC5322),
			"read-date", rtime.Format(header.TimeRFC5322),
			"size", "1024"))
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-Disposition: attachment; filename=foo.txt" +
		";\r\n creation-date=\"" + cdate + `"` +
		";\r\n modification-date=\"" + mdate + `"` +
		";\r\n read-date=\"" + rdate + `"` + "; size=1024" +
		"\r\n"
	assert.Equal(t, want, h.String())
}
