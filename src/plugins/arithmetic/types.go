package arithmetic

import "github.com/GoelandProver/Goeland/global"

type Evaluable[T any] interface {
	Form
	Evaluate() T
}

type String string

type Numeric interface {
	Evaluable[float64]

	Neg() Numeric

	Sum(Numeric) Numeric
	Diff(Numeric) Numeric
	Mult(Numeric) Numeric
	Div(Numeric) Numeric
	Mod(Numeric) Numeric

	Floor() Numeric
	Ceil() Numeric
	Trunc() Numeric
	Round() Numeric

	Eq(Numeric) bool
	Gr(Numeric) bool
	Geq(Numeric) bool
	Le(Numeric) bool
	Leq(Numeric) bool
	Neq(Numeric) bool
}

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

func (s String) getFactorMap() map[string]float64 {
	return make(map[string]float64)
}

func PanicOperation(typeName string) {
	global.PrintPanic("ARI", typeName+"s should only interact with other "+typeName+"s")
}
