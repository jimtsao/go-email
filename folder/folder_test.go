package folder_test

import (
	"fmt"
	"mime"
	"strings"
	"testing"

	"github.com/jimtsao/go-email/folder"
	"github.com/stretchr/testify/assert"
)

type testcase struct {
	input []interface{}
	want  string
	desc  string
}

func s(n int) string {
	return strings.Repeat("i", n)
}

func testCases(t *testing.T, header string, tcs []testcase) {
	for _, tc := range tcs {
		sb := &strings.Builder{}
		f := folder.New(sb)
		f.Write(header + ": ")
		f.Write(tc.input...)
		f.Close()
		assert.NoError(t, f.Err, tc.desc)
		want := fmt.Sprintf("%s: %s\r\n", header, tc.want)
		assert.Equal(t, want, sb.String(), tc.desc)
	}
}

func TestNoFolding(t *testing.T) {
	tcs := []testcase{
		{desc: "empty values", want: ""},
		{desc: "within line limit", input: []interface{}{"foo", 1, "bar"}, want: "foobar"},
		{desc: "no foldable locations", input: []interface{}{s(74), "world"}, want: s(74) + "world"},
	}
	testCases(t, "To", tcs)
}

func TestFolding(t *testing.T) {
	// fold after input
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write("Reply-To:")
	f.Write(folder.FWS(1), s(80))
	f.Close()
	assert.NoError(t, f.Err, "fold after input")
	want := fmt.Sprintf("Reply-To:\r\n %s\r\n", s(80))
	assert.Equal(t, want, sb.String(), "fold after input")

	// test cases
	tcs := []testcase{
		{desc: "0 priority ignored",
			input: []interface{}{s(40), 0, s(40), folder.WordEncodable{
				Decoded:      "iii",
				Enc:          mime.QEncoding,
				MustEncode:   true,
				FoldPriority: 0,
			}},
			want: strings.Repeat("i", 80) + "=?utf-8?q?iii?="},
		{desc: "priority respected",
			input: []interface{}{s(37), 2, s(37), 1, "foo"},
			want:  fmt.Sprintf("%s\r\n %s", s(74), "foo")},
		{desc: "unfoldable does not stop subsequent folding",
			input: []interface{}{s(80), 1, s(80)},
			want:  s(80) + "\r\n " + s(80)},
		{desc: "multiple folding",
			input: []interface{}{s(37), 1, s(37), 2, s(80)},
			want:  fmt.Sprintf("%s\r\n %s\r\n %s", s(37), s(37), s(80))},
		{desc: "whitespaces only before fold",
			input: []interface{}{"foo", 1, s(80)},
			want:  fmt.Sprintf("foo\r\n %s", s(80))},
		{desc: "whitespaces only after fold",
			input: []interface{}{s(37), 1, s(37), folder.FWS(1), 1, "\t", 1, " \t"},
			want:  fmt.Sprintf("%s\r\n %s%s", s(37), s(37), " \t \t")},
	}
	testCases(t, "To", tcs)
}

type EmptyFolder struct {
	val        string
	priority   int
	emptyLeft  bool
	emptyRight bool
}

func (e EmptyFolder) Value() string {
	return e.val
}

func (e EmptyFolder) Fold(limit int) string {
	if e.emptyLeft && e.emptyRight {
		return "\r\n " + e.val + "\r\n "
	} else if e.emptyLeft {
		return "\r\n " + e.val
	} else if e.emptyRight {
		return e.val + "\r\n "
	}

	return e.val + "\r\n " + e.val
}

func (e EmptyFolder) Priority() int {
	return e.priority
}

func TestWhitespaceCheck(t *testing.T) {
	// WSP only before fold
	desc := "WSP only before fold"
	wspLeft := EmptyFolder{
		val:       s(80),
		priority:  1,
		emptyLeft: true,
	}
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write("\t", " ", "\t ", wspLeft, "foo")
	f.Close()
	assert.NoError(t, f.Err, desc)
	want := "\t \t " + wspLeft.val + "foo\r\n"
	assert.Equal(t, want, sb.String(), desc)

	// test cases
	wspRight := EmptyFolder{
		val:        s(80),
		priority:   1,
		emptyRight: true,
	}
	wspLeftRight := EmptyFolder{
		val:        s(80),
		priority:   1,
		emptyLeft:  true,
		emptyRight: true,
	}

	tcs := []testcase{
		{desc: "WSP only after fold",
			input: []interface{}{"foo", wspRight, "\t", " ", "\t "},
			want:  "foo" + wspRight.val + "\t \t "},
		{desc: "WSP only before and after fold",
			input: []interface{}{s(80), 1, "\t", " ", "\t ", wspLeftRight, "\t", " ", "\t "},
			want:  s(80) + "\r\n \t \t " + wspLeftRight.val + "\t \t "},
	}
	testCases(t, "To", tcs)
}
