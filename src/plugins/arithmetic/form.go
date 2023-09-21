package arithmetic

import (
	"github.com/GoelandProver/Goeland/global"
)

type Form interface {
	global.Basic[Form]
	getFactorMap() map[string]float64
}

type SimpleForm[T Form] struct {
	value T
}

type SimpleFormable[T Form] interface {
	GetValue() T
}

func NewSimpleForm[T Form](value T) *SimpleForm[T] {
	return &SimpleForm[T]{value}
}

func (sf *SimpleForm[T]) GetValue() T {
	return sf.value
}

func (sf *SimpleForm[T]) ToString() string {
	return sf.value.ToString()
}

func (sf *SimpleForm[T]) Equals(other any) bool {
	if typed, ok := other.(SimpleFormable[T]); ok {
		return sf.value.Equals(typed.GetValue())
	}
	return false
}

func (sf *SimpleForm[T]) Copy() Form {
	return sf.value.Copy()
}

func (sf *SimpleForm[T]) TrueCopy() *SimpleForm[T] {
	if typed, ok := sf.value.Copy().(T); ok {
		return NewSimpleForm[T](typed)
	}
	return nil
}

func (sf *SimpleForm[T]) getFactorMap() map[string]float64 {
	return sf.value.getFactorMap()
}

type AnyConstant interface {
	Evaluable[Numeric]
	IsConstant() bool
}

type Constant[T Numeric] struct {
	*SimpleForm[T]
}

var Zero *Constant[Integer] = NewConstant[Integer](0)
var One *Constant[Integer] = NewConstant[Integer](1)

func NewConstant[T Numeric](value T) *Constant[T] {
	return &Constant[T]{NewSimpleForm(value)}
}

func (c *Constant[T]) Copy() Form {
	return &Constant[T]{c.SimpleForm.TrueCopy()}
}

func (c *Constant[T]) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = c.Evaluate().Evaluate()

	return factorMap
}

func (c *Constant[T]) Evaluate() Numeric {
	return c.value
}

func (c *Constant[T]) IsConstant() bool {
	return true
}

const varPrefix string = ""

type Variable struct {
	*SimpleForm[String]
}

var Unit *Variable = NewVariable("")

func NewVariable(value string) *Variable {
	return &Variable{NewSimpleForm(String(varPrefix + value))}
}

func (v *Variable) Copy() Form {
	return &Variable{v.SimpleForm.TrueCopy()}
}

func (v *Variable) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[v.value.ToString()] = 1

	return factorMap
}

func (v *Variable) Evaluate() Numeric {
	global.PrintPanic("ARI", "Trying to evaluate a Variable, this should never happen")
	return Zero.value
}

type Neg struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewNeg(value Evaluable[Numeric]) Evaluable[Numeric] {
	switch typed := value.(type) {
	case *Neg:
		return typed.value
	case *Constant[Integer]:
		return NewConstant(-typed.value)
	case *Constant[Rational]:
		return NewConstant(typed.value.Negate())
	case *Constant[Real]:
		return NewConstant(-typed.value)
	default:
		return &Neg{NewSimpleForm(value)}
	}
}

func (n *Neg) Copy() Form {
	return &Neg{n.SimpleForm.TrueCopy()}
}

func (n *Neg) ToString() string {
	return "-" + n.SimpleForm.ToString()
}

func (n *Neg) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)
	childMap := n.value.getFactorMap()

	for k, v := range childMap {
		factorMap[k] = -v
	}

	return factorMap
}

func (n *Neg) Evaluate() Numeric {
	return n.value.Evaluate().Neg()
}

type Floor struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewFloor(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &Floor{NewSimpleForm(value)}
}

func (f *Floor) Copy() Form {
	return &Floor{f.SimpleForm.TrueCopy()}
}

func (f *Floor) ToString() string {
	return "⌊" + f.SimpleForm.ToString() + "⌋"
}

func (f *Floor) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = f.Evaluate().Evaluate()

	return factorMap
}

func (f *Floor) Evaluate() Numeric {
	return f.SimpleForm.value.Evaluate().Floor()
}

type Ceil struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewCeil(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &Ceil{NewSimpleForm(value)}
}

func (c *Ceil) Copy() Form {
	return &Ceil{c.SimpleForm.TrueCopy()}
}

func (c *Ceil) ToString() string {
	return "⌈" + c.SimpleForm.ToString() + "⌉"
}

func (c *Ceil) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = c.Evaluate().Evaluate()

	return factorMap
}

func (c *Ceil) Evaluate() Numeric {
	return c.value.Evaluate().Ceil()
}

type Trunc struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewTrunc(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &Trunc{NewSimpleForm(value)}
}

func (t *Trunc) Copy() Form {
	return &Trunc{t.SimpleForm.TrueCopy()}
}

func (t *Trunc) ToString() string {
	return "trunc(" + t.SimpleForm.ToString() + ")"
}

func (t *Trunc) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = t.Evaluate().Evaluate()

	return factorMap
}

func (t *Trunc) Evaluate() Numeric {
	return t.value.Evaluate().Trunc()
}

type Round struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewRound(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &Round{NewSimpleForm(value)}
}

func (r *Round) Copy() Form {
	return &Round{r.SimpleForm.TrueCopy()}
}

func (r *Round) ToString() string {
	return "round(" + r.SimpleForm.ToString() + ")"
}

func (r *Round) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = r.Evaluate().Evaluate()

	return factorMap
}

func (r *Round) Evaluate() Numeric {
	return r.value.Evaluate().Round()
}
