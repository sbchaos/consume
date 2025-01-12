package base

import (
	"io"

	"github.com/sbchaos/consume/stream"
)

func Satisfy[S any](f stream.Predicate[S]) Parser[S, S] {
	return func(ss stream.SimpleStream[S]) (S, error) {
		var zero S
		token, err := ss.Take()
		if err != nil {
			return zero, err
		}
		if f(token) {
			return token, nil
		}

		return zero, ErrNotMatched
	}
}

// EOF - make sure the input stream has finished
func EOF[S any]() Parser[S, bool] {
	return func(ss stream.SimpleStream[S]) (bool, error) {
		_, err := ss.Take()
		if err == io.EOF {
			return true, nil
		}
		return false, nil
	}
}
