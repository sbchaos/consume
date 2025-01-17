package maps_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/consume/comb"
	"github.com/sbchaos/consume/comb/maps"
	"github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
	"github.com/sbchaos/consume/par/strings"
)

func TestMaps(t *testing.T) {
	t.Run("ObjectLiteral", func(t *testing.T) {
		input := ` {
		"name": "object",
		"property": "string"
		}`
		sp := comb.FMap(func(_ rune) string { return "" }, char.WhiteSpaces())
		str1 := strings.StringWithOptionalQuotes()
		p1 := maps.ObjectLiteral(str1, str1, sp)

		val, err := par.ParseString(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "object", val["name"])
		assert.Equal(t, "string", val["property"])
	})
	t.Run("AssociatedList", func(t *testing.T) {
		input := ` [
			(name, object),
			(property, string)]`
		sp := comb.FMap(func(_ rune) string { return "" }, char.WhiteSpaces())
		str1 := strings.StringWithOptionalQuotes()
		p1 := maps.AssociatedList(str1, str1, sp)

		val, err := par.ParseString(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "object", val["name"])
		assert.Equal(t, "string", val["property"])
	})
	t.Run("KVPair", func(t *testing.T) {
		input := `
			name=object
			property=string`
		sp := comb.FMap(func(_ rune) string { return "" }, char.WhiteSpaces())
		str1 := strings.StringWithOptionalQuotes()
		p1 := maps.KVPair(str1, str1, sp)

		val, err := par.ParseString(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "object", val["name"])
		assert.Equal(t, "string", val["property"])
	})
}
