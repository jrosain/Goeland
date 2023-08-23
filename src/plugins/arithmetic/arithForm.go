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

type PairForm struct {
	first  Form
	second Form
}

func NewPairForm(first Form, second Form) *PairForm {
	return &PairForm{first, second}
}

func (pf *PairForm) ToString() string {
	return pf.first.ToString() + " " + pf.second.ToString()
}

func (pf *PairForm) Equals(other any) bool {
	if typed, ok := other.(*PairForm); ok {
		return pf.first.Equals(typed.first) && pf.second.Equals(typed.second)
	}
	return false
}

func (pf *PairForm) Copy() Form {
	return NewPairForm(pf.first.Copy(), pf.second.Copy())
}

type Less struct {
	*PairForm
}

func NewLess(first Form, second Form) *Less {
	return &Less{NewPairForm(first, second)}
}

func (l *Less) ToString() string {
	return "(" + l.first.ToString() + " < " + l.second.ToString() + ")"
}

type LessEq struct {
	*PairForm
}

func NewLessEq(first Form, second Form) *LessEq {
	return &LessEq{NewPairForm(first, second)}
}

func (le *LessEq) ToString() string {
	return "(" + le.first.ToString() + " <= " + le.second.ToString() + ")"
}

type Great struct {
	*PairForm
}

func NewGreat(first Form, second Form) *Great {
	return &Great{NewPairForm(first, second)}
}

func (g *Great) ToString() string {
	return "(" + g.first.ToString() + " > " + g.second.ToString() + ")"
}

type GreatEq struct {
	*PairForm
}

func NewGreatEq(first Form, second Form) *GreatEq {
	return &GreatEq{NewPairForm(first, second)}
}

func (ge *GreatEq) ToString() string {
	return "(" + ge.first.ToString() + " >= " + ge.second.ToString() + ")"
}

type Sum struct {
	*PairForm
}

func NewSum(first Form, second Form) *Sum {
	return &Sum{NewPairForm(first, second)}
}

func (ge *Sum) ToString() string {
	return "(" + ge.first.ToString() + " + " + ge.second.ToString() + ")"
}

type Difference struct {
	*PairForm
}

func NewDifference(first Form, second Form) *Difference {
	return &Difference{NewPairForm(first, second)}
}

func (ge *Difference) ToString() string {
	return "(" + ge.first.ToString() + " - " + ge.second.ToString() + ")"
}

type Product struct {
	*PairForm
}

func NewProduct(first Form, second Form) *Product {
	return &Product{NewPairForm(first, second)}
}

func (ge *Product) ToString() string {
	return "(" + ge.first.ToString() + " * " + ge.second.ToString() + ")"
}
