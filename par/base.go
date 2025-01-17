package par

import (
	"io"

	s "github.com/sbchaos/consume/stream"
)

func Satisfy[S any](f s.Predicate[S]) Parser[S, S] {
	return func(ss s.SimpleStream[S]) (S, error) {
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

func TakeWhile[S any](f s.Predicate[S], escape S) Parser[S, []S] {
	return func(ss s.SimpleStream[S]) ([]S, error) {
		token := ss.TakeWhile(f, escape)
		return token, nil
	}
}

func TakeUntil[S any](seq []S) Parser[S, []S] {
	return func(ss s.SimpleStream[S]) ([]S, error) {
		toks := ss.TakeUntil(seq)
		return toks, nil
	}
}

// EOF - make sure the input s has finished
func EOF[S any]() Parser[S, bool] {
	return func(ss s.SimpleStream[S]) (bool, error) {
		_, err := ss.Take()
		if err == io.EOF {
			return true, nil
		}
		return false, nil
	}
}
