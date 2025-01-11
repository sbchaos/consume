package base

import "errors"

var (
	ErrNotEnoughElements = errors.New("not enough of elements")
	ErrNotMatched        = errors.New("nothing matched")
)
