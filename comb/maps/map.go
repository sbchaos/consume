package maps

import (
	c "github.com/sbchaos/consume/comb"
	p "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
)

// ObjectLiteral with parse JSON style object in a map.
// Most used form to represent a map, looks like {"Key":"value"}
func ObjectLiteral(key, val, spaces p.Parser[rune, string]) p.Parser[rune, map[string]string] {
	items := c.ToMap(
		c.Skip(spaces, key),
		c.Skip(spaces, char.Single(':')),
		c.Skip(spaces, val),
		c.Skip(spaces, char.Single(',')),
	)
	mp := c.Between(
		char.Single('{'),
		c.Skip(spaces, items),
		c.Skip(spaces, char.Single('}')),
	)
	return c.Skip(spaces, mp)
}

// AssociatedList will parse it to a map.
// It looks like [(key1,val1), (key2,val2),]
func AssociatedList(key, val, spaces p.Parser[rune, string]) p.Parser[rune, map[string]string] {
	items := c.ToMap(
		c.Skip(spaces, c.Skip(char.Single('('), key)),
		c.Skip(spaces, char.Single(',')),
		c.Skip(spaces, c.SkipAfter(val, char.Single(')'))),
		c.Skip(spaces, char.Single(',')),
	)

	mp := c.Between(
		char.Single('['),
		c.Skip(spaces, items),
		c.Skip(spaces, char.Single(']')),
	)
	return c.Skip(spaces, mp)
}

// KVPair will parse a ini style mapping of items.
// looks like PATH="./home"
func KVPair(key, val, spaces p.Parser[rune, string]) p.Parser[rune, map[string]string] {
	items := c.ToMap(
		c.Skip(spaces, key),
		c.Skip(spaces, char.Single('=')),
		c.Skip(spaces, val),
		c.FMap(func(_ string) rune { return 0 }, spaces),
	)

	return items
}
