package base64_test

import (
	"strings"
	"testing"

	"github.com/jimtsao/go-email/base64"
	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	sb := &strings.Builder{}
	enc := base64.NewEncoder(sb)
	for i := 0; i < 2; i++ {
		_, err := enc.Write([]byte(strings.Repeat("foo", 25)))
		assert.NoError(t, err)
	}

	want := "Zm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9v\n" +
		"Zm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9v\n" +
		"Zm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9vZm9v"
	assert.Equal(t, want, sb.String())
}
