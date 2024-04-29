package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	type testcase struct {
		header header.AddressField
		input  string
		want   string
		pass   bool
	}

	for _, c := range []testcase{
		{header.AddressFrom, "addr-name <addr@name.com>", "From: \"addr-name\" <addr@name.com>\r\n", true},
		{header.AddressSender, "addr@spec.com", "Sender: <addr@spec.com>\r\n", true},
		{header.AddressSender, "1@1.com, 2@2.com", "Sender: 1@1.com, 2@2.com\r\n", false},
		{header.AddressReplyTo, "alice@secret.com, Bob <bob@secret.com>", "Reply-To: <alice@secret.com>,\"Bob\" <bob@secret.com>\r\n", true},
		{header.AddressCc, "charlie@secret.com, Dmitri <dmitri@secret.com>", "Cc: <charlie@secret.com>,\"Dmitri\" <dmitri@secret.com>\r\n", true},
		{header.AddressBcc, "Eavesdrop Eve <eve@secret.com>", "Bcc: \"Eavesdrop Eve\" <eve@secret.com>\r\n", true},
		{header.AddressFrom, "invalid", "From: invalid\r\n", false},
	} {
		a := header.Address{Field: c.header, Value: c.input}
		err := a.Validate()
		if c.pass {
			assert.NoError(t, err, c.input)
		} else {
			assert.Error(t, err, c.input)
		}
		assert.Equal(t, c.want, a.String())
	}
}
