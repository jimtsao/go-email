package syntax

import (
	"strings"
	"unicode"
)

func checker(s string, fn func(r rune) bool) bool {
	for _, r := range s {
		if !fn(r) {
			return false
		}
	}
	return true
}

func IsASCII(s string) bool {
	return checker(s, isASCII)
}

func isASCII(r rune) bool {
	return r <= unicode.MaxASCII
}

// IsVchar:
//
//	VCHAR = %d33-126 ; printable ascii
func IsVchar(s string) bool {
	return checker(s, isVchar)
}

func isVchar(r rune) bool {
	return '!' <= r && r <= '~'
}

// IsWSP:
//
//	WSP = SP / HTAB
func IsWSP(s string) bool {
	return checker(s, isWSP)
}

func isWSP(r rune) bool {
	return r == ' ' || r == '\t'
}

// CTL:
//
//	CTL = %d0-31 / %d127 ; control characters
func IsCTL(s string) bool {
	return checker(s, isCTL)
}

func isCTL(r rune) bool {
	return r <= 31 || r == 127
}

// IsSpecials (RFC 5322):
//
//	specials   =   "(" / ")" / "<" / ">" / "[" / "]" /
//	               ":" / ";" / "@" / "\" / "," / "." /
//	               DQUOTE
//
// characters that do not appear in atext
func IsSpecials(s string) bool {
	return checker(s, isSpecials)
}

func isSpecials(r rune) bool {
	// fast check
	if r < '"' || r > ']' {
		return false
	}

	switch r {
	case '(', ')', '<', '>', '[', ']',
		':', ';', '@', '\\', ',', '.',
		'"':
		return true
	}

	return false
}

// IsTSpecials (RFC 2045):
//
//	tspecials :=  "(" / ")" / "<" / ">" / "@" /
//	              "," / ";" / ":" / "\" / <"> /
//	              "/" / "[" / "]" / "?" / "="
func IsTSpecials(s string) bool {
	return checker(s, isTSpecials)
}

func isTSpecials(r rune) bool {
	// fast check
	if r < '"' || r > ']' {
		return false
	}

	switch r {
	case '(', ')', '<', '>', '@',
		',', ';', ':', '\\', '"',
		'/', '[', ']', '?', '=':
		return true
	}

	return false
}

// IsAtext:
//
//	atext   =   ALPHA / DIGIT /
//				"!" / "#" / "$" / "%" / "&" / "'" /
//				"*" / "+" / "-" / "/" / "=" / "?" /
//				"^" / "_" / "`" / "{" / "|" / "}" /
//				"~"
//
// printable ascii excluding specials
func IsAtext(s string) bool {
	return checker(s, isAtext)
}

func isAtext(r rune) bool {
	if isSpecials(r) {
		return false
	}

	return isVchar(r)
}

// IsDtext:
//
//	dtext = %d33-90, %d94-126
//
// printable ascii excluding "[", "]" and "\"
func IsDtext(s string) bool {
	return checker(s, isDtext)
}

func isDtext(r rune) bool {
	switch r {
	case '[', ']', '\\':
		return false
	}
	return isVchar(r)
}

// IsFtext:
//
// ftext = %d33-57 / %d59-126
//
// printable ascii excluding ":"
func IsFtext(s string) bool {
	return checker(s, isFtext)
}

func isFtext(r rune) bool {
	return r != ':' && isVchar(r)
}

// IsMIMEParamAttributeChar:
//
//	attribute-char := <any (US-ASCII) CHAR except SPACE,
//	                  CTLs, "*", "'", "%", or tspecials>
func IsMIMEParamAttributeChar(s string) bool {
	return checker(s, isMIMEParamAttributeChar)
}

func isMIMEParamAttributeChar(r rune) bool {
	return r <= '~' && r != ' ' && r != '*' &&
		r != '\'' && r != '%' &&
		!isCTL(r) && !isTSpecials(r)
}

// IsMIMEToken (RFC 2231):
//
//	token := 1*<any (US-ASCII) CHAR except SPACE, CTLs, or tspecials>
func IsMIMEToken(s string) bool {
	if s == "" {
		return false
	}
	return checker(s, isMIMEToken)
}

func isMIMEToken(r rune) bool {
	return r > ' ' && r <= '~' && !isTSpecials(r)
}

// IsQuotedString:
//
//	quoted-string   =   [CFWS] DQUOTE ((1*([FWS] qcontent) [FWS]) / FWS) DQUOTE [CFWS]
//	qcontent        =   qtext / quoted-pair
//	qtext           =   %d32 / %d33 / %d35-91 / %d93-126
//	quoted-pair     =   ("\" (VCHAR / WSP))
//
// qtext: printable ascii except \ and "
func IsQuotedString(s string) bool {
	// quoted string with at least 1 char
	if len(s) < 3 {
		return false
	}

	// check qtext or quoted-pair
	escaped := false
	for i, r := range s {
		if !isQuotedString(r, i == 0 || i == len(s)-1, escaped) {
			return false
		}

		// toggle quoted-pair mode for next character
		if !escaped && r == '\\' {
			escaped = true
			continue
		}

		escaped = false
	}

	return true
}

func isQuotedString(r rune, dquote bool, escaped bool) bool {
	// beginning or end of string
	if dquote {
		return !escaped && r == '"'
	}

	// quoted pair
	if escaped {
		return isWSP(r) || isVchar(r)
	}

	// qtext
	return r != '"' && r >= ' ' && r <= '~'
}

// IsWordEncodable:
//
//	Only printable and white space character data should be
//	encoded using this scheme. RFC 2047 Section 5.
func IsWordEncodable(s string) bool {
	for _, r := range s {
		if !isWordEncodable(r) {
			return false
		}
	}
	return true
}

func isWordEncodable(r rune) bool {
	return r == '\t' || unicode.IsPrint(r)
}

// IsDotAtomText:
//
//	dot-atom-text = 1*atext *("." 1*atext)
func IsDotAtomText(s string) bool {
	dot := true
	var i int
	var r rune
	for i, r = range s {
		if r == '.' && (i == 0 || i == len(s)-1) {
			return false
		}
		if !isDotAtomText(r, dot) {
			return false
		}
		dot = r != '.'
	}

	// must be at least 1 valid character
	return i != 0
}

func isDotAtomText(r rune, dot bool) bool {
	if r == '.' {
		return dot
	}

	return isAtext(r)
}

// IsRFC2045Token:
//
// token := 1*<any (US-ASCII) CHAR except SPACE, CTLs, or tspecials>
func IsRFC2045Token(s string) bool {
	return checker(s, isRFC2045Token)
}

func isRFC2045Token(r rune) bool {
	if isTSpecials(r) {
		return false
	}
	return '!' <= r && r <= '~'
}

// IsNoFoldLiteral:
//
//	no-fold-literal = "[" *dtext "]"
func IsNoFoldLiteral(s string) bool {
	if len(s) < 2 {
		return false
	}

	for i, r := range s {
		switch i {
		case 0:
			if r != '[' {
				return false
			}
		case len(s) - 1:
			if r != ']' {
				return false
			}
		default:
			if !isDtext(r) {
				return false
			}
		}
	}
	return true
}

// IsMsgID:
//
//	msg-id          =   [CFWS] "<" id-left "@" id-right ">" [CFWS]
//	id-left         =   dot-atom-text
//	id-right        =   dot-atom-text / no-fold-literal
//	no-fold-literal =   "[" *dtext "]"
func IsMsgID(s string) bool {
	if len(s) < 3 || s[0] != '<' || s[len(s)-1] != '>' {
		return false
	}
	s = s[1 : len(s)-1]
	left, right, found := strings.Cut(s, "@")
	if !found ||
		!IsDotAtomText(left) ||
		(!IsDotAtomText(right) && !IsNoFoldLiteral(right)) {
		return false
	}

	return true
}
