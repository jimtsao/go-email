// types.go contains internally reusable types

package header

import (
	"strings"
	"time"
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
