package special

import (
	b "github.com/sbchaos/consume/base"
	"github.com/sbchaos/consume/base/higher"
	"github.com/sbchaos/consume/char"
)

// ObjectLiteral with parse JSON style object in a map.
// Most used form to represent a map, looks like {"Key":"value"}
func ObjectLiteral(key, val, spaces b.Parser[rune, string]) b.Parser[rune, map[string]string] {
	items := higher.ToMap(
		higher.Skip(spaces, key),
		higher.Skip(spaces, char.Single(':')),
		higher.Skip(spaces, val),
		higher.Skip(spaces, char.Single(',')),
	)
	mp := higher.Between(
		char.Single('{'),
		higher.Skip(spaces, items),
		higher.Skip(spaces, char.Single('}')),
	)
	return higher.Skip(spaces, mp)
}

// AssociatedList will parse it to a map.
// It looks like [(key1,val1), (key2,val2),]
func AssociatedList(key, val, spaces b.Parser[rune, string]) b.Parser[rune, map[string]string] {
	items := higher.ToMap(
		higher.Skip(spaces, higher.Skip(char.Single('('), key)),
		higher.Skip(spaces, char.Single(',')),
		higher.Skip(spaces, higher.SkipAfter(val, char.Single(')'))),
		higher.Skip(spaces, char.Single(',')),
	)

	mp := higher.Between(
		char.Single('['),
		higher.Skip(spaces, items),
		higher.Skip(spaces, char.Single(']')),
	)
	return higher.Skip(spaces, mp)
}

// KVPair will parse a ini style mapping of items.
// looks like PATH="./home"
func KVPair(key, val, spaces b.Parser[rune, string]) b.Parser[rune, map[string]string] {
	items := higher.ToMap(
		higher.Skip(spaces, key),
		higher.Skip(spaces, char.Single('=')),
		higher.Skip(spaces, val),
		higher.FMap(func(_ string) rune { return 0 }, spaces),
	)

	return items
}
