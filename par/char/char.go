package char

import (
	p "github.com/sbchaos/consume/par"
	s "github.com/sbchaos/consume/stream"
)

func Single(r rune) p.Parser[rune, rune] {
	return p.Satisfy(func(i rune) bool {
		return i == r
	})
}

func NewLine() p.Parser[rune, rune] {
	return Single('\n')
}

func WhiteSpaces() p.Parser[rune, rune] {
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

func Range[S rune](from, to S) p.Parser[S, S] {
	return p.Satisfy[S](func(x S) bool {
		return x >= from && x <= to
	})
}

func NotRange[S rune](from, to S) p.Parser[S, S] {
	return p.Satisfy[S](func(x S) bool {
		return x < from || x > to
	})
}
