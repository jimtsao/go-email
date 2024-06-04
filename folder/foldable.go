package folder

type Foldable interface {
	// length, before any optional transformations
	// i.e. if 'foo bar' can be optionally transformed into
	// =?utf-8?q?foo_bar?= it should return length 7
	Length() int
	Fold(limit int) (split string, rest Foldable, didFold bool)
}

// WordEncodable represents a string that can be conditionally encoded
// in order to utilise encoded word splitting where a foldable white space
// may not otherwise be present. Encoded word splitting has least priority.
// non us-ascii will trigger encoding.
type WordEncodable string
