package pkg

import "errors"

var errNegative = errors.New("we don't like negativity here")

// Doubler doubles its input, unless it's negative.
type Doubler interface {
	Double(int) (int, error)
}

// Squarer squaqres its input, unless it's negative.
type Squarer interface {
	Square(int) (int, error)
}

// MathService is an implementation of Doubler and Squarer interfaces.
type MathService struct{}

// Double implements Doubler interface.
func (*MathService) Double(a int) (int, error) {
	if a < 0 {
		return -1, errNegative
	}

	return 2 * a, nil
}

// Square implements Squarer interface.
func (*MathService) Square(a int) (int, error) {
	if a < 0 {
		return -1, errNegative
	}

	return a * a, nil
}
