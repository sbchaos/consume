package list

import (
	"github.com/sbchaos/consume/comb"
	b "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
)

// List It will parse a common list of items
// it will look something like [a, b, c]
func List[A any](item b.Parser[rune, A], spaces b.Parser[rune, string]) b.Parser[rune, []A] {
	items := comb.SepBy(char.Single(','), comb.Skip(spaces, item))

	lst := comb.Between(
		char.Single('['),
		items,
		char.Single(']'),
	)

	return comb.Skip(spaces, lst)
}

// Index will parse an expression for selection of item from a list
// will look something like list subscript lst[6]
func Index[A any]() b.Parser[rune, A] {
	return nil
}

// Slice will parse a slice expression for a list
// It will look something like lst[start:end]
func Slice[A any]() b.Parser[rune, A] {
	return nil
}
