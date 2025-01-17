package char

import (
	b "github.com/sbchaos/consume/par"
	s "github.com/sbchaos/consume/stream"
)

func Single(r rune) b.Parser[rune, rune] {
	return b.Satisfy(func(i rune) bool {
		return i == r
	})
}

func NewLine() b.Parser[rune, rune] {
	return Single('\n')
}

func WhiteSpaces() b.Parser[rune, rune] {
	return func(stream s.SimpleStream[rune]) (rune, error) {
		matcher := func(t rune) bool {
			if t == '\n' || t == '\r' {
				return true
			}
			if t == ' ' || t == '\t' {
				return true
			}
			return false
		}

		_ = stream.TakeWhile(matcher, '\\')
		return 0, nil
	}
}

func Range[S rune](from, to S) b.Parser[S, S] {
	return b.Satisfy[S](func(x S) bool {
		return x >= from && x <= to
	})
}

func NotRange[S rune](from, to S) b.Parser[S, S] {
	return b.Satisfy[S](func(x S) bool {
		return x < from || x > to
	})
}
