// types.go contains internally reusable types

package header

import (
	"errors"
	"mime"
	"strings"
	"time"

	"github.com/jimtsao/go-email/syntax"
)

const maxLineLen = 998

const time_RFC5322 = "Mon, 2 Jan 2006 15:04:05 -0700"

type datetime time.Time

func (dt datetime) String() string {
	return time.Time(dt).Format(time_RFC5322)
}

type msgid string

func (m msgid) String() string {
	return strings.TrimSpace(string(m))
}

type unstructured string

func (u unstructured) validate(encode bool) error {
	v := string(u)
	if encode {
		if strings.Contains(v, ":") || !syntax.IsWordEncodable(v) {
			return errors.New("must contain only printable or white space characters and no colon")
		}
	} else if !syntax.IsFtext(v) {
		return errors.New("invalid syntax")
	}
	return nil
}

func (u unstructured) string(encode bool) string {
	if encode {
		return mime.QEncoding.Encode("utf-8", string(u))
	}
	return string(u)
}
