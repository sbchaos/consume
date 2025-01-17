package par_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/stream/strings"
)

func TestBaseParsers(t *testing.T) {
	t.Run("Satisfy", func(t *testing.T) {
		t.Run("returns error on eof", func(t *testing.T) {
			ss := strings.NewStringStream("")

			_, err := par.Parse[rune, rune](ss, par.Satisfy(func(a rune) bool {
				return a == 'a'
			}))
			assert.Error(t, err)
		})
		t.Run("returns error when predicate is false", func(t *testing.T) {
			ss := strings.NewStringStream("b")

			_, err := par.Parse[rune, rune](ss, par.Satisfy(func(a rune) bool {
				return a == 'a'
			}))

			assert.Error(t, err)
			assert.EqualError(t, err, par.ErrNotMatched.Error())
		})
		t.Run("returns value when matched", func(t *testing.T) {
			ss := strings.NewStringStream("a")

			v, err := par.Parse[rune, rune](ss, par.Satisfy(func(a rune) bool {
				return a == 'a'
			}))

			assert.NoError(t, err)
			assert.Equal(t, 'a', v)
		})
	})

	t.Run("TakeWhile", func(t *testing.T) {
		ss := strings.NewStringStream("this is a stream text")
		p := par.TakeWhile(func(a rune) bool {
			return a != 'r'
		}, 0)

		v, err := par.Parse(ss, p)
		assert.NoError(t, err)
		assert.Equal(t, "this is a st", string(v))
	})

	t.Run("TakeUntil", func(t *testing.T) {
		ss := strings.NewStringStream("this is a stream text")
		p := par.TakeUntil([]rune("stream"))

		v, err := par.Parse(ss, p)
		assert.NoError(t, err)
		assert.Equal(t, "this is a ", string(v))
	})

	t.Run("EOF", func(t *testing.T) {
		ss := strings.NewStringStream("")

		eofp := par.EOF[rune]()
		v, err := par.Parse(ss, eofp)
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})
}
