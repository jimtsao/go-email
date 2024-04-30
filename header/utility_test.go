package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestCanonicalHeaderKey(t *testing.T) {
	cases := map[string]string{
		"to":                   "To",
		"SENDER":               "Sender",
		"reply-TO":             "Reply-To",
		"thRee-PaRt-hEADER":    "Three-Part-Header",
		"non-[mime]-compliant": "Non-[mime]-Compliant",
		"with space":           "with space",
		"non-ascii-🐈":          "non-ascii-🐈",
	}

	for input, want := range cases {
		got := header.CanonicalHeaderKey(input)
		assert.Equal(t, want, got)
	}
}

func TestIsValidHeaderName(t *testing.T) {
	for i := 0; i <= 255; i++ {
		s := string(rune(i))
		got := header.IsValidHeaderName(s)
		want := false
		if 33 <= i && i <= 126 && i != 58 {
			want = true
		}
		assert.Equalf(t, want, got, "%q", s)
	}
}

func TestIsValidHeaderValue(t *testing.T) {
	for i := 0; i <= 255; i++ {
		s := string(rune(i))
		got := header.IsValidHeaderValue(s)
		want := false
		if (32 <= i && i <= 126) || i == 9 {
			want = true
		}
		assert.Equalf(t, want, got, "%q", s)
	}
}
