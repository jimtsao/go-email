package header

// Header field
type Header interface {
	// Name of header field, MUST be in canonical form
	// eg. 'Reply-To' rather than 'reply-to'
	Name() string
	// Validate against relevant RFCs
	Validate() error
	// String outputs in format complaint with relevant RFCs.
	// The general format is:
	//
	//	Header fields are lines beginning with a field name, followed by a
	//	colon (":"), followed by a field body, and terminated by CRLF
	//
	// In case of validation error, String() should default to user
	// supplied value without any additional formatting
	String() string
}
