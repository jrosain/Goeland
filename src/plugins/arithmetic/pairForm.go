package arithmetic

type PairOperator string

const (
	EqOperator      PairOperator = "="
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

type PairForm[T, U Form] struct {
	first  T
	second U
	symbol PairOperator
}

func NewPairForm[T, U Form](first T, second U, symbol PairOperator) *PairForm[T, U] {
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

func (pf *PairForm[T, U]) getFactorMap() map[string]int {
	factorMap := make(map[string]int)
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

func (pf *PairForm[T, U]) TrueCopy() *PairForm[T, U] {
	if typedFirst, ok := pf.first.Copy().(T); ok {
		if typedSecond, ok := pf.second.Copy().(U); ok {
			return NewPairForm[T, U](typedFirst, typedSecond, pf.symbol)
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
	*PairForm[Form, Form]
}

func NewSum(first, second Form) *Sum {
	return &Sum{NewPairForm(first, second, SumOperator)}
}

func NewDiff(first, second Form) *Sum {
	return &Sum{NewPairForm(first, Form(NewNeg(second)), SumOperator)}
}

func (s *Sum) Copy() Form {
	return &Sum{s.PairForm.TrueCopy()}
}

var sum = func(first, second int) int {
	return first + second
}

var diff = func(first, second int) int {
	return first - second
}

func (s *Sum) getFactorMap() map[string]int {
	return getFactorMapForFunc[Form, Form](s.PairForm, sum)
}

func getFactorMapForFunc[T, U Form](pf *PairForm[T, U], op func(int, int) int) map[string]int {
	factorMap := make(map[string]int)
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
	*PairForm[*Constant, Form]
}

func NewFactor(factor *Constant, value Form) *Factor {
	return &Factor{NewPairForm[*Constant, Form](factor, value, NoOperator)}
}

func (f *Factor) Copy() Form {
	return &Factor{f.PairForm.TrueCopy()}
}

func (f *Factor) getFactorMap() map[string]int {
	factorMap := make(map[string]int)
	childMap := f.GetSecond().getFactorMap()

	for k, v := range childMap {
		factorMap[k] = v * int(f.first.value)
	}

	return factorMap
}
