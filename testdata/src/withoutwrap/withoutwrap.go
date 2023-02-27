package main

import "github.com/pkg/errors"

func main() {
	foo(0)
	foo(1)
}

var ErrZeroValue = errors.New("zero value")

var ErrNonZeroValue = errors.New("non zero value")

func foo(bar int) (int, error) {
	if bar == 0 {
		return 0, ErrZeroValue
	}
	if bar == 1 {
		return 0, errors.Wrap(ErrNonZeroValue, "from bar == 1")
	}
	return 0, nil
}
