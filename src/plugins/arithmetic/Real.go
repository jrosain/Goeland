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

func (r Real) ToRat() Numeric {
	current := float64(r)
	counter := 1

	for math.Trunc(current) != current {
		counter *= 10
		current *= 10
	}

	return Rational{int(current), counter}
}

func (r Real) ToReal() Numeric {
	return r
}

func (r Real) GetHint() string {
	return "$real"
}
