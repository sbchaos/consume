package stream

type Predicate[T any] func(t T) bool

// Stream defines a regular stream of tokens, it can be a stream of chars for a string
type Stream[S any] interface {
	Take() (S, error)
}

type PeekStream[S any] interface {
	Stream[S]
	Peek() (S, error)
}

type MultiStream[S any] interface {
	Stream[S]
	// TakeN To use with functions with fixed length tokens
	TakeN(num int) ([]S, error)
	// TakeWhile To get tokens until a condition is met
	TakeWhile(p Predicate[S], escape S) []S
	// TakeUntil will get the values until the specified sequence is found
	TakeUntil(seq []S) []S
}

type ObservableStream[S any] interface {
	Stream[S]
	// Offset will return a position in stream
	Offset() int
	// Seek is used to move the offset in case of unwinding
	Seek(n int)
}

type ObservableMultiStream[S any] interface {
	MultiStream[S]
	ObservableStream[S]
}

type SimpleStream[S any] interface {
	ObservableMultiStream[S]
	PeekStream[S]
}
