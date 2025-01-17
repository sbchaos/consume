package strings

import (
	"unicode"

	"github.com/sbchaos/consume/comb"
	"github.com/sbchaos/consume/par"
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

func AlphaNumeric() par.Parser[rune, string] {
	return CustomString(func(a rune) bool {
		return unicode.IsLetter(a) || unicode.IsDigit(a)
	})
}

func StringWithOptionalQuotes() par.Parser[rune, string] {
	return comb.Choice(
		QuotedString(0, Quotes...),
		AlphaNumeric(),
	)
}

func CustomString(fn func(a rune) bool) par.Parser[rune, string] {
	return func(ss stream.SimpleStream[rune]) (string, error) {
		str := ss.TakeWhile(fn, 0)
		return string(str), nil
	}
}
