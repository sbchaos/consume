package special_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/consume/base"
	"github.com/sbchaos/consume/base/higher"
	"github.com/sbchaos/consume/run"
	"github.com/sbchaos/consume/stream/strings"
	sp "github.com/sbchaos/consume/strings"
	special "github.com/sbchaos/consume/strings/special"
)

func TestComments(t *testing.T) {
	t.Run("LineComment", func(t *testing.T) {
		t.Run("parses line comment", func(t *testing.T) {
			p1 := special.LineComment("//")
			input := strings.NewStringStream("// A head line comment\nActual text")

			val, err := run.Parse(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " A head line comment", val)
		})
		t.Run("parses multiple line comment", func(t *testing.T) {
			p1 := special.LineComment("//")
			input := strings.NewStringStream("// First line comment\n// Second line comment\nActual text")

			val, err := run.Parse(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " First line comment\n Second line comment", val)
		})
		t.Run("parses other style line comment", func(t *testing.T) {
			p1 := special.LineComment("--")
			input := strings.NewStringStream("-- A line comment\nActual text")

			val, err := run.Parse(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " A line comment", val)
		})
	})

	t.Run("BlockComment", func(t *testing.T) {
		t.Run("parses block comment", func(t *testing.T) {
			p1 := special.BlockComment("/*", "*/")
			input := strings.NewStringStream(" /* Inline comment*/Actual text")

			val, err := run.Parse(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " Inline comment", val)
		})
		t.Run("parses multiple comment", func(t *testing.T) {
			p1 := special.BlockComment("/*", "*/")
			input := strings.NewStringStream("/* First line comment\n */\n/* Second line comment*/\nActual text")

			val, err := run.Parse(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, " First line comment\n \n Second line comment", val)
		})
		t.Run("parses other style line comment", func(t *testing.T) {
			p1 := special.BlockComment("`", "`")
			input := strings.NewStringStream("`A multiline comment\nActual text`End")

			val, err := run.Parse(input, p1)
			assert.NoError(t, err)
			assert.Equal(t, "A multiline comment\nActual text", val)
		})
	})

	t.Run("BuildSpaceConsumer", func(t *testing.T) {
		t.Run("consumes only space when no comment style", func(t *testing.T) {
			input := strings.NewStringStream("   // After space")
			p1 := special.BuildSpaceConsumer("", "", "")

			val, err := run.Parse(input, higher.Skip(p1, sp.String("// After space", sp.Equals)))
			assert.NoError(t, err)
			assert.Equal(t, "// After space", val)
		})
		t.Run("can define other newline types", func(t *testing.T) {
			input := strings.NewStringStream("   // After space\rText")
			p1 := special.BuildSpaceConsumer("//", "/*", "*/", '\r')

			val, err := run.Parse(input,
				higher.Skip(p1, sp.CustomString(base.Anything[rune])))
			assert.NoError(t, err)
			assert.Equal(t, "Text", val)
		})
		t.Run("returns comments", func(t *testing.T) {
			input := strings.NewStringStream("   // After space\nText")
			p1 := special.BuildSpaceConsumer("//", "/*", "*/")

			val, err := run.Parse(input,
				higher.And(p1, sp.CustomString(base.Anything[rune]), func(comm string, txt string) string {
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
			input := strings.NewStringStream(txt)
			p1 := special.BuildSpaceConsumer("//", "/*", "*/")

			val, err := run.Parse(input,
				higher.And(p1, sp.CustomString(base.Anything[rune]), func(comm string, txt string) string {
					return txt + ":comment:" + comm
				}))
			assert.NoError(t, err)
			assert.Equal(t, "Text:comment: After space\n one more \n last one", val)
		})
	})
}
