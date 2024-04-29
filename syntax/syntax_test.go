package syntax_test

import (
	"strings"
	"testing"

	"github.com/jimtsao/go-email/syntax"
	"github.com/stretchr/testify/assert"
)

const (
	alpha    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit    = "0123456789"
	atext    = alpha + digit + "!#$%&'*+-/=?^_`{|}~"
	specials = "()<>[]:;@\\,.\""
	vchar    = atext + specials
	dtext    = atext + "()<>:;@,.\"" // vchar minus []\
)

func check(t *testing.T, fn func(string) bool, allowed string) {
	for i := 0; i < 255; i++ {
		s := string([]byte{uint8(i)})
		want := strings.Contains(allowed, s)
		got := fn(s)
		assert.Equalf(t, want, got, "%q", s)
	}
}

func TestIsVchar(t *testing.T) {
	check(t, syntax.IsVchar, vchar)
}

func TestIsSpecials(t *testing.T) {
	check(t, syntax.IsSpecials, specials)
}

func TestIsAtext(t *testing.T) {
	check(t, syntax.IsAtext, atext)
}

func TestIsDtext(t *testing.T) {
	check(t, syntax.IsDtext, dtext)
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
