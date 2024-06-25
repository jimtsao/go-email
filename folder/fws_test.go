package folder_test

import (
	"fmt"
	"testing"

	"github.com/jimtsao/go-email/folder"
)

func TestFWS(t *testing.T) {
	tcs := []testcase{
		{desc: "foldable whitespace preserved",
			input: []interface{}{"foo", folder.FWS(1), 1, "bar"},
			want:  "foo bar"},
		{desc: "foldable whitespace consumed",
			input: []interface{}{s(74), folder.FWS(1), "foo"},
			want:  fmt.Sprintf("%s\r\n %s", s(74), "foo")},
	}

	testCases(t, "To", tcs)
}
