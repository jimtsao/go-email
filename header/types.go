// types.go contains internally reusable types

package header

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jimtsao/go-email/syntax"
)

type msgid string

func (m msgid) validate() error {
	// folding not permitted within actual content of msg-id
	// smtp allows only 78 octets excluding crlf
	// folding allowed before actual message id, so max content length is 78 - folding white space
	maxOctetLen := 78 - 1
	id := string(m)
	id = strings.TrimSpace(id)
	if len(id) > maxOctetLen {
		return fmt.Errorf("id must not exceed %d octets, has %d octets", maxOctetLen, len(id))
	}

	// msg-id syntax
	if !syntax.IsMsgID(id) {
		return errors.New("id invalid syntax")
	}

	return nil
}

func (m msgid) string() string {
	return strings.TrimSpace(string(m))
}
