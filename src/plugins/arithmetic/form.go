package arithmetic

import (
	"fmt"
	"math"

	"github.com/GoelandProver/Goeland/global"
)

type Form interface {
	global.Basic[Form]
	getFactorMap() map[string]Numeric
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

func (sf *SimpleForm[T]) getFactorMap() map[string]Numeric {
	return sf.value.getFactorMap()
}

type Constant struct {
	*SimpleForm[Numeric]
}

var Zero *Constant = NewConstant(0)
var One *Constant = NewConstant(1)

func NewConstant(value Numeric) *Constant {
	return &Constant{NewSimpleForm(value)}
}

func (c *Constant) Copy() Form {
	return &Constant{c.SimpleForm.TrueCopy()}
}

func (c *Constant) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)

	factorMap[Unit.ToString()] = Numeric(c.value)

	return factorMap
}

func (c *Constant) Evaluate() Numeric {
	return c.value.Evaluate()
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

func (v *Variable) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)

	factorMap[v.value.ToString()] = 1

	return factorMap
}

func (v *Variable) Evaluate() Numeric {
	global.PrintPanic("ARI", "Trying to evaluate a Variable, this should never happen")
	return 0
}

type Neg struct {
	*SimpleForm[Evaluable[Numeric]]
}

func NewNeg(value Evaluable[Numeric]) Evaluable[Numeric] {
	switch typed := value.(type) {
	case *Neg:
		return typed.value
	case *Constant:
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

func (n *Neg) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)
	childMap := n.value.getFactorMap()

	for k, v := range childMap {
		factorMap[k] = -v
	}

	return factorMap
}

func (n *Neg) Evaluate() Numeric {
	return -n.value.Evaluate()
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

func (f *Floor) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)

	factorMap[Unit.ToString()] = f.Evaluate()
	fmt.Printf("--%v\n", factorMap)

	return factorMap
}

func (f *Floor) Evaluate() Numeric {
	return Numeric(math.Floor(float64(f.SimpleForm.value.Evaluate())))
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

func (c *Ceil) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)

	factorMap[Unit.ToString()] = c.Evaluate()

	return factorMap
}

func (c *Ceil) Evaluate() Numeric {
	return Numeric(math.Ceil(float64(c.SimpleForm.value.Evaluate())))
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

func (t *Trunc) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)

	factorMap[Unit.ToString()] = t.Evaluate()

	return factorMap
}

func (t *Trunc) Evaluate() Numeric {
	return Numeric(math.Trunc(float64(t.SimpleForm.value.Evaluate())))
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

func (r *Round) getFactorMap() map[string]Numeric {
	factorMap := make(map[string]Numeric)

	factorMap[Unit.ToString()] = r.Evaluate()

	return factorMap
}

func (r *Round) Evaluate() Numeric {
	return Numeric(math.RoundToEven(float64(r.SimpleForm.value.Evaluate())))
}
