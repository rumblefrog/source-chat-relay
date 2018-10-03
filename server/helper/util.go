package helper

import (
	"strings"
	"unicode"
)

func StripSymbol(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSymbol(r) {
			return rune(-1)
		}

		return r
	}, s)
}
