package special_test

import (
	"testing"

	"github.com/sbchaos/consume/base/higher"
	"github.com/sbchaos/consume/char"
	"github.com/sbchaos/consume/run"
	"github.com/sbchaos/consume/special"
	"github.com/sbchaos/consume/strings"
	"github.com/stretchr/testify/assert"
)

func TestMaps(t *testing.T) {
	t.Run("ObjectLiteral", func(t *testing.T) {
		input := ` {
		"name": "object",
		"property": "string"
		}`
		sp := higher.FMap(func(_ rune) string { return "" }, char.WhiteSpaces())
		str1 := strings.StringWithOptionalQuotes()
		p1 := special.ObjectLiteral(str1, str1, sp)

		val, err := run.ParseString(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "object", val["name"])
		assert.Equal(t, "string", val["property"])
	})
	t.Run("AssociatedList", func(t *testing.T) {
		input := ` [
			(name, object),
			(property, string)]`
		sp := higher.FMap(func(_ rune) string { return "" }, char.WhiteSpaces())
		str1 := strings.StringWithOptionalQuotes()
		p1 := special.AssociatedList(str1, str1, sp)

		val, err := run.ParseString(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "object", val["name"])
		assert.Equal(t, "string", val["property"])
	})
	t.Run("KVPair", func(t *testing.T) {
		input := `
			name=object
			property=string`
		sp := higher.FMap(func(_ rune) string { return "" }, char.WhiteSpaces())
		str1 := strings.StringWithOptionalQuotes()
		p1 := special.KVPair(str1, str1, sp)

		val, err := run.ParseString(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "object", val["name"])
		assert.Equal(t, "string", val["property"])
	})
}
