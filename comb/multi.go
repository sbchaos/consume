package comb

import (
	p "github.com/sbchaos/consume/par"
	s "github.com/sbchaos/consume/stream"
)

func Count[S, A any](num int, p1 p.Parser[S, A]) p.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		result := make([]A, 0, num)

		for i := 0; i < num; i++ {
			n, err := p1(ss)
			if err != nil {
				return nil, err
			}

			result = append(result, n)
		}

		return result, nil
	}
}

func Many[S, A any](p1 p.Parser[S, A]) p.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		result := make([]A, 0)

		for {
			x, err := p1(ss)
			if err != nil {
				break
			}

			result = append(result, x)
		}

		return result, nil
	}
}

func Some[S, A any](p1 p.Parser[S, A]) p.Parser[S, []A] {
	return func(ss s.SimpleStream[S]) ([]A, error) {
		cc := Many(p1)

		result, _ := cc(ss)
		if len(result) == 0 {
			return nil, p.ErrNotEnoughElements
		}

		return result, nil
	}
}

func OneOf[S comparable](data ...S) p.Parser[S, S] {
	m := make(map[S]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return p.Satisfy[S](func(x S) bool {
		_, exists := m[x]
		return exists
	})
}

func NoneOf[S comparable](data ...S) p.Parser[S, S] {
	m := make(map[S]struct{})
	for _, x := range data {
		m[x] = struct{}{}
	}

	return p.Satisfy[S](func(x S) bool {
		_, exists := m[x]
		return !exists
	})
}

// ToMap is a special parser to parse a key value list to golang map
// Created for performance reasons
func ToMap[S any, A comparable, B any](keyP p.Parser[S, A], connP p.Parser[S, S],
	valP p.Parser[S, B], sep p.Parser[S, S],
) p.Parser[S, map[A]B] {
	return func(ss s.SimpleStream[S]) (map[A]B, error) {
		mapping := make(map[A]B)

		for {
			idx := ss.Offset()

			key, err := keyP(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, nil
			}

			_, err = connP(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, nil
			}

			val, err := valP(ss)
			if err != nil {
				ss.Seek(idx)
				return mapping, nil
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
