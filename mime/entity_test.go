package mime_test

import (
	"fmt"
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/jimtsao/go-email/mime"
	"github.com/stretchr/testify/assert"
)

type aString string

func (a aString) String() string {
	return string(a)
}

func TestEntity(t *testing.T) {
	// message entity
	msg := &mime.Entity{
		Headers: []fmt.Stringer{
			header.Address{Field: header.AddressFrom, Value: "a@a.com"},
			header.Address{Field: header.AddressTo, Value: "b@b.com"},
			header.Subject("Foo Bar"),
			header.MIMEVersion{},
		},
		Body: aString("Hello World"),
	}
	want := "From: <a@a.com>\r\n" +
		"To: <b@b.com>\r\n" +
		"Subject: Foo Bar\r\n" +
		"MIME-Version: 1.0\r\n" +
		"\r\n" +
		"Hello World"
	got := msg.String()
	assert.Equal(t, want, got)
}
