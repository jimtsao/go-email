package folder_test

import (
	"fmt"
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
		{desc: "foldable whitespace preserved", input: []interface{}{"foo", 1, " ", 1, "bar"}, want: "foo bar"},
		{desc: "no foldable locations", input: []interface{}{s(74), "world"}, want: s(74) + "world"},
	}
	testCases(t, "To", tcs)
}

func TestFolding(t *testing.T) {
	// fold after input
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write("Reply-To:")
	f.Write(1, " ", s(80))
	f.Close()
	assert.NoError(t, f.Err, "fold after input")
	want := fmt.Sprintf("Reply-To:\r\n %s\r\n", s(80))
	assert.Equal(t, want, sb.String(), "fold after input")

	// test cases
	tcs := []testcase{
		{desc: "priority respected",
			input: []interface{}{s(37), 2, s(37), 1, "foo"},
			want:  fmt.Sprintf("%s\r\n %s", s(74), "foo")},
		{desc: "foldable whitespace consumed",
			input: []interface{}{s(74), 2, 1, " ", "foo"},
			want:  fmt.Sprintf("%s\r\n %s", s(74), "foo")},
		{desc: "multiple folding",
			input: []interface{}{s(37), 1, s(37), 2, s(80)},
			want:  fmt.Sprintf("%s\r\n %s\r\n %s", s(37), s(37), s(80))},
		{desc: "whitespaces only after fold",
			input: []interface{}{s(37), 1, s(37), 1, " ", 1, "\t", 1, " \t"},
			want:  fmt.Sprintf("%s\r\n %s%s", s(37), s(37), " \t \t")},
		{desc: "whitespaces only before fold",
			input: []interface{}{"foo", 1, 2, s(80)},
			want:  fmt.Sprintf("foo\r\n %s", s(80))},
	}
	testCases(t, "To", tcs)
}
