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

func TestFoldingEncodedWords(t *testing.T) {
	encword19q := "=?utf-8?q?foo_bar?="
	encword20qm := "=?utf-8?q?=C3=A9oo?="
	encword20b := "=?utf-8?b?Zm9vYmFy?="
	encword20bm := "=?utf-8?b?ZmbDqQ==?="
	tcs := []testcase{
		{desc: "multiple encoded word no split (priority)",
			input: []interface{}{s(50), 1, encword19q + " " + encword19q},
			want:  fmt.Sprintf("%s\r\n %s %s", s(50), encword19q, encword19q)},
		{desc: "multiple encoded word split",
			input: []interface{}{s(50), encword19q + " " + encword19q},
			want:  fmt.Sprintf("%s%s\r\n %s", s(50), encword19q, encword19q)},
		{desc: "q encoded no split (priority)",
			input: []interface{}{s(37), 1, s(20), encword19q},
			want:  fmt.Sprintf("%s\r\n %s%s", s(37), s(20), encword19q)},
		{desc: "q encoded no split (limit too little)",
			input: []interface{}{s(62), encword19q},
			want:  s(62) + encword19q},
		{desc: "q encoded no split (multibyte char)",
			input: []interface{}{s(61), encword20qm},
			want:  s(61) + encword20qm},
		{desc: "q encoded split",
			input: []interface{}{s(61), encword19q},
			want:  s(61) + "=?utf-8?q?f?=\r\n =?utf-8?q?oo_bar?="},
		{desc: "q encoded split (multibyte char)",
			input: []interface{}{s(56), encword20qm},
			want:  s(56) + "=?utf-8?q?=C3=A9?=\r\n =?utf-8?q?oo?="},
		{desc: "b encoded no split (limit too little)",
			input: []interface{}{s(59), encword20b},
			want:  s(59) + encword20b},
		{desc: "b encoded split",
			input: []interface{}{s(58), encword20b},
			want:  s(58) + "=?utf-8?b?Zm9v?=\r\n =?utf-8?b?YmFy?="},
		{desc: "b encoded split (multibyte char)",
			input: []interface{}{s(58), encword20bm},
			want:  s(58) + "=?utf-8?b?ZmY=?=\r\n =?utf-8?b?w6k=?="},
	}
	testCases(t, "To", tcs)
}

func TestWordEncodable(t *testing.T) {
	sb := &strings.Builder{}
	s := strings.Repeat("i", 182)
	f := folder.New(sb)
	f.Write("Reply-To: ")
	f.Write(folder.WordEncodable(s))
	f.Close()

	assert.NoError(t, f.Err)
	want := "Reply-To: =?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=\r\n =?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=\r\n =?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=\r\n"
	assert.Equal(t, want, sb.String())
}
