package header_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestAddressValidate(t *testing.T) {
	// sender no more than 1 address
	addr := header.Address{
		Field: header.AddressSender,
		Value: "a@b.com, c@d.com",
	}
	err := addr.Validate()
	assert.Error(t, err)

	// local exceeds 64 octets
	addr.Value = fmt.Sprintf("%s@b.com", strings.Repeat("a", 65))
	err = addr.Validate()
	assert.Error(t, err)

	// domain exceeds 255
	addr.Value = fmt.Sprintf("a@%s.com", strings.Repeat("b", 256))
	err = addr.Validate()
	assert.Error(t, err)
}

func TestAddressNoFold(t *testing.T) {
	type testcase struct {
		header header.AddressField
		input  string
		want   string
	}

	for _, c := range []testcase{
		{header.AddressFrom, "addr-name <addr@name.com>", "\"addr-name\" <addr@name.com>"},
		{header.AddressSender, "addr@spec.com", "<addr@spec.com>"},
		{header.AddressReplyTo, "alice@secret.com, Bob <bob@secret.com>", "<alice@secret.com>,\"Bob\" <bob@secret.com>"},
		{header.AddressCc, "charlie@secret.com, Dmitri <dmitri@secret.com>", "<charlie@secret.com>,\"Dmitri\" <dmitri@secret.com>"},
		{header.AddressBcc, "Eavesdrop Eve <eve@secret.com>", "\"Eavesdrop Eve\" <eve@secret.com>"},
	} {
		a := header.Address{Field: c.header, Value: c.input}
		err := a.Validate()
		assert.NoError(t, err, c.input)
		assert.Equal(t, fmt.Sprintf("%s: %s\r\n", c.header, c.want), a.String())
	}
}

func TestAddressFold(t *testing.T) {
	for _, c := range []struct {
		desc  string
		input string
		want  string
	}{
		{desc: "fold after header",
			input: "<a@iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii.com>",
			want:  "\r\n <a@iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii.com>"},

		// angle-addr: [1]angle-addr[1]
		{desc: "fold domain",
			input: "<iiiiiiiiiiiiiii@iiiiiiiiiiiiiii>,<iiiiiiiiiiiiiii@iiiiiiiiiiiiiii>,<iiiiiiiiiiiiiii@iiiiiiiiiiiiiii>",
			want:  "<iiiiiiiiiiiiiii@iiiiiiiiiiiiiii>,<iiiiiiiiiiiiiii@iiiiiiiiiiiiiii>,\r\n <iiiiiiiiiiiiiii@iiiiiiiiiiiiiii>"},

		// quoted-string: [1]quoted[3][space]string[2][space]angle-addr[1]
		{desc: "fold quoted-string",
			input: "\"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii\" <a@b>, \"sssssssssssssssssssssssssssssssssssss sssssssssssssssssssssssssssssssssssss\" <a@b>",
			want:  "\"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii\" <a@b>,\r\n \"sssssssssssssssssssssssssssssssssssss sssssssssssssssssssssssssssssssssssss\"\r\n <a@b>"},

		// encoded-word: encoded-word[2][space]angle-addr[1]
		{desc: "encoded word",
			input: "éiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii <a@b>",
			want:  "=?utf-8?q?=C3=A9iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=\r\n =?utf-8?q?iiiiii?= <a@b>"},
	} {
		a := header.Address{Field: header.AddressTo, Value: c.input}
		err := a.Validate()
		assert.NoError(t, err, c.desc)
		assert.Equal(t, fmt.Sprintf("To: %s\r\n", c.want), a.String(), c.desc)
	}
}
