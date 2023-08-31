package arithmetic

import (
	"strconv"

	"github.com/GoelandProver/Goeland/global"
)

type Form interface {
	global.Basic[Form]
	getFactorMap() map[string]int
}

type Evaluable[T any] interface {
	Form
	Evaluate() T
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

func (i Integer) getFactorMap() map[string]int {
	return make(map[string]int)
}

func (i Integer) Evaluate() int {
	return int(i)
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

func (s String) getFactorMap() map[string]int {
	return make(map[string]int)
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

func (sf *SimpleForm[T]) getFactorMap() map[string]int {
	return sf.value.getFactorMap()
}

type Constant struct {
	*SimpleForm[Integer]
}

var Zero *Constant = NewConstant(0)
var One *Constant = NewConstant(1)

func NewConstant(value int) *Constant {
	return &Constant{NewSimpleForm(Integer(value))}
}

func (c *Constant) getFactorMap() map[string]int {
	factorMap := make(map[string]int)

	factorMap[Unit.ToString()] = int(c.value)

	return factorMap
}

func (c *Constant) Evaluate() int {
	return c.value.Evaluate()
}

type Variable struct {
	*SimpleForm[String]
}

var Unit *Variable = NewVariable("")

func NewVariable(value string) *Variable {
	return &Variable{NewSimpleForm(String(value))}
}

func (v *Variable) getFactorMap() map[string]int {
	factorMap := make(map[string]int)

	factorMap[v.value.ToString()] = 1

	return factorMap
}

func (v *Variable) Evaluate() int {
	global.PrintPanic("ARI", "Trying to evaluate a Variable, this should never happen")
	return 0
}

type Neg struct {
	*SimpleForm[Evaluable[int]]
}

func NewNeg(value Evaluable[int]) Evaluable[int] {
	switch typed := value.(type) {
	case *Neg:
		return typed.value
	case *Constant:
		if typed.value < 0 {
			return NewConstant(-int(typed.value))
		} else {
			return &Neg{NewSimpleForm(value)}
		}
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

func (n *Neg) getFactorMap() map[string]int {
	factorMap := make(map[string]int)
	childMap := n.value.getFactorMap()

	for k, v := range childMap {
		factorMap[k] = -v
	}

	return factorMap
}

func (n *Neg) Evaluate() int {
	return -n.value.Evaluate()
}
