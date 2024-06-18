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
	Attribute    string
	Val          string
	FoldPriority int
}

func (m MIMEParam) Value() string {
	if m.requiresExtendedFormat() {
		return m.extendedVal()
	}
	return m.regularVal()
}

func (m MIMEParam) Priority() int {
	return m.FoldPriority
}

// Fold performs both encoding and folding, though it probably should not
// according to Single Responsibility Principle. However since both encoded
// and folded form have the same extended parameter syntax, we deduplicate
// work by only accepting unencoded values and deciding whether to encode/fold
// at the same time
func (m MIMEParam) Fold(limit int) string {
	// no fold - empty val
	if m.Val == "" || m.Val == `""` {
		return ""
	}

	// no fold - within limit or unattainable limit
	minLen := len(fmt.Sprintf("%s*0*=utf-8''", m.Attribute))
	if len(m.Value()) <= limit || limit <= minLen {
		return ""
	}

	// requires folding
	sb := strings.Builder{}
	var iteration int
	remaining := m.dequote()

	// for each iteration, i.e. *0*, *1*, *2*...
	for {
		// begin new iteration
		var part string
		if iteration == 0 {
			part = fmt.Sprintf("%s*%d*=utf-8''", m.Attribute, iteration)
			limit -= len(part)
		} else {
			part = fmt.Sprintf("%s%s*%d*=", fwsToken, m.Attribute, iteration)
			limit -= len(part) - len("\r\n")
		}

		// find index where length exceeds limit
		var r rune
		var i, runeLen int
		for i = 0; i < len(remaining); i += runeLen {
			r, runeLen = utf8.DecodeRuneInString(remaining[i:])
			encRune := url.PathEscape(string(r))
			encLen := len(encRune)

			// limit exceeded
			if encLen > limit {
				// cant split even 1 single rune
				if i == 0 {
					return ""
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

		// reset limit
		limit = maxLineLen
	}

	return sb.String()
}

// regular param value := token / quoted-string
func (m MIMEParam) requiresExtendedFormat() bool {
	// technically token and quoted string requires at least 1 char
	// but we allow garbage in garbage out, it is more correct
	// to display it as garbage regular param than garbage extended param
	return m.Val != "" && m.Val != `""` &&
		!syntax.IsMIMEToken(m.Val) && !syntax.IsQuotedString(m.Val) &&
		!syntax.IsQuotedString(`"`+m.Val+`"`)
}

func (m MIMEParam) regularVal() string {
	// value valid only if quoted string
	if !syntax.IsMIMEToken(m.Val) &&
		!syntax.IsQuotedString(m.Val) &&
		syntax.IsQuotedString(`"`+m.Val+`"`) {
		return fmt.Sprintf("%s=\"%s\"", m.Attribute, m.Val)
	}

	// return value unchanged
	return fmt.Sprintf("%s=%s", m.Attribute, m.Val)
}

func (m MIMEParam) extendedVal() string {
	val := m.dequote()
	val = url.PathEscape(val)
	return fmt.Sprintf("%s*=utf-8''%s", m.Attribute, val)
}

func (m MIMEParam) dequote() string {
	if len(m.Val) >= 2 && m.Val[0] == '"' && m.Val[len(m.Val)-1] == '"' {
		return m.Val[1 : len(m.Val)-1]
	}
	return m.Val
}
