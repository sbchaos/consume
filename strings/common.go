package strings

import (
	"unicode"

	"github.com/sbchaos/consume/base"
	"github.com/sbchaos/consume/base/higher"
	"github.com/sbchaos/consume/stream"
)

var Quotes = []struct {
	Start rune
	End   rune
}{
	{Start: '\'', End: '\''},
	{Start: '"', End: '"'},
	{Start: '`', End: '`'},
}

func AlphaNumeric() base.Parser[rune, string] {
	return CustomString(func(a rune) bool {
		return unicode.IsLetter(a) || unicode.IsDigit(a)
	})
}

func StringWithOptionalQuotes() base.Parser[rune, string] {
	return higher.Choice(
		QuotedString(0, Quotes...),
		AlphaNumeric(),
	)
}

func CustomString(fn func(a rune) bool) base.Parser[rune, string] {
	return func(ss stream.SimpleStream[rune]) (string, error) {
		str := ss.TakeWhile(fn, 0)
		return string(str), nil
	}
}
