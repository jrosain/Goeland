package arithmetic

import (
	"fmt"
	"math"
)

type Real float64

func (r Real) ToString() string {
	return fmt.Sprintf("%v", r)
}

func (r Real) Equals(other any) bool {
	if typed, ok := other.(Real); ok {
		return float64(r) == float64(typed)
	}
	return false
}

func (r Real) Copy() Form {
	return r
}

func (r Real) getFactorMap() map[string]float64 {
	return make(map[string]float64)
}

func (r Real) Evaluate() float64 {
	return float64(r)
}

func (r Real) ContainsVar() int {
	return 0
}

func (r Real) Sum(other Numeric) Numeric {
	switch typed := other.(type) {
	case Real:
		return Real(r + typed)
	default:
		PanicOperation("Real")
		return nil
	}
}

func (r Real) Diff(other Numeric) Numeric {
	switch typed := other.(type) {
	case Real:
		return Real(r - typed)
	default:
		PanicOperation("Real")
		return nil
	}
}

func (r Real) Mult(other Numeric) Numeric {
	switch typed := other.(type) {
	case Real:
		return Real(r * typed)
	default:
		PanicOperation("Real")
		return nil
	}
}

func (r Real) Div(other Numeric) Numeric {
	switch typed := other.(type) {
	case Real:
		return Real(r / typed)
	default:
		PanicOperation("Real")
		return nil
	}
}

func (r Real) Mod(other Numeric) Numeric {
	switch typed := other.(type) {
	case Real:
		return Real(math.Mod(float64(r), float64(typed)))
	default:
		PanicOperation("Real")
		return nil
	}
}

func (r Real) Neg() Numeric {
	return Real(-r)
}

func (r Real) Floor() Numeric {
	return Real(math.Floor(float64(r)))
}

func (r Real) Ceil() Numeric {
	return Real(math.Ceil(float64(r)))
}

func (r Real) Trunc() Numeric {
	return Real(math.Trunc(float64(r)))
}

func (r Real) Round() Numeric {
	return Real(math.RoundToEven(float64(r)))
}

func (r Real) Eq(other Numeric) bool {
	switch typed := other.(type) {
	case Real:
		return float64(r) == float64(typed)
	default:
		PanicOperation("Real")
		return false
	}
}

func (r Real) Gr(other Numeric) bool {
	switch typed := other.(type) {
	case Real:
		return float64(r) > float64(typed)
	default:
		PanicOperation("Real")
		return false
	}
}

func (r Real) Geq(other Numeric) bool {
	switch typed := other.(type) {
	case Real:
		return float64(r) >= float64(typed)
	default:
		PanicOperation("Real")
		return false
	}
}

func (r Real) Le(other Numeric) bool {
	switch typed := other.(type) {
	case Real:
		return float64(r) < float64(typed)
	default:
		PanicOperation("Real")
		return false
	}
}

func (r Real) Leq(other Numeric) bool {
	switch typed := other.(type) {
	case Real:
		return float64(r) <= float64(typed)
	default:
		PanicOperation("Real")
		return false
	}
}

func (r Real) Neq(other Numeric) bool {
	switch typed := other.(type) {
	case Real:
		return float64(r) != float64(typed)
	default:
		PanicOperation("Real")
		return false
	}
}

func (r Real) IsInt() bool {
	return false
}

func (r Real) IsRat() bool {
	return false
}

func (r Real) ToInt() Numeric {
	return Integer(r.Floor().Evaluate())
}

var d = math.SmallestNonzeroFloat64

// This is a simple approximation that finds a rational in the interval [r - d; r + d]
func (r Real) ToRat() Numeric {
	sign := 1
	if r < 0 {
		sign = -1
		r = Real(float64(r) * float64(sign))
	}

	if r == 0 {
		return Rational{0, 1}
	}

	result := r.getStrictlyPositiveRational()

	result.top *= sign

	return result
}

func (r Real) getStrictlyPositiveRational() Rational {
	result := Rational{1, 1}
	found := false

	for !found {
		if result.Evaluate() > float64(r)+d {
			result.bot += 1
		} else if result.Evaluate() < float64(r)-d {
			result.top += 1
		} else {
			found = true
		}
	}

	return result
}

func (r Real) ToReal() Numeric {
	return r
}
