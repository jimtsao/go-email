package email

// Header represents an entity header field
type Header interface {
	// Name of header field, MUST be in canonical form
	// eg. 'Reply-To' rather than 'reply-to'
	Name() string
	// Validate against RFC 5322 syntax
	Validate() error
	// String outputs RFC 5322 compliant format:
	//
	//	Header fields are lines beginning with a field name, followed by a
	//	colon (":"), followed by a field body, and terminated by CRLF
	//
	// In case of validation error, String() should default to user
	// supplied value without any additional formatting
	String() string
}
