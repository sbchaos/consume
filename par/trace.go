package par

import (
	"fmt"

	"github.com/sbchaos/consume/stream"
)

var Debug = false

type Logger interface {
	Log(args ...any)
}

func Trace[S any, A any](l Logger, name string, c Parser[S, A]) Parser[S, A] {
	return func(ss stream.SimpleStream[S]) (A, error) {
		if Debug {
			l.Log(name)
			l.Log("\toffset before:", ss.Offset())
		}

		result, err := c(ss)
		if Debug {
			l.Log("\tposition after:", ss.Offset())
			if err != nil {
				var zero A
				l.Log("\tnot parsed:", name, result, err)
				return zero, err
			}

			l.Log("\tparsed:", fmt.Sprintf("%#v", result))
		}
		return result, err
	}
}
