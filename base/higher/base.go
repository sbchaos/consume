package higher

import (
	b "github.com/sbchaos/consume/base"
	s "github.com/sbchaos/consume/stream"
)

func FMap[S, A, B any](f func(A) B, p b.Parser[S, A]) b.Parser[S, B] {
	return func(ss s.SimpleStream[S]) (B, error) {
		var zero B

		result, err := p(ss)
		if err != nil {
			return zero, err
		}

		return f(result), nil
	}
}

// FMap1 - The Functor instance for the Parser
func FMap1[S, A, B any](f func(A) (B, error), p b.Parser[S, A]) b.Parser[S, B] {
	return func(ss s.SimpleStream[S]) (B, error) {
		var zero B

		result, err := p(ss)
		if err != nil {
			return zero, err
		}

		return f(result)
	}
}

// FlatMap - function to take a parser and run it on function returning parsers
func FlatMap[S, A, B any](p1 b.Parser[S, A], f func(A) b.Parser[S, B]) b.Parser[S, B] {
	return func(ss s.SimpleStream[S]) (B, error) {
		var zero B

		result, err := p1(ss) // A, err
		if err != nil {
			return zero, err
		}

		return f(result)(ss)
	}
}

// Wrap will make a parser out of a function
func Wrap[S, A, B any](fn func(A) B) b.Parser[S, func(A) B] {
	return func(_ s.SimpleStream[S]) (func(A) B, error) {
		return fn, nil
	}
}

// Const will create parser for a value
func Const[S any, A any](value A) b.Parser[S, A] {
	return func(_ s.SimpleStream[S]) (A, error) {
		return value, nil
	}
}

// Fail will create a failing parser
func Fail[S any, A any](err error) b.Parser[S, A] {
	return func(_ s.SimpleStream[S]) (A, error) {
		var x A
		return x, err
	}
}

func Apply[S, A, B any](p1 b.Parser[S, func(A) B], p2 b.Parser[S, A]) b.Parser[S, B] {
	return func(ss s.SimpleStream[S]) (B, error) {
		var zero B

		result, err := p2(ss)
		if err != nil {
			return zero, err
		}

		f, err := p1(ss)
		if err != nil {
			return zero, err
		}

		return f(result), nil
	}
}
