package folder_test

import (
	"fmt"
	"mime"
	"strings"
	"testing"

	"github.com/jimtsao/go-email/folder"
)

func TestWordEncodable(t *testing.T) {
	encwordq := folder.WordEncodable{"foo bar", mime.QEncoding, true, 2}
	encwordqm := folder.WordEncodable{"éoo", mime.QEncoding, true, 2}
	encwordqq := folder.WordEncodable{strings.Repeat("q", 132), mime.QEncoding, true, 2}
	encwordb := folder.WordEncodable{"foo bar", mime.BEncoding, true, 2}
	encwordbm := folder.WordEncodable{"ffé", mime.BEncoding, true, 2}
	encwordbb := folder.WordEncodable{strings.Repeat("b", 99), mime.BEncoding, true, 2}

	tcs := []testcase{
		// plain strings
		{desc: "plain string (no encode)",
			input: []interface{}{folder.WordEncodable{"foo bar", mime.QEncoding, false, 2}},
			want:  "foo bar"},
		{desc: "plain string (encode)",
			input: []interface{}{folder.WordEncodable{strings.Repeat("i", 182), mime.QEncoding, false, 2}},
			want: fmt.Sprintf("%s\r\n %[2]s\r\n %[2]s",
				"=?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=",
				"=?utf-8?q?iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii?=")},

		// quoted-printable
		{desc: "q encoded no fold",
			input: []interface{}{s(37), 1, s(20), encwordq},
			want:  fmt.Sprintf("%s\r\n %s%s", s(37), s(20), encwordq.Value())},
		{desc: "q encoded no fold (limit too little)",
			input: []interface{}{s(62), encwordq},
			want:  s(62) + encwordq.Value()},
		{desc: "q encoded no fold (multibyte char)",
			input: []interface{}{s(61), encwordqm},
			want:  s(61) + encwordqm.Value()},
		{desc: "q encoded fold (simple string)",
			input: []interface{}{s(55), encwordq},
			want:  s(55) + "=?utf-8?q?f?=\r\n =?utf-8?q?oo_bar?="},
		{desc: "q encoded fold (multibyte char)",
			input: []interface{}{s(50), encwordqm},
			want:  s(50) + "=?utf-8?q?=C3=A9?=\r\n =?utf-8?q?oo?="},
		{desc: "q encoded fold (multiple times)",
			input: []interface{}{s(50), encwordqq},
			want: fmt.Sprintf("%s%s\r\n %[3]s\r\n %[3]s",
				s(50), "=?utf-8?q?qqqqqq?=",
				"=?utf-8?q?qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq?=")},

		// base64
		{desc: "b encoded no fold(limit too little)",
			input: []interface{}{s(59), encwordb},
			want:  s(59) + encwordb.Value()},
		{desc: "b encoded fold",
			input: []interface{}{s(52), encwordb},
			want:  s(52) + "=?utf-8?b?Zm9v?=\r\n =?utf-8?b?IGJhcg==?="},
		{desc: "b encoded fold (multibyte char)",
			input: []interface{}{s(52), encwordbm},
			want:  s(52) + "=?utf-8?b?ZmY=?=\r\n =?utf-8?b?w6k=?="},
		{desc: "q encoded fold (multiple times)",
			input: []interface{}{s(50), encwordbb},
			want: fmt.Sprintf("%s%s\r\n %[3]s\r\n %[3]s",
				s(50), "=?utf-8?b?YmJi?=",
				"=?utf-8?b?YmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJiYmJi?=")},
	}

	testCases(t, "X-Header", tcs)
}
