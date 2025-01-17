package comb

import (
	p "github.com/sbchaos/consume/par"
	s "github.com/sbchaos/consume/stream"
)

// Try will convert any parser to one with ability of record and reset in case of error
func Try[S any, A any](p1 p.Parser[S, A]) p.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		idx := ss.Offset()

		result, err := p1(ss)
		if err != nil {
			var zero A
			ss.Seek(idx)
			return zero, err
		}

		return result, nil
	}
}

// Choice - searches for a combinator that works successfully on the input data.
func Choice[S any, A any](ps ...p.Parser[S, A]) p.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		var zero A
		for _, p := range ps {
			result, err := Try(p)(ss)
			if err == nil {
				return result, err
			}
		}

		return zero, p.ErrNotMatched
	}
}

func Sequence[S, A any](ps ...p.Parser[S, A]) p.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		arr := make([]A, len(ps))

		for _, c := range ps {
			v, err := c(ss)
			if err != nil {
				return nil, err
			}

			arr = append(arr, v)
		}

		return arr, nil
	}
}

func Optional[S any, A any](p1 p.Parser[S, A], def A) p.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		idx := ss.Offset()
		result, err := p1(ss)
		if err != nil {
			ss.Seek(idx)
			return def, nil
		}

		return result, nil
	}
}

// And is simplified version of combinator for 2 parsers, alternative to flatmap -> flatmap
// A special case of 2 successive apply of flatmaps
func And[S, A, B, C any](x p.Parser[S, A], y p.Parser[S, B], compose p.Composer[A, B, C]) p.Parser[S, C] {
	return func(ss s.SimpleStream[S]) (C, error) {
		var zero C
		first, err := x(ss)
		if err != nil {
			return zero, err
		}

		second, err := y(ss)
		if err != nil {
			return zero, err
		}

		return compose(first, second), nil
	}
}

// SepBy - Parses a list of items separated by sep
// eg a, p, c
func SepBy[S, A, B any](sep p.Parser[S, B], item p.Parser[S, A]) p.Parser[S, []A] {
	c := Try(And(sep, item, func(_ B, x A) A { return x }))

	return func(ss s.SimpleStream[S]) ([]A, error) {
		result := make([]A, 0)

		token, err := item(ss)
		if err != nil {
			return result, nil
		}
		result = append(result, token)

		for {
			token, err = c(ss)
			if err != nil {
				break
			}

			result = append(result, token)
		}

		return result, nil
	}
}

// Between parses a sequence of input combinators, skip first and last
// eg ( items_in_between )
func Between[S, A, B, C any](pre p.Parser[S, A], c p.Parser[S, B], suf p.Parser[S, C]) p.Parser[S, B] {
	return func(ss s.SimpleStream[S]) (B, error) {
		var zero B
		_, err := pre(ss)
		if err != nil {
			return zero, err
		}

		item, err := c(ss)
		if err != nil {
			return zero, err
		}

		_, err = suf(ss)
		if err != nil {
			return zero, err
		}

		return item, nil
	}
}

// Surround a parser with another, generally used to wrap whitespace around a parser
func Surround[S, A, B any](w p.Parser[S, B], p p.Parser[S, A]) p.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		var zero A

		_, err := w(ss)
		if err != nil {
			return zero, err
		}

		res, err := p(ss)
		if err != nil {
			return zero, err
		}

		_, err = w(ss)
		if err != nil {
			return zero, err
		}

		return res, nil
	}
}

// Skip is a Special case combinator
func Skip[S, A, B any](skip p.Parser[S, B], p p.Parser[S, A]) p.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		var zero A
		_, err := skip(ss)
		if err != nil {
			return zero, err
		}

		return p(ss)
	}
}

func SkipAfter[S, A, B any](p1 p.Parser[S, A], skip p.Parser[S, B]) p.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		var zero A
		res, err := p1(ss)
		if err != nil {
			return zero, err
		}

		_, err = skip(ss)
		if err != nil {
			return zero, err
		}

		return res, nil
	}
}
