package header

// CanonicalHeaderKey returns a canonical form of the key
// whereby the first letter of each word is capitalised
//
// eg, "reply-to" becomes "Reply-To"
//
// if key contains any invalid characters, key is returned unchanged
//
// note this implementation differs to net/textproto CanonicalMIMEHeaderKey
// which follows RFC 7230
func CanonicalHeaderKey(key string) string {
	const toLower = 'a' - 'A'

	copy := []byte(key)
	upper := true
	for i, c := range copy {
		// return if invalid character
		if !isValidHeaderNameByte(c) {
			return key
		}

		// convert case if needed
		if upper && 'a' <= c && c <= 'z' {
			// to uppercase
			copy[i] = c - toLower
		} else if !upper && 'A' <= c && c <= 'Z' {
			// to lowercase
			copy[i] = c + toLower
		}

		// convert next letter to upper if after hyphen
		upper = c == '-'
	}

	return string(copy)
}

// IsValidHeaderName reports whether a header
// field name contains only valid characters
//
//	A field name MUST be composed of printable US-ASCII characters
//	(i.e. characters that have values between 33 and 126, inclusive),
//	except colon (58)
//
// Errata 5918:
//
//	...field names should be limited to 77 characters – the field name and
//	a trailing : – after which the field body can start after FWS on the next line.
func IsValidHeaderName(s string) bool {
	if len(s) > 77 {
		return false
	}

	for _, c := range []byte(s) {
		if !isValidHeaderNameByte(c) {
			return false
		}
	}
	return true
}

// IsValidHeaderValue reports whether a header
// field body contains only valid characters
//
//	A field body may be composed of printable US-ASCII characters as well
//	as the space (SP, ASCII value 32) and horizontal tab (HTAB, ASCII value 9)
//	characters (together known as the white space characters, WSP)
func IsValidHeaderValue(s string) bool {
	for _, c := range []byte(s) {
		if !isValidHeaderValueByte(c) {
			return false
		}
	}
	return true
}

// valid range from !(33) to ~(126) except :(58)
func isValidHeaderNameByte(c byte) bool {
	return '!' <= c && c <= '~' && c != ':'
}

// valid range from !(33) to ~(126), plus space(32) and htab(9)
func isValidHeaderValueByte(c byte) bool {
	return c == 9 || (' ' <= c && c <= '~')
}
