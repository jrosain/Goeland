package arithmetic

import (
	"fmt"
	"math"
)

type Rational struct {
	top int
	bot int
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func (r Rational) Simplify() Rational {
	gcd := gcd(r.top, r.bot)
	return Rational{r.top / gcd, r.bot / gcd}
}

func (r Rational) Negate() Rational {
	return Rational{-r.top, r.bot}
}

func (r Rational) panicOperation() {
	PanicOperation("Rational")
}

func (r Rational) ToString() string {
	return fmt.Sprintf("%v/%v", r.top, r.bot)
}

func (r Rational) Equals(other any) bool {
	if typed, ok := other.(Rational); ok {
		simpleFirst := r.Simplify()
		simpleSecond := typed.Simplify()
		return simpleFirst.top == simpleSecond.top && simpleFirst.bot == simpleSecond.bot
	}
	return false
}

func (r Rational) Copy() Form {
	return r
}

func (r Rational) getFactorMap() map[string]float64 {
	return make(map[string]float64)
}

func (r Rational) Evaluate() float64 {
	return float64(r.top) / float64(r.bot)
}

func (r Rational) ContainsVar() int {
	return 0
}

func (r Rational) Sum(other Numeric) Numeric {
	switch typed := other.(type) {
	case Rational:
		return Rational{(r.top*typed.bot + typed.top*r.bot), r.bot * typed.bot}.Simplify()
	default:
		r.panicOperation()
		return nil
	}
}

func (r Rational) Diff(other Numeric) Numeric {
	switch typed := other.(type) {
	case Rational:
		return Rational{(r.top*typed.bot - typed.top*r.bot), r.bot * typed.bot}.Simplify()
	default:
		r.panicOperation()
		return nil
	}
}

func (r Rational) Mult(other Numeric) Numeric {
	switch typed := other.(type) {
	case Rational:
		return Rational{r.top * typed.top, r.bot * typed.bot}.Simplify()
	default:
		r.panicOperation()
		return nil
	}
}

func (r Rational) Div(other Numeric) Numeric {
	switch typed := other.(type) {
	case Rational:
		return r.Mult(Rational{typed.bot, typed.top})
	default:
		r.panicOperation()
		return nil
	}
}

func (r Rational) Mod(other Numeric) Numeric {
	switch typed := other.(type) {
	case Rational:
		return Rational{r.top % typed.top, r.bot * typed.bot}.Simplify()
	default:
		r.panicOperation()
		return nil
	}
}

func (r Rational) Neg() Numeric {
	return Rational{-r.top, r.bot}.Simplify()
}

func (r Rational) Floor() Numeric {
	r.top = int(math.Floor(float64(r.top) / float64(r.bot)))
	r.bot = 1
	return r
}

func (r Rational) Ceil() Numeric {
	r.top = int(math.Ceil(float64(r.top) / float64(r.bot)))
	r.bot = 1
	return r
}

func (r Rational) Trunc() Numeric {
	r.top = int(math.Trunc(float64(r.top) / float64(r.bot)))
	r.bot = 1
	return r
}

func (r Rational) Round() Numeric {
	r.top = int(math.RoundToEven(float64(r.top) / float64(r.bot)))
	r.bot = 1
	return r
}

func (r Rational) Eq(other Numeric) bool {
	switch typed := other.(type) {
	case Rational:
		return r.Evaluate() == typed.Evaluate()
	default:
		r.panicOperation()
		return false
	}
}

func (r Rational) Gr(other Numeric) bool {
	switch typed := other.(type) {
	case Rational:
		return r.Evaluate() > typed.Evaluate()
	default:
		r.panicOperation()
		return false
	}
}

func (r Rational) Geq(other Numeric) bool {
	switch typed := other.(type) {
	case Rational:
		return r.Evaluate() >= typed.Evaluate()
	default:
		r.panicOperation()
		return false
	}
}

func (r Rational) Le(other Numeric) bool {
	switch typed := other.(type) {
	case Rational:
		return r.Evaluate() < typed.Evaluate()
	default:
		r.panicOperation()
		return false
	}
}

func (r Rational) Leq(other Numeric) bool {
	switch typed := other.(type) {
	case Rational:
		return r.Evaluate() <= typed.Evaluate()
	default:
		r.panicOperation()
		return false
	}
}

func (r Rational) Neq(other Numeric) bool {
	switch typed := other.(type) {
	case Rational:
		return r.Evaluate() != typed.Evaluate()
	default:
		r.panicOperation()
		return false
	}
}

func (r Rational) IsInt() bool {
	return false
}

func (r Rational) IsRat() bool {
	return true
}

func (r Rational) ToInt() Numeric {
	return Integer(r.Floor().Evaluate())
}

func (r Rational) ToRat() Numeric {
	return r
}

func (r Rational) ToReal() Numeric {
	return Real(r.Evaluate())
}

func (r Rational) GetHint() string {
	return "$rat"
}
