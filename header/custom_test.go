package header_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestCustomHeaderValidation(t *testing.T) {
	// header name invalid syntax
	h := header.CustomHeader{
		FieldName: "X Header",
	}
	err := h.Validate()
	assert.Error(t, err)

	// value contains colon
	h.FieldName = "X-Experimental"
	h.Value = "a:b"
	err = h.Validate()
	assert.Error(t, err)

	// not word encodable
	h.WordEncodable = true
	h.Value = "hello\nworld"
	err = h.Validate()
	assert.Error(t, err)
}

func TestCustomHeader(t *testing.T) {
	for _, c := range []struct {
		desc      string
		name      string
		encodable bool
		input     string
		want      string
	}{
		{desc: "simple", name: "X-Foo", input: "bar", want: "X-Foo: bar\r\n"},
		{desc: "no canonical header", name: "x-fOO", input: "bar", want: "x-fOO: bar\r\n"},
		{desc: "no encoding", name: "X-Foo", input: "tést", want: "X-Foo: tést\r\n"},
		{desc: "yes encoding", name: "X-Foo", encodable: true, input: "tést", want: "X-Foo: =?utf-8?q?t=C3=A9st?=\r\n"},
		{desc: "folding", name: "Subject", input: strings.Repeat("i", 75), want: fmt.Sprintf("Subject:\r\n %s\r\n", strings.Repeat("i", 75))},
		{desc: "encodable word", name: "Subject", encodable: true,
			input: strings.Repeat("i", 126),
			want:  "Subject:\r\n =?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=\r\n =?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=\r\n"},
	} {
		h := header.CustomHeader{
			FieldName:     c.name,
			Value:         c.input,
			WordEncodable: c.encodable,
		}
		err := h.Validate()
		assert.NoError(t, err, h.FieldName, h.Value, c.desc)
		assert.Equal(t, c.want, h.String(), c.desc)
	}
}
