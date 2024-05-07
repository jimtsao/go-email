package header

// Subject represents the 'Date' header field
//
// Syntax:
//
//	subject         =   "Subject:" unstructured CRLF
//	unstructured    =   (*([FWS] VCHAR) *WSP)
type Subject string

func (s Subject) Name() string {
	return "Subject"
}

// Validate
//
// since any word encoding will produce printable ascii
// and satisfy 'unstructured' definition, we check that
// it can be word encoded instead
func (s Subject) Validate() error {
	return CustomHeader{
		FieldName:    s.Name(),
		Value:        string(s),
		WordEncoding: true,
	}.Validate()
}

func (s Subject) String() string {
	return CustomHeader{
		FieldName:    s.Name(),
		Value:        string(s),
		WordEncoding: true,
	}.String()
}
