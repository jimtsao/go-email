package mime_test

import (
	"testing"

	"github.com/jimtsao/go-email/mime"
	"github.com/stretchr/testify/assert"
)

func multipartAlt() (*mime.Entity, string) {
	text := mime.NewEntity("text/plain", "foo bar", "charset", "us-ascii")
	html := mime.NewEntity("text/html", "<b>foo bar</b>", "charset", "utf-8")
	alt := mime.NewMultipartAlternative(nil, []*mime.Entity{text, html})
	want := "Content-Type: multipart/alternative; boundary=.*?\r\n" +
		"\r\n--.*?\r\n" +
		"Content-Type: text/plain; charset=us-ascii\r\n" +
		"\r\n" +
		"foo bar" +
		"\r\n--.*?\r\n" +
		"Content-Type: text/html; charset=utf-8" + "\r\n" +
		"\r\n" +
		"<b>foo bar</b>" +
		"\r\n--.*?--"
	return alt, want
}

func TestMultipart(t *testing.T) {
	alt, want := multipartAlt()
	got := alt.String()
	assert.Regexp(t, want, got)
}

func TestMultipartNested(t *testing.T) {
	alt, altWant := multipartAlt()
	mixed := mime.NewMultipartMixed(nil, []*mime.Entity{alt})
	got := mixed.String()
	want := "Content-Type: multipart/mixed; boundary=.*?\r\n" +
		"\r\n--.*?\r\n" +
		altWant +
		"\r\n--.*?--"
	assert.Regexp(t, want, got)
}
