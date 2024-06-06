package folder

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/jimtsao/go-email/syntax"
)

// MIMEParam implements RFC 2231 MIME parameter encoding and folding (continuation)
//
//	parameter              :=   regular-parameter / extended-parameter
//	regular-parameter      :=   regular-parameter-name "=" value
//	regular-parameter-name :=   attribute
//	attribute              :=   1*attribute-char
//	attribute-char         :=   <any (US-ASCII) CHAR except SPACE, CTLs,
//	                            "*", "'", "%", or tspecials>
//	value                  :=   token / quoted-string
//	token                  :=   1*<any (US-ASCII) CHAR except SPACE, CTLs, or tspecials>
//	extended-parameter     :=   (extended-initial-name "=" extended-initial-value) /
//	                            (extended-other-names  "=" extended-other-values)
//	extended-initial-name  :=   attribute [initial-section] "*"
//	extended-initial-value :=   [charset] "'" [language] "'"extended-other-values
//	extended-other-names   :=   attribute other-sections "*"
//	extended-other-values  :=   *(ext-octet / attribute-char)
//	initial-section        :=   "*0"
//	other-sections         :=   "*" ("1" / "2" / "3" / "4" / "5" /
//	                            "6" / "7" / "8" / "9") *DIGIT
//	ext-octet              :=   "%" 2(DIGIT / "A" / "B" / "C" / "D" / "E" / "F")
//	charset                :=   <registered character set name>
//	language               :=   <registered language tag [RFC-1766]>
type MIMEParam struct {
	Name string
	Val  string
}

func (m MIMEParam) Value() string {
	if rv := m.regularVal(); rv != "" {
		return rv
	}
	return m.extendedVal()
}

func (m MIMEParam) Length() int {
	return len(m.Value())
}

func (m MIMEParam) Fold(limit int) (string, Foldable, bool) {
	// no fold
	if rv := m.regularVal(); rv != "" && len(rv) <= limit {
		return rv, nil, false
	}
	ev := m.extendedVal()
	if len(ev) <= limit {
		return ev, nil, false
	}

	// extended parameter format
	sb := strings.Builder{}
	sb.WriteString(";")
	var iteration int
	remaining := m.dequote()

	// for each iteration, i.e. *0*, *1*, *2*...
	for {
		// begin new iteration
		var part string
		if iteration == 0 {
			part = fmt.Sprintf("%s%s*%d*=utf-8''", fwsToken, m.Name, iteration)
		} else {
			part = fmt.Sprintf("%s%s*%d*=", fwsToken, m.Name, iteration)

		}
		limit = maxLineLen - len(part) + len("\r\n")
		var r rune
		var i, runeLen int

		// find index where length exceeds limit
		for i = 0; i < len(remaining); i += runeLen {
			r, runeLen = utf8.DecodeRuneInString(remaining[i:])
			encRune := url.PathEscape(string(r))
			encLen := len(encRune)

			// limit exceeded
			if encLen > limit {
				// cant split even 1 single rune
				if i == 0 {
					return "", nil, false
				}

				// stop at current index
				break
			}

			limit -= encLen
		}

		// close iteration
		part += url.PathEscape(remaining[:i])
		remaining = remaining[i:]
		if iteration == 0 && remaining == "" {
			// remove *0* if only 1 total part
			part = strings.Replace(part, "*0*=", "*=", 1)
		}
		sb.WriteString(part)
		iteration++

		// nothing left to write
		if remaining == "" {
			break
		}
	}

	return sb.String(), nil, true
}

// checks if val is token, convertable to token, quoted string or convertable to quoted string
// returns empty string if none of these options possible
func (m MIMEParam) regularVal() string {
	// token
	if syntax.IsMIMEToken(m.dequote()) {
		return fmt.Sprintf(";%s=%s", m.Name, m.dequote())
	}

	// quoted string
	if syntax.IsQuotedString(m.Val) {
		return fmt.Sprintf(";%s=%s", m.Name, m.Val)
	}
	if syntax.IsQuotedString(fmt.Sprintf("\"%s\"", m.Val)) {
		return fmt.Sprintf(";%s=\"%s\"", m.Name, m.Val)
	}
	return ""
}

func (m MIMEParam) extendedVal() string {
	val := url.PathEscape(m.dequote())
	return fmt.Sprintf(";%s*=utf-8''%s", m.Name, val)
}

// dequote if param val is quoted string
func (m MIMEParam) dequote() string {
	if len(m.Val) >= 2 &&
		m.Val[0] == '"' &&
		m.Val[len(m.Val)-1] == '"' {
		return m.Val[1 : len(m.Val)-1]
	}
	return m.Val
}
