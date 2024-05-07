package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestCustomHeader(t *testing.T) {
	for _, c := range []struct {
		name     string
		val      string
		encoding bool
		want     string
		pass     bool
	}{
		{name: "X-Foo", val: "bar", want: "X-Foo: bar\r\n", pass: true},
		{name: "x-fOO", val: "bar", want: "x-fOO: bar\r\n", pass: true}, // no canonical header
		{name: "X-Foo", val: "tést", encoding: true, want: "X-Foo: =?utf-8?q?t=C3=A9st?=\r\n", pass: true},
		{name: "X-Colon", val: "t:est", want: "X-Colon: t:est\r\n", pass: false},
		{name: "X-Colon", val: "t:est", encoding: true, want: "X-Colon: t:est\r\n", pass: false},
	} {
		h := header.CustomHeader{
			FieldName:    c.name,
			Value:        c.val,
			WordEncoding: c.encoding,
		}
		err := h.Validate()
		if c.pass {
			assert.NoError(t, err, h.FieldName, h.Value)
		} else {
			assert.Error(t, err, h.FieldName, h.Value)
		}
		assert.Equal(t, c.want, h.String())
	}
}
