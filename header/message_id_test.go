package header_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestMessageID(t *testing.T) {
	var longStr strings.Builder
	maxContentLen := 998 - len("Message-ID: <@>")
	for i := 0; i < maxContentLen; i++ {
		longStr.WriteString("-")
	}
	ls := longStr.String()
	max := fmt.Sprintf("<%s@%s>", ls[maxContentLen/2:], ls[:maxContentLen/2])
	over := fmt.Sprintf("<%s@%s->", ls[maxContentLen/2:], ls[:maxContentLen/2])

	type testcase struct {
		input string
		want  string
		pass  bool
	}

	for _, c := range []testcase{
		{"<one@two>", "Message-ID: <one@two>\r\n", true},
		{max, fmt.Sprintf("Message-ID: %s\r\n", max), true},
		{over, fmt.Sprintf("Message-ID: %s\r\n", over), false},
	} {
		{
			m := header.MessageID(c.input)
			err := m.Validate()
			if c.pass {
				assert.NoError(t, err, c.input)
			} else {
				assert.Error(t, err, c.input)
			}
			assert.Equal(t, c.want, m.String(), c.input)
		}
	}
}
