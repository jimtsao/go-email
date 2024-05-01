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
	// trim whitespace and inserts @ and < > if needed
	s := string(m)
	s = strings.TrimSpace(s)
	if !strings.Contains(s, "@") {
		s = s + "@"
	}
	if !strings.Contains(s, "<") && !strings.Contains(s, ">") {
		return "<" + s + ">"
	}
	return s
}
