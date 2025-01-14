package run

import (
	"github.com/sbchaos/consume/base"
	"github.com/sbchaos/consume/stream"
	"github.com/sbchaos/consume/stream/strings"
)

func Parse[S any, A any](ss stream.SimpleStream[S], p base.Parser[S, A]) (A, error) {
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func Parse1[S any, A any](ss stream.ObservableMultiStream[S], p base.Parser1[S, A]) (A, error) {
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func Parse2[S any, A any](ss stream.PeekStream[S], p base.Parser2[S, A]) (A, error) {
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func Parse3[S any, A any](ss stream.ObservableStream[S], p base.Parser3[S, A]) (A, error) {
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func Parse4[S any, A any](ss stream.MultiStream[S], p base.Parser4[S, A]) (A, error) {
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func Parse5[S any, A any](ss stream.Stream[S], p base.Parser5[S, A]) (A, error) {
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}

func ParseString[A any](content string, p base.Parser[rune, A]) (A, error) {
	ss := strings.NewStringStream(content)
	var zero A
	result, err := p(ss)
	if err != nil {
		return zero, err
	}

	return result, nil
}
