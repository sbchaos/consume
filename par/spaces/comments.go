package spaces

import (
	"strings"

	"github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/stream"
)

// LineComment will parse a line comment
// It will consume everything until `\n`, but it will not include newline in content
// eg start of this line in golang or -- in sql
func LineComment(seq string) par.Parser[rune, string] {
	return BuildSpaceConsumer(seq, "", "")
}

// BlockComment will parse everything until the end block is found, it does not allow nesting
// For parsing things like /* in golang or """ in python
func BlockComment(start, end string) par.Parser[rune, string] {
	return BuildSpaceConsumer("", start, end)
}

// BuildSpaceConsumer can be used to build a complex space consumer which can
// consume spaces, tabs, newlines, block and line comments
// returns the content of comment or else the empty string, \n is default newline
func BuildSpaceConsumer(lineStart, blockStart, blockEnd string, newline ...rune) par.Parser[rune, string] {
	moreNewLine := len(newline) > 0
	ll := len(lineStart)
	bl := len(blockStart)

	return func(ss stream.SimpleStream[rune]) (string, error) {
		var b strings.Builder
		for {
			_ = ss.TakeWhile(func(a rune) bool {
				if a == ' ' || a == '\t' || a == '\n' {
					return true
				}

				return moreNewLine && strings.ContainsRune(string(newline), a)
			}, 0)

			if ll > 0 {
				idx := ss.Offset()
				v, err := ss.TakeN(ll)
				if err != nil || !strings.EqualFold(string(v), lineStart) {
					ss.Seek(idx)
				} else {
					lc := ss.TakeWhile(func(t rune) bool {
						if t == '\n' {
							return false
						}
						if moreNewLine && strings.ContainsRune(string(newline), t) {
							return false
						}
						return true
					}, 0)

					if b.Len() > 0 {
						b.WriteRune('\n')
					}
					b.WriteString(string(lc))

					continue
				}
			}

			if bl > 0 {
				idx := ss.Offset()
				v, err := ss.TakeN(bl)
				if err != nil || !strings.EqualFold(string(v), blockStart) {
					ss.Seek(idx)
				} else {
					lc := ss.TakeUntil([]rune(blockEnd))

					_, err := ss.TakeN(len(blockEnd))
					if err == nil {
						if b.Len() > 0 {
							b.WriteRune('\n')
						}
						b.WriteString(string(lc))

						continue // continue loop on match
					} else {
						// unlikely, we checked on this above
						ss.Seek(idx)
					}
				}
			}

			break
		}

		return b.String(), nil
	}
}
