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

type IsInt struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewIsInt(value Evaluable[Numeric]) *IsInt {
	return &IsInt{NewSimpleForm(value)}
}

func (ii *IsInt) Copy() Form {
	return &IsInt{ii.SimpleForm.TrueCopy()}
}

func (ii *IsInt) ToString() string {
	return "isInt(" + ii.SimpleForm.ToString() + ")"
}

func (ii *IsInt) getFactorMap() map[string]float64 {
	global.PrintPanic("ARI", "Should not try to get the factor map of a IsInt formula as it does not make any sense.")
	return nil
}

func (ii *IsInt) Evaluate() bool {
	return ii.value.Evaluate().IsInt()
}

type IsRat struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewIsRat(value Evaluable[Numeric]) *IsRat {
	return &IsRat{NewSimpleForm(value)}
}

func (ir *IsRat) Copy() Form {
	return &IsRat{ir.SimpleForm.TrueCopy()}
}

func (ir *IsRat) ToString() string {
	return "isRat(" + ir.SimpleForm.ToString() + ")"
}

func (ir *IsRat) getFactorMap() map[string]float64 {
	global.PrintPanic("ARI", "Should not try to get the factor map of a IsRat formula as it does not make any sense.")
	return nil
}

func (ir *IsRat) Evaluate() bool {
	return ir.value.Evaluate().IsRat()
}

type ToInt struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewToInt(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &ToInt{NewSimpleForm(value)}
}

func (ti *ToInt) Copy() Form {
	return &ToInt{ti.SimpleForm.TrueCopy()}
}

func (ti *ToInt) ToString() string {
	return "toInt(" + ti.SimpleForm.ToString() + ")"
}

func (ti *ToInt) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = ti.Evaluate().ToInt().Evaluate()

	return factorMap
}

func (ti *ToInt) Evaluate() Numeric {
	return ti.value.Evaluate().ToInt()
}

type ToRat struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewToRat(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &ToRat{NewSimpleForm(value)}
}

func (tr *ToRat) Copy() Form {
	return &ToRat{tr.SimpleForm.TrueCopy()}
}

func (tr *ToRat) ToString() string {
	return "toRat(" + tr.SimpleForm.ToString() + ")"
}

func (tr *ToRat) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = tr.Evaluate().ToRat().Evaluate()

	return factorMap
}

func (tr *ToRat) Evaluate() Numeric {
	return tr.value.Evaluate().ToRat()
}

type ToReal struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewToReal(value Evaluable[Numeric]) Evaluable[Numeric] {
	return &ToReal{NewSimpleForm(value)}
}

func (tr *ToReal) Copy() Form {
	return &ToReal{tr.SimpleForm.TrueCopy()}
}

func (tr *ToReal) ToString() string {
	return "toReal(" + tr.SimpleForm.ToString() + ")"
}

func (tr *ToReal) getFactorMap() map[string]float64 {
	factorMap := make(map[string]float64)

	factorMap[Unit.ToString()] = tr.Evaluate().ToReal().Evaluate()

	return factorMap
}

func (tr *ToReal) Evaluate() Numeric {
	return tr.value.Evaluate().ToReal()
}
