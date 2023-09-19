package arithmetic

type ComparisonForm interface {
	Evaluable[bool]
	EvaluablePair[Numeric]
	Normalize() ComparisonForm
	Reverse() ComparisonForm
	Equalize() ComparisonForm
	Simplify() ComparisonForm
}

func NewComparaisonForm(first, second Evaluable[Numeric], symbol PairOperator) ComparisonForm {
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
	normalized := compForm.Normalize()
	reversed := normalized.Reverse()
	equalized := reversed.Equalize()
	factorMap := equalized.getFactorMap()

	firstDone := false
	var constant *Constant
	var form Evaluable[Numeric]

	for k, v := range factorMap {
		factor := NewConstant(v)

		if k == Unit.ToString() {
			constant = factor
		} else if firstDone {
			form = NewSum(form, NewFactor(factor, NewVariable(k)))
		} else {
			form = NewFactor(factor, NewVariable(k))
			firstDone = true
		}
	}

	return NewComparaisonForm(form, NewNeg(constant), equalized.GetSymbol())
}

type Eq struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewEq(first, second Evaluable[Numeric]) *Eq {
	return &Eq{NewPairForm(first, second, EqOperator)}
}

func (e *Eq) TrueCopy() *Eq {
	return &Eq{e.PairForm.TrueCopy()}
}

func (e *Eq) Copy() Form {
	return e.TrueCopy()
}

func (e *Eq) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc(e.PairForm, diff)
}

func (e *Eq) Evaluate() bool {
	return e.first.Evaluate() == e.second.Evaluate()
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

type NotEq struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewNotEq(first, second Evaluable[Numeric]) *NotEq {
	return &NotEq{NewPairForm(first, second, EqOperator)}
}

func (d *NotEq) TrueCopy() *NotEq {
	return &NotEq{d.PairForm.TrueCopy()}
}

func (d *NotEq) Copy() Form {
	return d.TrueCopy()
}

func (d *NotEq) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc(d.PairForm, diff)
}

func (d *NotEq) Evaluate() bool {
	return d.first.Evaluate() == d.second.Evaluate()
}

func (d *NotEq) Normalize() ComparisonForm {
	return NewNotEq(NewSum(d.first, NewNeg(d.second)), Zero)
}

func (d *NotEq) Reverse() ComparisonForm {
	return d.TrueCopy()
}

func (d *NotEq) Equalize() ComparisonForm {
	return d.TrueCopy()
}

func (d *NotEq) Simplify() ComparisonForm {
	return buildComparisonComponentsFrom(d)
}

type Less struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewLess(first, second Evaluable[Numeric]) *Less {
	return &Less{NewPairForm(first, second, LessOperator)}
}

func (l *Less) TrueCopy() *Less {
	return &Less{l.PairForm.TrueCopy()}
}

func (l *Less) Copy() Form {
	return l.TrueCopy()
}

func (l *Less) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc(l.PairForm, diff)
}

func (l *Less) Evaluate() bool {
	return l.first.Evaluate() < l.second.Evaluate()
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

type LessEq struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewLessEq(first, second Evaluable[Numeric]) *LessEq {
	return &LessEq{NewPairForm(first, second, LessEqOperator)}
}

func (le *LessEq) TrueCopy() *LessEq {
	return &LessEq{le.PairForm.TrueCopy()}
}

func (le *LessEq) Copy() Form {
	return le.TrueCopy()
}

func (le *LessEq) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc(le.PairForm, diff)
}

func (le *LessEq) Evaluate() bool {
	return le.first.Evaluate() <= le.second.Evaluate()
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

type Great struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewGreat(first, second Evaluable[Numeric]) *Great {
	return &Great{NewPairForm(first, second, GreatOperator)}
}

func (g *Great) TrueCopy() *Great {
	return &Great{g.PairForm.TrueCopy()}
}

func (g *Great) Copy() Form {
	return g.TrueCopy()
}

func (g *Great) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc(g.PairForm, diff)
}

func (g *Great) Evaluate() bool {
	return g.first.Evaluate() > g.second.Evaluate()
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

type GreatEq struct {
	*PairForm[Evaluable[Numeric], Evaluable[Numeric]]
}

func NewGreatEq(first, second Evaluable[Numeric]) *GreatEq {
	return &GreatEq{NewPairForm(first, second, GreatEqOperator)}
}

func (ge *GreatEq) TrueCopy() *GreatEq {
	return &GreatEq{ge.PairForm.TrueCopy()}
}

func (ge *GreatEq) Copy() Form {
	return ge.TrueCopy()
}

func (ge *GreatEq) getFactorMap() map[string]Numeric {
	return getFactorMapForFunc(ge.PairForm, diff)
}

func (ge *GreatEq) Evaluate() bool {
	return ge.first.Evaluate() >= ge.second.Evaluate()
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
