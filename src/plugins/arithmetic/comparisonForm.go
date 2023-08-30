package arithmetic

type ComparisonForm interface {
	Form
	Paired
	Normalize() ComparisonForm
	Reverse() ComparisonForm
	Equalize() ComparisonForm
	Simplify() ComparisonForm
	isClosure() bool
}

func NewComparaisonForm(first, second Form, symbol PairOperator) ComparisonForm {
	switch symbol {
	case EqOperator:
		return NewEq(first, second)
	case LessOperator:
		return NewLess(first, second)
	case LessEqOperator:
		return NewLessEq(first, second)
	case GreatOperator:
		return NewGreat(first, second)
	case GreatEqOperator:
		return NewGreatEq(first, second)
	default:
		return nil
	}
}

func buildComparisonComponentsFrom(compForm ComparisonForm) ComparisonForm {
	formatted := compForm.Normalize().Reverse().Equalize()
	factorMap := formatted.getFactorMap()

	firstDone := false
	var constant *Constant
	var form Form

	for k, v := range factorMap {
		factor := NewConstant(v)

		if k == Unit.ToString() {
			constant = factor
		} else if firstDone {
			form = NewSum(form, NewFactor(factor, NewVariable(k)))
		} else {
			form = NewFactor(factor, NewVariable(k))
		}
	}

	return NewComparaisonForm(form, NewNeg(constant), formatted.GetSymbol())
}

type Eq struct {
	*PairForm[Form, Form]
}

func NewEq(first, second Form) *Eq {
	return &Eq{NewPairForm(first, second, EqOperator)}
}

func (e *Eq) TrueCopy() *Eq {
	return &Eq{e.PairForm.TrueCopy()}
}

func (e *Eq) Copy() Form {
	return e.TrueCopy()
}

func (e *Eq) getFactorMap() map[string]int {
	return getFactorMapForFunc[Form, Form](e.PairForm, diff)
}

func (e *Eq) Normalize() ComparisonForm {
	return NewEq(NewSum(e.first, NewNeg(e.second)), Zero)
}

func (e *Eq) Reverse() ComparisonForm {
	return e.TrueCopy()
}

func (e *Eq) Equalize() ComparisonForm {
	return e.TrueCopy()
}

func (e *Eq) Simplify() ComparisonForm {
	return buildComparisonComponentsFrom(e)
}

func (e *Eq) isClosure() bool {
	return !e.first.Equals(e.second)
}

func getBothIntegers(comp ComparisonForm) (first int, second int, areBothInts bool) {
	if firstTyped, ok := comp.GetFirst().(*Constant); ok {
		if secondTyped, ok := comp.GetSecond().(*Constant); ok {
			return int(firstTyped.value), int(secondTyped.value), true
		}
	}

	return 0, 0, false
}

type Less struct {
	*PairForm[Form, Form]
}

func NewLess(first, second Form) *Less {
	return &Less{NewPairForm(first, second, LessOperator)}
}

func (l *Less) TrueCopy() *Less {
	return &Less{l.PairForm.TrueCopy()}
}

func (l *Less) Copy() Form {
	return l.TrueCopy()
}

func (l *Less) getFactorMap() map[string]int {
	return getFactorMapForFunc[Form, Form](l.PairForm, diff)
}

func (l *Less) Normalize() ComparisonForm {
	return NewLess(NewDiff(l.first, l.second), Zero)
}

func (l *Less) Reverse() ComparisonForm {
	return NewGreatEq(l.first, l.second)
}

func (l *Less) Equalize() ComparisonForm {
	return NewLessEq(NewSum(l.first, NewConstant(1)), l.second)
}

func (l *Less) Simplify() ComparisonForm {
	return buildComparisonComponentsFrom(l)
}

func (l *Less) isClosure() bool {
	first, second, areBothIntegers := getBothIntegers(l)
	if areBothIntegers {
		return first >= second
	}

	return false
}

type LessEq struct {
	*PairForm[Form, Form]
}

func NewLessEq(first, second Form) *LessEq {
	return &LessEq{NewPairForm(first, second, LessEqOperator)}
}

func (le *LessEq) TrueCopy() *LessEq {
	return &LessEq{le.PairForm.TrueCopy()}
}

func (le *LessEq) Copy() Form {
	return le.TrueCopy()
}

func (le *LessEq) getFactorMap() map[string]int {
	return getFactorMapForFunc[Form, Form](le.PairForm, diff)
}

func (le *LessEq) Normalize() ComparisonForm {
	return NewLessEq(NewDiff(le.first, le.second), Zero)
}

func (le *LessEq) Reverse() ComparisonForm {
	return NewGreat(le.first, le.second)
}

func (le *LessEq) Equalize() ComparisonForm {
	return le.TrueCopy()
}

func (le *LessEq) Simplify() ComparisonForm {
	return buildComparisonComponentsFrom(le)
}

func (le *LessEq) isClosure() bool {
	first, second, areBothIntegers := getBothIntegers(le)
	if areBothIntegers {
		return first > second
	}

	return false
}

type Great struct {
	*PairForm[Form, Form]
}

func NewGreat(first, second Form) *Great {
	return &Great{NewPairForm(first, second, GreatOperator)}
}

func (g *Great) TrueCopy() *Great {
	return &Great{g.PairForm.TrueCopy()}
}

func (g *Great) Copy() Form {
	return g.TrueCopy()
}

func (g *Great) getFactorMap() map[string]int {
	return getFactorMapForFunc[Form, Form](g.PairForm, diff)
}

func (g *Great) Normalize() ComparisonForm {
	return NewGreat(NewDiff(g.first, g.second), Zero)
}

func (g *Great) Reverse() ComparisonForm {
	return NewLessEq(g.first, g.second)
}

func (g *Great) Equalize() ComparisonForm {
	return NewGreatEq(NewDiff(g.first, NewConstant(1)), g.second)
}

func (g *Great) Simplify() ComparisonForm {
	return buildComparisonComponentsFrom(g)
}

func (g *Great) isClosure() bool {
	first, second, areBothIntegers := getBothIntegers(g)
	if areBothIntegers {
		return first <= second
	}

	return false
}

type GreatEq struct {
	*PairForm[Form, Form]
}

func NewGreatEq(first, second Form) *GreatEq {
	return &GreatEq{NewPairForm(first, second, GreatEqOperator)}
}

func (ge *GreatEq) TrueCopy() *GreatEq {
	return &GreatEq{ge.PairForm.TrueCopy()}
}

func (ge *GreatEq) Copy() Form {
	return ge.TrueCopy()
}

func (ge *GreatEq) getFactorMap() map[string]int {
	return getFactorMapForFunc[Form, Form](ge.PairForm, diff)
}

func (ge *GreatEq) Normalize() ComparisonForm {
	return NewGreatEq(NewDiff(ge.first, ge.second), Zero)
}

func (ge *GreatEq) Reverse() ComparisonForm {
	return NewLess(ge.first, ge.second)
}

func (ge *GreatEq) Equalize() ComparisonForm {
	return ge.TrueCopy()
}

func (ge *GreatEq) Simplify() ComparisonForm {
	return buildComparisonComponentsFrom(ge)
}

func (ge *GreatEq) isClosure() bool {
	first, second, areBothIntegers := getBothIntegers(ge)
	if areBothIntegers {
		return first < second
	}

	return false
}
