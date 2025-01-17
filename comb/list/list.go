package list

import (
	c "github.com/sbchaos/consume/comb"
	p "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
)

// List It will parse a common list of items
// it will look something like [a, p, c]
func List[A any](item p.Parser[rune, A], spaces p.Parser[rune, string]) p.Parser[rune, []A] {
	items := c.SepBy(char.Single(','), c.Skip(spaces, item))

	lst := c.Between(
		char.Single('['),
		items,
		char.Single(']'),
	)

	return c.Skip(spaces, lst)
}

// Index will parse an expression for selection of item from a list
// will look something like list subscript lst[6]
func Index[A any]() p.Parser[rune, A] {
	return nil
}

// Slice will parse a slice expression for a list
// It will look something like lst[start:end]
func Slice[A any]() p.Parser[rune, A] {
	return nil
}
