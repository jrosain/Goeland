package arithmetic

import (
	"strconv"

	"github.com/GoelandProver/Goeland/global"
)

type Form interface {
	global.Basic[Form]
}

type Integer int

func (i Integer) ToString() string {
	return strconv.Itoa(int(i))
}

func (i Integer) Equals(other any) bool {
	if typed, ok := other.(Integer); ok {
		return i == typed
	}
	return false
}

func (i Integer) Copy() Form {
	return i
}

type String string

func (s String) ToString() string {
	return string(s)
}

func (s String) Equals(other any) bool {
	if typed, ok := other.(String); ok {
		return s == typed
	}
	return false
}

func (s String) Copy() Form {
	return s
}

type SimpleForm[T global.Basic[Form]] struct {
	value T
}

func NewSimpleForm[T global.Basic[Form]](value T) *SimpleForm[T] {
	return &SimpleForm[T]{value}
}

func (sf *SimpleForm[T]) ToString() string {
	return sf.value.ToString()
}

func (sf *SimpleForm[T]) Equals(other any) bool {
	if typed, ok := other.(*SimpleForm[T]); ok {
		return sf.value.Equals(typed.value)
	}
	return false
}

func (sf *SimpleForm[T]) Copy() Form {
	return sf.value.Copy()
}

type Constant struct {
	*SimpleForm[Integer]
}

func NewConstant(value int) *Constant {
	return &Constant{NewSimpleForm(Integer(value))}
}

type Variable struct {
	*SimpleForm[String]
}

func NewVariable(value string) *Variable {
	return &Variable{NewSimpleForm(String(value))}
}

type PairOperator string

const (
	EqOperator      PairOperator = "="
	LessOperator    PairOperator = "<"
	LessEqOperator  PairOperator = "<="
	GreatOperator   PairOperator = ">"
	GreatEqOperator PairOperator = ">="
	SumOperator     PairOperator = "+"
	DiffOperator    PairOperator = "-"
	ProdOperator    PairOperator = "*"
)

type PairForm struct {
	first  Form
	second Form
	symbol PairOperator
}

func NewPairForm(first, second Form, symbol PairOperator) *PairForm {
	return &PairForm{first, second, symbol}
}

func (pf *PairForm) ToString() string {
	return "(" + pf.first.ToString() + " " + string(pf.symbol) + " " + pf.second.ToString() + ")"
}

func (pf *PairForm) Equals(other any) bool {
	if typed, ok := other.(*PairForm); ok {
		return pf.first.Equals(typed.first) && pf.second.Equals(typed.second)
	}
	return false
}

func (pf *PairForm) Copy() Form {
	return pf.TrueCopy()
}

func (pf *PairForm) TrueCopy() *PairForm {
	return NewPairForm(pf.first.Copy(), pf.second.Copy(), pf.symbol)
}

type Eq struct {
	*PairForm
}

func NewEq(first, second Form) *Eq {
	return &Eq{NewPairForm(first, second, EqOperator)}
}

func (e *Eq) Copy() Form {
	return &Eq{e.PairForm.TrueCopy()}
}

type Less struct {
	*PairForm
}

func NewLess(first, second Form) *Less {
	return &Less{NewPairForm(first, second, LessOperator)}
}

func (l *Less) Copy() Form {
	return &Less{l.PairForm.TrueCopy()}
}

type LessEq struct {
	*PairForm
}

func NewLessEq(first, second Form) *LessEq {
	return &LessEq{NewPairForm(first, second, LessEqOperator)}
}

func (le *LessEq) Copy() Form {
	return &LessEq{le.PairForm.TrueCopy()}
}

type Great struct {
	*PairForm
}

func NewGreat(first, second Form) *Great {
	return &Great{NewPairForm(first, second, GreatOperator)}
}

func (g *Great) Copy() Form {
	return &Great{g.PairForm.TrueCopy()}
}

type GreatEq struct {
	*PairForm
}

func NewGreatEq(first, second Form) *GreatEq {
	return &GreatEq{NewPairForm(first, second, GreatEqOperator)}
}

func (ge *GreatEq) Copy() Form {
	return &GreatEq{ge.PairForm.TrueCopy()}
}

type Sum struct {
	*PairForm
}

func NewSum(first, second Form) *Sum {
	return &Sum{NewPairForm(first, second, SumOperator)}
}

func (s *Sum) Copy() Form {
	return &Sum{s.PairForm.TrueCopy()}
}

type Diff struct {
	*PairForm
}

func NewDiff(first, second Form) *Diff {
	return &Diff{NewPairForm(first, second, DiffOperator)}
}

func (d *Diff) Copy() Form {
	return &Diff{d.PairForm.TrueCopy()}
}

type Prod struct {
	*PairForm
}

func NewProd(first, second Form) *Prod {
	return &Prod{NewPairForm(first, second, ProdOperator)}
}

func (p *Prod) Copy() Form {
	return &Prod{p.PairForm.TrueCopy()}
}
