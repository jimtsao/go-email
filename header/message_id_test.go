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

	type testcase struct {
		input string
		want  string
		pass  bool
	}

	for _, c := range []testcase{
		{"<message@id>", "Message-ID: <message@id>\r\n", true},
		{"message-id", "Message-ID: <message-id@>\r\n", true},
		{"message@id", "Message-ID: <message@id>\r\n", true},
		{ls, "Message-ID: " + fmt.Sprintf("<%s@>\r\n", ls), true},
		{"messa<ge@i>d", "Message-ID: messa<ge@i>d\r\n", false},
		{ls + "-", "Message-ID: " + fmt.Sprintf("<%s-@>\r\n", ls), false},
		{"<in:valid@char>", "Message-ID: <in:valid@char>\r\n", false},
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
