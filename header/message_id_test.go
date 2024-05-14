package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestMessageID(t *testing.T) {
	// no folding
	m := header.MessageID("<a@b>")
	want := "Message-ID: <a@b>\r\n"
	assert.Equal(t, want, m.String(), "no folding")

	// folding
	m = header.MessageID("<i@iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii>")
	want = "Message-ID:\r\n <i@iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii>\r\n"
	assert.Equal(t, want, m.String(), "folding")
}
