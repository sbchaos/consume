package par

import "github.com/sbchaos/consume/stream"

type Parser[S any, A any] func(stream stream.SimpleStream[S]) (A, error)

// Parser1 to N are Parsers for each type of stream
type Parser1[S any, A any] func(stream stream.ObservableMultiStream[S]) (A, error)
type Parser2[S any, A any] func(stream stream.PeekStream[S]) (A, error)
type Parser3[S any, A any] func(stream stream.ObservableStream[S]) (A, error)
type Parser4[S any, A any] func(stream stream.MultiStream[S]) (A, error)
type Parser5[S any, A any] func(stream stream.Stream[S]) (A, error)

// Anything - return true anyway.
func Anything[T any](x T) bool { return true }

// Nothing - return false anyway.
func Nothing[T any](x T) bool { return false }

type Composer[A, B, C any] func(A, B) C
