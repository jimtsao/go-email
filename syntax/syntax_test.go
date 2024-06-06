package syntax_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/jimtsao/go-email/syntax"
	"github.com/stretchr/testify/assert"
)

const (
	alpha         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit         = "0123456789"
	atext         = alpha + digit + atextSpecials
	atextSpecials = "!#$%&'*+-/=?^_`{|}~"
	specials      = "()<>[]:;@\\,.\""
	tSpecials     = "()<>@,;:\\\"/[]?="
	tokenSpecials = "!#$%&'*+-^_`{|}~."
	token         = alpha + digit + tokenSpecials
	vchar         = atext + specials
	dtext         = atext + "()<>:;@,.\"" // vchar minus []\
)

func check(t *testing.T, fn func(string) bool, allowed ...string) {
	// build allowed charset
	sb := strings.Builder{}
	for _, a := range allowed {
		min, max, _ := strings.Cut(a, "-")
		if mins, err := strconv.Atoi(min); err == nil {
			// is range
			if maxs, err := strconv.Atoi(max); err == nil {
				for i := mins; i <= maxs; i++ {
					sb.WriteRune(rune(i))
				}
			}

			// is single number
			sb.WriteRune(rune(mins))
		} else {
			// is string
			sb.WriteString(a)
		}
	}

	// test some unicode code point ranges
	// 0-127: control and basic latin (ascii)
	// 128-255: control and latin-1 supplement
	// 256-383: latin extended-a
	for i := 0; i < 383; i++ {
		r := rune(i)
		want := strings.ContainsRune(sb.String(), r)
		got := fn(string(r))
		assert.Equalf(t, want, got, "ascii (%d): %q", i, r)
	}
}

func TestIsASCII(t *testing.T) {
	check(t, syntax.IsASCII, "0-127")
}

func TestIsVchar(t *testing.T) {
	check(t, syntax.IsVchar, "33-126")
}

func TestIsWSP(t *testing.T) {
	check(t, syntax.IsWSP, " ", "\t")
}

func TestIsCTL(t *testing.T) {
	check(t, syntax.IsCTL, "0-31", "127")
}

func TestIsSpecials(t *testing.T) {
	check(t, syntax.IsSpecials, specials)
}

func TestIsTSpecials(t *testing.T) {
	check(t, syntax.IsTSpecials, tSpecials)
}

func TestIsAtext(t *testing.T) {
	check(t, syntax.IsAtext, "48-57", "65-90", "97-122", atextSpecials)
}

func TestIsDtext(t *testing.T) {
	check(t, syntax.IsDtext, "33-90", "94-126")
}

func TestIsFtext(t *testing.T) {
	check(t, syntax.IsFtext, "33-57", "59-126")
}

func TestIsMIMEParamAttributeChar(t *testing.T) {
	check(t, syntax.IsMIMEParamAttributeChar,
		"33", "35", "36", "38", "43", "45",
		"46", "48-57", "65-90", "94-126")
}

func TestIsQuotedString(t *testing.T) {
	for input, want := range map[string]bool{
		`""`:                true,
		`"foo"`:             true,
		`"foo\\"`:           true,
		`"\ \` + "\t" + `"`: true,
		``:                  false,
		`"foo`:              false,
		`foo"`:              false,
		`"foo\"`:            false,
		`"foo bar"`:         false,
	} {
		got := syntax.IsQuotedString(input)
		assert.Equalf(t, want, got, "%q", input)
	}
}

func TestIsWordEncodable(t *testing.T) {
	for input, want := range map[string]bool{
		vchar:    true,
		"\t ":    true,
		"\r":     false,
		"\n":     false,
		"\u009C": false,
	} {
		got := syntax.IsWordEncodable(input)
		assert.Equalf(t, want, got, "%q", input)
	}
}

func TestIsDotAtomText(t *testing.T) {
	for s, pass := range map[string]bool{
		atext:        true,
		"2024.01.24": true,
		"127.0.0.1":  true,
		"d.ot":       true,
		".dot":       false,
		"dot.":       false,
		"dot..dot":   false,

		dtext:    false,
		specials: false,
		vchar:    false,
	} {
		got := syntax.IsDotAtomText(s)
		if pass {
			assert.Truef(t, got, "%+q", s)
		} else {
			assert.Falsef(t, got, "%+q", s)
		}
	}
}

func TestIsRFC2045Token(t *testing.T) {
	check(t, syntax.IsRFC2045Token, token)
}

func TestIsNoFoldLiteral(t *testing.T) {
	for s, want := range map[string]bool{
		"[]":                 true,
		"[" + dtext + "]":    true,
		"a":                  false,
		"[a":                 false,
		"a]":                 false,
		"[" + specials + "]": false,
		"":                   false,
		"[":                  false,
		"]":                  false,
	} {
		got := syntax.IsNoFoldLiteral(s)
		assert.Equalf(t, want, got, "%q", s)
	}
}

func TestIsMessageID(t *testing.T) {
	for s, want := range map[string]bool{
		// valid
		"<left@right>":           true,
		"<2024.12@127.0.0.1>":    true,
		"<left@[]>":              true,
		"<left@[" + dtext + "]>": true,
		"<test@[]>":              true,

		// invalid
		"":                         false,
		"one":                      false,
		"<>":                       false,
		"@":                        false,
		"<@>":                      false,
		"<one>":                    false,
		"<left@>":                  false,
		"<@right>":                 false,
		"<.left@right>":            false,
		"<left.@right>":            false,
		"le<ft@ri>ght":             false,
		"<" + specials + "@right>": false,
		"<left@" + specials + ">":  false,
	} {
		got := syntax.IsMsgID(s)
		assert.Equalf(t, want, got, "%q", s)
	}
}
