package spaces_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/consume/comb"
	"github.com/sbchaos/consume/comb/spaces"
	"github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/strings"
)

func TestComments(t *testing.T) {
	t.Run("LineComment", func(t *testing.T) {
		t.Run("parses line comment", func(t *testing.T) {
			p1 := spaces.LineComment("//")
			input := "// A head line comment\nActual text"

			val, err := par.ParseString(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " A head line comment", val)
		})
		t.Run("parses multiple line comment", func(t *testing.T) {
			p1 := spaces.LineComment("//")
			input := "// First line comment\n// Second line comment\nActual text"

			val, err := par.ParseString(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " First line comment\n Second line comment", val)
		})
		t.Run("parses other style line comment", func(t *testing.T) {
			p1 := spaces.LineComment("--")
			input := "-- A line comment\nActual text"

			val, err := par.ParseString(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " A line comment", val)
		})
	})

	t.Run("BlockComment", func(t *testing.T) {
		t.Run("parses block comment", func(t *testing.T) {
			p1 := spaces.BlockComment("/*", "*/")
			input := " /* Inline comment*/Actual text"

			val, err := par.ParseString(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " Inline comment", val)
		})
		t.Run("parses multiple comment", func(t *testing.T) {
			p1 := spaces.BlockComment("/*", "*/")
			input := "/* First line comment\n */\n/* Second line comment*/\nActual text"

			val, err := par.ParseString(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " First line comment\n \n Second line comment", val)
		})
		t.Run("parses other style line comment", func(t *testing.T) {
			p1 := spaces.BlockComment("`", "`")
			input := "`A multiline comment\nActual text`End"

			val, err := par.ParseString(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, "A multiline comment\nActual text", val)
		})
	})

	t.Run("BuildSpaceConsumer", func(t *testing.T) {
		t.Run("consumes only space when no comment style", func(t *testing.T) {
			input := "   // After space"
			p1 := spaces.BuildSpaceConsumer("", "", "")

			val, err := par.ParseString(input, comb.Skip(p1, strings.String("// After space", strings.Equals)))
			assert.NoError(t, err)
			assert.Equal(t, "// After space", val)
		})
		t.Run("can define other newline types", func(t *testing.T) {
			input := "   // After space\rText"
			p1 := spaces.BuildSpaceConsumer("//", "/*", "*/", '\r')

			val, err := par.ParseString(input,
				comb.Skip(p1, strings.CustomString(par.Anything[rune])))
			assert.NoError(t, err)
			assert.Equal(t, "Text", val)
		})
		t.Run("returns comments", func(t *testing.T) {
			input := "   // After space\nText"
			p1 := spaces.BuildSpaceConsumer("//", "/*", "*/")

			val, err := par.ParseString(input,
				comb.And(p1, strings.CustomString(par.Anything[rune]), func(comm string, txt string) string {
					return txt + ":comment:" + comm
				}))
			assert.NoError(t, err)
			assert.Equal(t, "Text:comment: After space", val)
		})
		t.Run("runs on sparse content file", func(t *testing.T) {
			txt := `   // After space
			         /* one more */    
	// last one
			Text`
			input := txt
			p1 := spaces.BuildSpaceConsumer("//", "/*", "*/")

			val, err := par.ParseString(input,
				comb.And(p1, strings.CustomString(par.Anything[rune]), func(comm string, txt string) string {
					return txt + ":comment:" + comm
				}))
			assert.NoError(t, err)
			assert.Equal(t, "Text:comment: After space\n one more \n last one", val)
		})
	})
}
