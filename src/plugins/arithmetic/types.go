package arithmetic

import (
	"fmt"
	"strconv"
)

type Evaluable[T any] interface {
	Form
	Evaluate() T
}

type String string

type Numeric float64

type Integer int

type Rational struct {
	top int
	bot int
}

type Real float64

func (s String) ToString() string {
	return string(s)
}

func (s String) Equals(other any) bool {
	if typed, ok := other.(String); ok {
		return s.ToString() == typed.ToString()
	}
	return false
}

func (s String) Copy() Form {
	return s
}

func (s String) getFactorMap() map[string]Numeric {
	return make(map[string]Numeric)
}

func (n Numeric) ToString() string {
	return fmt.Sprintf("%v", float64(n))
}

func (n Numeric) Equals(other any) bool {
	if typed, ok := other.(Numeric); ok {
		return float64(n) == float64(typed)
	}
	return false
}

func (n Numeric) Copy() Form {
	return n
}

func (n Numeric) getFactorMap() map[string]Numeric {
	return make(map[string]Numeric)
}

func (n Numeric) Evaluate() Numeric {
	return n
}

func (i Integer) ToString() string {
	return strconv.Itoa(int(i))
}

func (i Integer) Equals(other any) bool {
	if typed, ok := other.(Integer); ok {
		return int(i) == int(typed)
	}
	return false
}

func (i Integer) Copy() Form {
	return i
}

func (i Integer) getFactorMap() map[string]Numeric {
	return make(map[string]Numeric)
}

func (i Integer) Evaluate() Numeric {
	return Numeric(i)
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

func (r Rational) getFactorMap() map[string]Numeric {
	return make(map[string]Numeric)
}

func (r Rational) Evaluate() Numeric {
	return Numeric(float64(r.top) / float64(r.bot))
}

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

func (r Real) getFactorMap() map[string]Numeric {
	return make(map[string]Numeric)
}

func (r Real) Evaluate() Numeric {
	return Numeric(r)
}
