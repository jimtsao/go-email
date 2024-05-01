package syntax

import (
	"strings"
)

func checker(s string, fn func(b byte) bool) bool {
	for i := 0; i < len(s); i++ {
		if !fn(s[i]) {
			return false
		}
	}
	return true
}

// IsVchar definition:
//
//	VCHAR = %d33-126
//
// printable ascii
func IsVchar(s string) bool {
	return checker(s, isVchar)
}

func isVchar(b byte) bool {
	return '!' <= b && b <= '~'
}

// CTL   =   %d0-31 / %d127
func IsCTL(s string) bool {
	return checker(s, isCTL)
}

func isCTL(b byte) bool {
	return b <= 31 || b == 127
}

// IsSpecials RFC 5322 definition:
//
//	specials   =   "(" / ")" / "<" / ">" / "[" / "]" /
//	               ":" / ";" / "@" / "\" / "," / "." /
//	               DQUOTE
//
// characters that do not appear in atext
func IsSpecials(s string) bool {
	return checker(s, isSpecials)
}

func isSpecials(b byte) bool {
	// fast check
	if b < '"' || b > ']' {
		return false
	}

	switch b {
	case '(', ')', '<', '>', '[', ']',
		':', ';', '@', '\\', ',', '.',
		'"':
		return true
	}

	return false
}

// IsTSpecials RFC2045 definition:
//
//	tspecials :=  "(" / ")" / "<" / ">" / "@" /
//	              "," / ";" / ":" / "\" / <"> /
//	              "/" / "[" / "]" / "?" / "="
func IsTSpecials(s string) bool {
	return checker(s, isTSpecials)
}

func isTSpecials(b byte) bool {
	// fast check
	if b < '"' || b > ']' {
		return false
	}

	switch b {
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
//	ALPHA   =   %d65-90 / %d97-122
//	DIGIT   =   %d48-57
//
// atext: VCHAR excluding specials
//
// ALPHA: A-Z a-z
//
// DIGIT: 0-9
func IsAtext(s string) bool {
	return checker(s, isAtext)
}

func isAtext(b byte) bool {
	if isSpecials(b) {
		return false
	}

	return isVchar(b)
}

// IsDtext:
//
//	dtext = %d33-90, %d94-126
//
// printable ascii excluding "[", "]", or "\"
func IsDtext(s string) bool {
	return checker(s, isDtext)
}

func isDtext(b byte) bool {
	switch b {
	case '[', ']', '\\':
		return false
	}
	return isVchar(b)
}

// IsDotAtomText:
//
//	dot-atom-text = 1*atext *("." 1*atext)
func IsDotAtomText(s string) bool {
	dot := true
	var i int
	for i = 0; i < len(s); i++ {
		if s[i] == '.' && (i == 0 || i == len(s)-1) {
			return false
		}
		if !isDotAtomText(s[i], dot) {
			return false
		}
		dot = s[i] != '.'
	}

	// must be at least 1 valid character
	return i != 0
}

func isDotAtomText(b byte, dot bool) bool {
	if b == '.' {
		return dot
	}

	return isAtext(b)
}

// IsRFC2045Token:
//
// token := 1*<any (US-ASCII) CHAR except SPACE, CTLs, or tspecials>
func IsRFC2045Token(s string) bool {
	return checker(s, isRFC2045Token)
}

func isRFC2045Token(b byte) bool {
	if isTSpecials(b) {
		return false
	}
	return '!' <= b && b <= '~'
}

// IsNoFoldLiteral:
//
//	no-fold-literal = "[" *dtext "]"
func IsNoFoldLiteral(s string) bool {
	if len(s) < 2 {
		return false
	}

	for i := 0; i < len(s); i++ {
		switch i {
		case 0:
			if s[i] != '[' {
				return false
			}
		case len(s) - 1:
			if s[i] != ']' {
				return false
			}
		default:
			if !isDtext(s[i]) {
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
