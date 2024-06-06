package folder_test

import (
	"strings"
	"testing"

	"github.com/jimtsao/go-email/folder"
)

func TestMIMEParam(t *testing.T) {
	tcs := []testcase{
		{desc: "no fold - simple string",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", "foobar.txt"}},
			want:  "attachment;filename=foobar.txt"},
		{desc: "no fold - quoted string",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", `"foo bar.txt"`}},
			want:  "attachment;filename=\"foo bar.txt\""},
		{desc: "no fold - convert to quoted string",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", "foo bar.txt"}},
			want:  "attachment;filename=\"foo bar.txt\""},
		{desc: "no fold - dequote (quote not necessary)",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", "\"" + strings.Repeat("s", 37) + "\""}},
			want:  "attachment;filename=sssssssssssssssssssssssssssssssssssss"},
		{desc: "no fold - dequote (need encoding)",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", `"méow.txt"`}},
			want:  "attachment;filename*=utf-8''m%C3%A9ow.txt"},
		{desc: "no fold - multibyte char",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", "méow.txt"}},
			want:  "attachment;filename*=utf-8''m%C3%A9ow.txt"},
		{desc: "fold - simple string",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", strings.Repeat("s", 38)}},
			want:  "attachment;\r\n filename*=utf-8''ssssssssssssssssssssssssssssssssssssss"},
		{desc: "fold - multibyte char",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", strings.Repeat("s", 55) + "é"}},
			want: "attachment;\r\n filename*0*=utf-8''sssssssssssssssssssssssssssssssssssssssssssssssssssssss\r\n" +
				" filename*1*=%C3%A9"},
		{desc: "fold - quoted string",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", "\"" + strings.Repeat("s", 69) + "\""}},
			want: "attachment;\r\n filename*0*=utf-8''ssssssssssssssssssssssssssssssssssssssssssssssssssssssssss\r\n" +
				" filename*1*=sssssssssss"},
		{desc: "fold - multiple times",
			input: []interface{}{"attachment", folder.MIMEParam{"filename", strings.Repeat("s", 188)}},
			want: "attachment;\r\n filename*0*=utf-8''ssssssssssssssssssssssssssssssssssssssssssssssssssssssssss\r\n" +
				" filename*1*=sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss\r\n" +
				" filename*2*=sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss"},
	}

	testCases(t, "Content-Disposition", tcs)
}
