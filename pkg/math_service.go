package pkg

import (
	"errors"
	"fmt"

	"github.com/go-kit/kit/metrics"
)

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
type MathService struct {
	Operations metrics.Counter
}

// Double implements Doubler interface.
func (m *MathService) Double(a int) (ret int, err error) {
	ret = -1

	defer func() {
		m.Operations.With(
			"operation", "double",
			"success", fmt.Sprintf("%t", err == nil),
		).Add(1)
	}()

	if a < 0 {
		err = errNegative
		return
	}

	ret = 2 * a

	return
}

// Square implements Squarer interface.
func (m *MathService) Square(a int) (ret int, err error) {
	ret = -1

	defer func() {
		m.Operations.With(
			"operation", "square",
			"success", fmt.Sprintf("%t", err == nil),
		).Add(1)
	}()

	if a < 0 {
		err = errNegative
		return
	}

	ret = a * a

	return
}
