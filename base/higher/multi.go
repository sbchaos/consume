package higher

import (
	"fmt"

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

// A special parser to parse a key value list to golang map
// Created for performance reasons
func ToMap[S any, A comparable, B any](keyP b.Parser[S, A], connP b.Parser[S, S],
	valP b.Parser[S, B], sep b.Parser[S, S],
) b.Parser[S, map[A]B] {
	return func(ss s.SimpleStream[S]) (map[A]B, error) {
		mapping := make(map[A]B)

		for {
			idx := ss.Offset()

			key, err := keyP(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, fmt.Errorf("error in parsing map key: %w", err)
			}

			_, err = connP(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, err
			}

			val, err := valP(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, fmt.Errorf("error parsing map value: %w", err)
			}

			mapping[key] = val

			idx = ss.Offset()
			_, err = sep(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, nil
			}
		}
	}
}
