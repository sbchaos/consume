package strings

import (
	"errors"
	"strings"

	"github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
	"github.com/sbchaos/consume/stream"
)

type StringParser par.Parser[rune, string]

type StringMatcher func(s string, t string) bool

var EqualIgnoreCase StringMatcher = strings.EqualFold
var Equals StringMatcher = func(s, t string) bool {
	return s == t
}

func String(expected string, fn StringMatcher) par.Parser[rune, string] {
	return func(ss stream.SimpleStream[rune]) (string, error) {
		n := len(expected)

		tokens, err := ss.TakeN(n)
		if err != nil {
			return "", err
		}

		if fn(string(tokens), expected) {
			return expected, nil
		}

		return "", par.ErrNotMatched
	}
}

func QuotedString(escape rune, quotes ...struct {
	Start rune
	End   rune
}) par.Parser[rune, string] {
	seq := '\\'
	if escape > 0 {
		seq = escape
	}
	return func(ss stream.SimpleStream[rune]) (string, error) {
		first, err := ss.Take()
		if err != nil {
			return "", err
		}

		var end rune = -1
		for _, q1 := range quotes {
			if q1.Start == first {
				end = q1.End
				break
			}
		}

		if end == -1 {
			return "", errors.New("no matching start quote")
		}

		content := ss.TakeWhile(func(t rune) bool {
			return t != end
		}, seq)

		endToken, err := ss.Take()
		if err != nil || endToken != end {
			return "", errors.New("missing end quote for string")
		}

		return string(content), nil
	}
}

func Choice(options []string, fn StringMatcher) par.Parser[rune, string] {
	return func(strm stream.SimpleStream[rune]) (string, error) {
		offset := strm.Offset()
		for _, op1 := range options {
			n := len(op1)

			tokens, err := strm.TakeN(n)
			if err != nil {
				return "", err
			}

			if fn(op1, string(tokens)) {
				return op1, nil
			} else {
				// Reset the offset if no match
				strm.Seek(offset)
			}
		}

		return "", par.ErrNotMatched
	}
}

func Sequence(values []string, fn StringMatcher) par.Parser[rune, string] {
	spaces := char.WhiteSpaces()
	return func(ss stream.SimpleStream[rune]) (string, error) {
		for _, val := range values {
			n := len(val)

			tokens, err := ss.TakeN(n)
			if err != nil {
				return "", err
			}

			if !fn(val, string(tokens)) {
				return "", par.ErrNotMatched
			}

			_, _ = spaces(ss)
		}
		return "", nil
	}
}
