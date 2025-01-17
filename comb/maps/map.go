package maps

import (
	higher2 "github.com/sbchaos/consume/comb"
	b "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
)

// ObjectLiteral with parse JSON style object in a map.
// Most used form to represent a map, looks like {"Key":"value"}
func ObjectLiteral(key, val, spaces b.Parser[rune, string]) b.Parser[rune, map[string]string] {
	items := higher2.ToMap(
		higher2.Skip(spaces, key),
		higher2.Skip(spaces, char.Single(':')),
		higher2.Skip(spaces, val),
		higher2.Skip(spaces, char.Single(',')),
	)
	mp := higher2.Between(
		char.Single('{'),
		higher2.Skip(spaces, items),
		higher2.Skip(spaces, char.Single('}')),
	)
	return higher2.Skip(spaces, mp)
}

// AssociatedList will parse it to a map.
// It looks like [(key1,val1), (key2,val2),]
func AssociatedList(key, val, spaces b.Parser[rune, string]) b.Parser[rune, map[string]string] {
	items := higher2.ToMap(
		higher2.Skip(spaces, higher2.Skip(char.Single('('), key)),
		higher2.Skip(spaces, char.Single(',')),
		higher2.Skip(spaces, higher2.SkipAfter(val, char.Single(')'))),
		higher2.Skip(spaces, char.Single(',')),
	)

	mp := higher2.Between(
		char.Single('['),
		higher2.Skip(spaces, items),
		higher2.Skip(spaces, char.Single(']')),
	)
	return higher2.Skip(spaces, mp)
}

// KVPair will parse a ini style mapping of items.
// looks like PATH="./home"
func KVPair(key, val, spaces b.Parser[rune, string]) b.Parser[rune, map[string]string] {
	items := higher2.ToMap(
		higher2.Skip(spaces, key),
		higher2.Skip(spaces, char.Single('=')),
		higher2.Skip(spaces, val),
		higher2.FMap(func(_ string) rune { return 0 }, spaces),
	)

	return items
}
