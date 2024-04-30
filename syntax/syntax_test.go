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

	// test against this charset
	for i := 0; i < 255; i++ {
		r := rune(i)
		want := strings.ContainsRune(sb.String(), r)
		got := fn(string(r))
		assert.Equalf(t, want, got, "ascii (%d): %q", i, r)
	}
}

func TestIsVchar(t *testing.T) {
	check(t, syntax.IsVchar, "33-126")
}

func TestIsSpecials(t *testing.T) {
	check(t, syntax.IsSpecials, specials)
}

func TestIsAtext(t *testing.T) {
	check(t, syntax.IsAtext, "48-57", "65-90", "97-122", atextSpecials)
}

func TestIsDtext(t *testing.T) {
	check(t, syntax.IsDtext, "33-90", "94-126")
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
		"<20241231.123=@example.net>":  true,
		"<123@localhost.com>":          true,
		"<2024.01.24@[" + dtext + "]>": true,
		"<test@[]>":                    true,
		"<@>":                          true,
		"a":                            false,
		"test@test>":                   false,
		"<test@test":                   false,
		"<" + specials + "@test>":      false,
		"<test@" + specials + ">":      false,
	} {
		got := syntax.IsMessageID(s)
		assert.Equalf(t, want, got, "%q", s)
	}
}
