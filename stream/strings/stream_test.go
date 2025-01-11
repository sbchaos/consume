package strings_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/consume/stream/strings"
)

func TestStringStream(t *testing.T) {
	content := "this is a stream with lots of string content"
	t.Run("Peek", func(t *testing.T) {
		t.Run("gets one char without advancing", func(t *testing.T) {
			ss := strings.NewStringStream(content)
			peek1, err := ss.Peek()
			assert.NoError(t, err)
			assert.Equal(t, 't', peek1)
			assert.Equal(t, 0, ss.Offset())

			peek2, err := ss.Peek()
			assert.NoError(t, err)
			assert.Equal(t, 't', peek2)
			assert.Equal(t, 0, ss.Offset())
		})
	})
	t.Run("Take", func(t *testing.T) {
		ss := strings.NewStringStream(content)
		take1, err := ss.Take()
		assert.NoError(t, err)
		assert.Equal(t, 't', take1)
		assert.Equal(t, 1, ss.Offset())

		take2, err := ss.Take()
		assert.NoError(t, err)
		assert.Equal(t, 'h', take2)
		assert.Equal(t, 2, ss.Offset())

	})
	t.Run("TakeN", func(t *testing.T) {
		ss := strings.NewStringStream(content)
		take4, err := ss.TakeN(4)
		assert.NoError(t, err)
		assert.Equal(t, 4, ss.Offset())

		assert.Equal(t, "this", string(take4))
		take0, err := ss.TakeN(0)
		assert.NoError(t, err)
		assert.Equal(t, "", string(take0))
		assert.Equal(t, 4, ss.Offset())
	})
	t.Run("TakeWhile", func(t *testing.T) {
		ss := strings.NewStringStream(content)
		while := ss.TakeWhile(func(t rune) bool {
			return t != 'o'
		}, 0)

		assert.Equal(t, "this is a stream with l", string(while))
		assert.Equal(t, 23, ss.Offset())
	})
	t.Run("TakeWhile with escape sequence", func(t *testing.T) {
		ss := strings.NewStringStream(`this is \"a\" string"`)
		while := ss.TakeWhile(func(t rune) bool {
			return t != '"'
		}, '\\')

		assert.Equal(t, "this is \\\"a\\\" string", string(while))
		assert.Equal(t, 20, ss.Offset())
	})
}
