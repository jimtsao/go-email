package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestMIMEContentType(t *testing.T) {
	h := &header.MIMEContentType{
		ContentType: "text/plain",
		Charset:     "utf-8",
	}
	assert.NoError(t, h.Validate())
	assert.Equal(t, "Content-Type: text/plain; charset=utf-8\r\n", h.String())

	h.DetectFromContent([]byte("<html>foo</html>"))
	assert.NoError(t, h.Validate())
	assert.Equal(t, "Content-Type: text/html; charset=utf-8\r\n", h.String())
}

func TestMIMEEncoding(t *testing.T) {
	h := header.MIMEEncoding("7bit")
	err := h.Validate()
	assert.NoError(t, err)
	want := "Content-Transfer-Encoding: 7bit\r\n"
	assert.Equal(t, want, h.String())
}
