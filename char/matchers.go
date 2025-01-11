package char

import (
	"unicode"

	b "github.com/sbchaos/consume/base"
)

// IsDigit parses decimal digit UTF-8 characters.
func IsDigit() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsDigit)
}

// IsNumber parses UTF-8 number characters.
func IsNumber() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsNumber)
}

// IsLetter parses letter UTF-8 characters.
func IsLetter() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsLetter)
}

// IsLower parses UTF-8 in upper case char
func IsUpper() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsUpper)
}

// IsLower parses UTF-8 character in lower case.
func IsLower() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsLower)
}

// IsTitle parses UTF-8 character in title case.
func IsTitle() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsTitle)
}

// IsSpace parses UTF-8 space character.
func IsSpace() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsSpace)
}

// IsPunct parses UTF-8 punctuation character.
func IsPunct() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsPunct)
}

// IsPrint parses printable UTF-8 characters.
func IsPrint() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsPrint)
}

// IsSymbol parsse UTF-8 symbolic character.
func IsSymbol() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsSymbol)
}

// IsControl parses control UTF-8 characters.
func IsControl() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsControl)
}

// IsGraphic parses graphic UTF-8 characters.
func IsGraphic() b.Parser[rune, rune] {
	return b.Satisfy(unicode.IsGraphic)
}
