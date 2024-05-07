package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestSubject(t *testing.T) {
	for _, c := range []struct {
		input string
		want  string
		valid bool
	}{
		{"secret message", "Subject: secret message\r\n", true},
		{"éve is\tlistening", "Subject: =?utf-8?q?=C3=A9ve_is=09listening?=\r\n", true},
		{"\v\f", "Subject: =?utf-8?q?=0B=0C?=\r\n", false},
	} {
		h := header.Subject(c.input)
		err := h.Validate()
		if c.valid {
			assert.NoError(t, err, c.input)
		} else {
			assert.Error(t, err, c.input)
		}
		got := h.String()
		assert.Equal(t, c.want, got, c.input)
	}
}
