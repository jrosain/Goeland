package arithmetic

import (
	"fmt"

	"github.com/GoelandProver/Goeland/global"
)

type Comparaison string

const (
	Equals    Comparaison = "="
	GreaterEq Comparaison = ">="
	LesserEq  Comparaison = "<="
)

type Constraint interface {
	global.Basic[Constraint]
}

type Network global.List[Constraint]

func (n *Network) ToString() string {
	str := ""

	for _, element := range *n {
		str += element.ToString() + "\n"
	}

	return str
}

type SimpleConstraint struct {
	form   Form
	symbol Comparaison
	value  *Constant
}

func MakeSimpleConstraint(form Form, symbol Comparaison, value *Constant) SimpleConstraint {
	return SimpleConstraint{form, symbol, value}
}

func NewSimpleConstraint(form Form, symbol Comparaison, value *Constant) *SimpleConstraint {
	result := MakeSimpleConstraint(form, symbol, value)
	return &result
}

func (sc *SimpleConstraint) Equals(other any) bool {
	if typed, ok := other.(*SimpleConstraint); ok {
		return sc.form.Equals(typed.form) &&
			sc.symbol == typed.symbol &&
			sc.value.Equals(typed.value)
	}
	return false
}

func (sc *SimpleConstraint) Copy() Constraint {
	var typedConst *Constant
	if typed, ok := sc.value.Copy().(*Constant); ok {
		typedConst = typed
	}
	return NewSimpleConstraint(sc.form.Copy(), sc.symbol, typedConst)
}

func (sc *SimpleConstraint) ToString() string {
	return fmt.Sprintf("%s %s %v", sc.form.ToString(), sc.symbol, sc.value.ToString())
}
