package arithmetic

import (
	"github.com/GoelandProver/Goeland/global"
)

type PairOperator string

const (
	EqOperator      PairOperator = "="
	DiffOperator    PairOperator = "!="
	LessOperator    PairOperator = "<"
	LessEqOperator  PairOperator = "<="
	GreatOperator   PairOperator = ">"
	GreatEqOperator PairOperator = ">="
	SumOperator     PairOperator = "+"
	NegOperator     PairOperator = "-"
	NoOperator      PairOperator = ""
)

type Paired interface {
	GetFirst() Form
	GetSecond() Form
	GetSymbol() PairOperator
}

type EvaluablePair[T any] interface {
	GetFirst() Evaluable[T]
	GetSecond() Evaluable[T]
	GetSymbol() PairOperator
}

type PairForm[T, U Evaluable[Numeric]] struct {
	first  T
	second U
	symbol PairOperator
}

func NewPairForm[T, U Evaluable[Numeric]](first T, second U, symbol PairOperator) *PairForm[T, U] {
	return &PairForm[T, U]{first, second, symbol}
}

func (pf *PairForm[T, U]) ToString() string {
	return pf.first.ToString() + " " + string(pf.symbol) + " " + pf.second.ToString()
}

func (pf *PairForm[T, U]) Equals(other any) bool {
	if typed, ok := other.(*PairForm[T, U]); ok {
		return pf.first.Equals(typed.first) && pf.second.Equals(typed.second)
	}
	return false
}

func (pf *PairForm[T, U]) Copy() Form {
	return pf.TrueCopy()
}

func (pf *PairForm[T, U]) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)
	firstChildMap := pf.GetFirst().getFactorMap()
	secondChildMap := pf.GetSecond().getFactorMap()

	for k, v := range firstChildMap {
		factorMap[k] = v
	}

	for k, v := range secondChildMap {
		factorMap[k] = v
	}

	return factorMap
}

func (pf *PairForm[T, U]) ContainsVar() int {
	return pf.first.ContainsVar() + pf.second.ContainsVar()
}

func (pf *PairForm[T, U]) areTwoPartsEqual() bool {
	firstMap := pf.first.getFactorMap()
	secondMap := pf.second.getFactorMap()

	mapsAreEqual := true

	for k, v := range firstMap {
		if secondMap[k] != v {
			mapsAreEqual = false
		}
	}

	for k, v := range secondMap {
		if firstMap[k] != v {
			mapsAreEqual = false
		}
	}

	return mapsAreEqual
}

func (pf *PairForm[T, U]) TrueCopy() *PairForm[Evaluable[Numeric], Evaluable[Numeric]] {
	if typedFirst, ok := pf.first.Copy().(Evaluable[Numeric]); ok {
		if typedSecond, ok := pf.second.Copy().(Evaluable[Numeric]); ok {
			return NewPairForm[Evaluable[Numeric], Evaluable[Numeric]](typedFirst, typedSecond, pf.symbol)
		}
	}

	return nil
}

func (pf *PairForm[T, U]) GetFirst() T {
	return pf.first
}

func (pf *PairForm[T, U]) GetSecond() U {
	return pf.second
}

func (pf *PairForm[T, U]) GetSymbol() PairOperator {
	return pf.symbol
}

type Sum struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewSum(first, second Evaluable[Numeric]) *Sum {
	return &Sum{NewPairForm(first, second, SumOperator)}
}

func NewDiff(first, second Evaluable[Numeric]) *Sum {
	return &Sum{NewPairForm(first, Evaluable[Numeric](NewNeg(second)), SumOperator)}
}

func (s *Sum) Copy() Form {
	return &Sum{s.PairForm.TrueCopy()}
}

var sum = func(first, second float64) float64 {
	return first + second
}

var diff = func(first, second float64) float64 {
	return first - second
}

func (s *Sum) getFactorMap() map[string]float64 {
	return getFactorMapForFunc[Evaluable[Numeric], Evaluable[Numeric]](s.PairForm, sum)
}

func (s *Sum) Evaluate() Numeric {
	return s.first.Evaluate().Sum(s.second.Evaluate())
}

func getFactorMapForFunc[T, U Evaluable[Numeric]](pf *PairForm[T, U], op func(float64, float64) float64) map[string]float64 {
	factorMap := make(map[string]float64)
	firstChildMap := pf.GetFirst().getFactorMap()
	secondChildMap := pf.GetSecond().getFactorMap()

	for k, v := range firstChildMap {
		factorMap[k] = v
	}

	for k, v := range secondChildMap {
		factorMap[k] = op(factorMap[k], v)
	}

	return factorMap
}

type Factor struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewProduct(first, second Evaluable[Numeric]) *Factor {
	if typed, ok := first.(AnyConstant); ok {
		return NewFactor(typed, second)
	} else if typed, ok := second.(AnyConstant); ok {
		return NewFactor(typed, first)
	} else {
		global.PrintPanic("ARI", "Trying to make a product that has two variables. This is forbidden, non-linear arithmetic formulas are not supported.")
		return nil
	}
}

func NewFactor(factor AnyConstant, value Evaluable[Numeric]) *Factor {
	return &Factor{NewPairForm[Evaluable[Numeric], Evaluable[Numeric]](factor, value, NoOperator)}
}

func (f *Factor) Copy() Form {
	return &Factor{f.PairForm.TrueCopy()}
}

func (f *Factor) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)
	childMap := f.GetSecond().getFactorMap()

	if f.GetSecond().ContainsVar() != 0 && f.GetFirst().ContainsVar() != 0 {
		global.PrintPanic("ARI", "Cannot get the factor map of a Factor function with variables on both sides")
	}

	for k, v := range childMap {
		factorMap[k] = v * f.first.Evaluate().Evaluate()
	}

	return factorMap
}

func (f *Factor) Evaluate() Numeric {
	return f.first.Evaluate().Mult(f.second.Evaluate())
}

type Quotient struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
	divisionAlgorithm func(Numeric, Numeric) Numeric
}

func NewQuotient(numerator, denominator Evaluable[Numeric]) *Quotient {
	return &Quotient{NewPairForm(numerator, denominator, "/"), quotient}
}

func NewQuotientE(numerator, denominator Evaluable[Numeric]) *Quotient {
	return &Quotient{NewPairForm(numerator, denominator, "/"), quotientEuclidean}
}

func NewQuotientT(numerator, denominator Evaluable[Numeric]) *Quotient {
	return &Quotient{NewPairForm(numerator, denominator, "/"), quotientTruncation}
}

func NewQuotientF(numerator, denominator Evaluable[Numeric]) *Quotient {
	return &Quotient{NewPairForm(numerator, denominator, "/"), quotientFloor}
}

func (f *Quotient) Copy() Form {
	return &Quotient{f.PairForm.TrueCopy(), f.divisionAlgorithm}
}

func (f *Quotient) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	if f.second.ContainsVar() == 0 {
		firstMap := f.first.getFactorMap()
		for k, v := range firstMap {
			factorMap[k] = v / f.second.Evaluate().Evaluate()
		}
	} else {
		if f.areTwoPartsEqual() {
			factorMap[Unit.ToString()] = One.Evaluate().Evaluate()
		} else {
			global.PrintPanic("ARI", "Trying to get the factor map of a non-linear formula in a quotient function.")
		}
	}

	return factorMap
}

func quotient(f, s Numeric) Numeric {
	return f.Div(s)
}

func quotientEuclidean(f, s Numeric) Numeric {
	res := f.Div(s)

	if s.Gr(ZeroOfType(f)) {
		res = res.Floor()
	} else {
		res = res.Ceil()
	}

	return res
}

func quotientTruncation(f, s Numeric) Numeric {
	return f.Div(s).Trunc()
}

func quotientFloor(f, s Numeric) Numeric {
	return f.Div(s).Floor()
}

func (f *Quotient) Evaluate() Numeric {
	return f.divisionAlgorithm(f.first.Evaluate(), f.second.Evaluate())
}

type Remainder struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
	divisionAlgorithm func(Numeric, Numeric) Numeric
}

func NewRemainderE(numerator, denominator Evaluable[Numeric]) *Remainder {
	return &Remainder{NewPairForm(numerator, denominator, "%"), quotientEuclidean}
}

func NewRemainderT(numerator, denominator Evaluable[Numeric]) *Remainder {
	return &Remainder{NewPairForm(numerator, denominator, "%"), quotientTruncation}
}

func NewRemainderF(numerator, denominator Evaluable[Numeric]) *Remainder {
	return &Remainder{NewPairForm(numerator, denominator, "%"), quotientFloor}
}

func (r *Remainder) Copy() Form {
	return &Remainder{r.PairForm.TrueCopy(), r.divisionAlgorithm}
}

func (r *Remainder) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	if r.ContainsVar() == 0 {
		factorMap[Unit.ToString()] = r.Evaluate().Evaluate()
	} else {
		if r.areTwoPartsEqual() {
			factorMap[Unit.ToString()] = Zero.Evaluate().Evaluate()
		} else {
			global.PrintPanic("ARI", "Trying to get the factor map of a non-linear formula in a remainder function.")
		}
	}

	return factorMap
}

func (r *Remainder) Evaluate() Numeric {
	return r.first.Evaluate().Mod(r.second.Evaluate())
}
