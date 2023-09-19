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

func (pf *PairForm[T, U]) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)
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

var sum = func(first, second Numeric) Numeric {
	return first + second
}

var diff = func(first, second Numeric) Numeric {
	return first - second
}

func (s *Sum) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc[Evaluable[Numeric], Evaluable[Numeric]](s.PairForm, sum)
}

func (s *Sum) Evaluate() Numeric {
	return s.first.Evaluate() + s.second.Evaluate()
}

func getFactorMapForFunc[T, U Evaluable[Numeric]](pf *PairForm[T, U], op func(Numeric, Numeric) Numeric) map[string]Numeric {
	factorMap := make(map[string]Numeric)
	firstChildMap := pf.GetFirst().getFactorMap()
	secondChildMap := pf.GetSecond().getFactorMap()

	for k, v := range firstChildMap {
		factorMap[k] = v
	}

	for k, v := range secondChildMap {
		if _, found := factorMap[k]; !found {
			factorMap[k] = 0
		}

		factorMap[k] = op(factorMap[k], v)
	}

	return factorMap
}

type Factor struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewProduct(first, second Evaluable[Numeric]) *Factor {
	if typed, ok := first.(*Constant); ok {
		return NewFactor(typed, second)
	} else if typed, ok := second.(*Constant); ok {
		return NewFactor(typed, first)
	} else {
		global.PrintPanic("ARI", "Trying to make a product that has two variables. This is forbidden, non-linear arithmetic formulas are not supported.")
		return nil
	}
}

func NewFactor(factor *Constant, value Evaluable[Numeric]) *Factor {
	return &Factor{NewPairForm[Evaluable[Numeric], Evaluable[Numeric]](factor, value, NoOperator)}
}

func (f *Factor) Copy() Form {
	return &Factor{f.PairForm.TrueCopy()}
}

func (f *Factor) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)
	childMap := f.GetSecond().getFactorMap()

	for k, v := range childMap {
		factorMap[k] = v * f.first.Evaluate()
	}

	return factorMap
}

func (f *Factor) Evaluate() Numeric {
	return f.first.Evaluate() * f.second.Evaluate()
}
