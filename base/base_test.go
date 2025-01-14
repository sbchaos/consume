package base_test

import (
	"testing"

	"github.com/sbchaos/consume/base"
	"github.com/sbchaos/consume/run"
	"github.com/sbchaos/consume/stream/strings"
	"github.com/stretchr/testify/assert"
)

func TestBaseParsers(t *testing.T) {
	t.Run("Satisfy", func(t *testing.T) {
		t.Run("returns error on eof", func(t *testing.T) {
			ss := strings.NewStringStream("")

			_, err := run.Parse[rune, rune](ss, base.Satisfy(func(a rune) bool {
				return a == 'a'
			}))
			assert.Error(t, err)
		})
		t.Run("returns error when predicate is false", func(t *testing.T) {
			ss := strings.NewStringStream("b")

			_, err := run.Parse[rune, rune](ss, base.Satisfy(func(a rune) bool {
				return a == 'a'
			}))

			assert.Error(t, err)
			assert.EqualError(t, err, base.ErrNotMatched.Error())
		})
		t.Run("returns value when matched", func(t *testing.T) {
			ss := strings.NewStringStream("a")

			v, err := run.Parse[rune, rune](ss, base.Satisfy(func(a rune) bool {
				return a == 'a'
			}))

			assert.NoError(t, err)
			assert.Equal(t, 'a', v)
		})
	})

	t.Run("TakeWhile", func(t *testing.T) {
		ss := strings.NewStringStream("this is a stream text")
		p := base.TakeWhile(func(a rune) bool {
			return a != 'r'
		}, 0)

		v, err := run.Parse(ss, p)
		assert.NoError(t, err)
		assert.Equal(t, "this is a st", string(v))
	})

	t.Run("TakeUntil", func(t *testing.T) {
		ss := strings.NewStringStream("this is a stream text")
		p := base.TakeUntil([]rune("stream"))

		v, err := run.Parse(ss, p)
		assert.NoError(t, err)
		assert.Equal(t, "this is a ", string(v))
	})

	t.Run("EOF", func(t *testing.T) {
		ss := strings.NewStringStream("")

		eofp := base.EOF[rune]()
		v, err := run.Parse(ss, eofp)
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})
}
