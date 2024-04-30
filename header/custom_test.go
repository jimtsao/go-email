package header_test

import (
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestCustomHeader(t *testing.T) {
	h := header.CustomHeader{
		FieldName: "my-header",
		Value:     "!hola	amigo~",
	}
	err := h.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "My-Header: !hola	amigo~\r\n", h.String())
}
