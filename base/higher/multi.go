package higher

import (
	b "github.com/sbchaos/consume/base"
	s "github.com/sbchaos/consume/stream"
)

func Count[S, A any](num int, p b.Parser[S, A]) b.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		result := make([]A, 0, num)

		for i := 0; i < num; i++ {
			n, err := p(ss)
			if err != nil {
				return nil, err
			}

			result = append(result, n)
		}

		return result, nil
	}
}

func Many[S, A any](p b.Parser[S, A]) b.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		result := make([]A, 0)

		for {
			x, err := p(ss)
			if err != nil {
				break
			}

			result = append(result, x)
		}

		return result, nil
	}
}

func Some[S, A any](p b.Parser[S, A]) b.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		cc := Many(p)

		result, _ := cc(ss)
		if len(result) == 0 {
			return nil, b.ErrNotEnoughElements
		}

		return result, nil
	}
}

// Special case combinator
func Skip[S, A, B any](skip b.Parser[S, B], p b.Parser[S, A]) b.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		var zero A
		_, err := skip(ss)
		if err != nil {
			return zero, err
		}

		return p(ss)
	}
}

// Special case combinator
func SkipAfter[S, A, B any](p b.Parser[S, A], skip b.Parser[S, B]) b.Parser[S, A] {
	return func(ss s.SimpleStream[S]) (A, error) {
		var zero A
		res, err := p(ss)
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

func OneOf[S comparable](data ...S) b.Parser[S, S] {
	m := make(map[S]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return b.Satisfy[S](func(x S) bool {
		_, exists := m[x]
		return exists
	})
}

func NoneOf[S comparable](data ...S) b.Parser[S, S] {
	m := make(map[S]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return b.Satisfy[S](func(x S) bool {
		_, exists := m[x]
		return !exists
	})
}
