package arithmetic

import (
	"strconv"

	"github.com/GoelandProver/Goeland/global"
)

type Form interface {
	global.Basic[Form]
}

type Constant struct {
	value int
}

type Variable struct {
	value string
}

func (c *Constant) ToString() string {
	return strconv.Itoa(c.value)
}

func (c *Constant) Equals(other any) bool {
	if typed, ok := other.(*Constant); ok {
		return c.value == typed.value
	}
	return false
}

func (c *Constant) Copy() Form {
	return &Constant{c.value}
}

func (v *Variable) ToString() string {
	return v.value
}

func (v *Variable) Equals(other any) bool {
	if typed, ok := other.(*Variable); ok {
		return v.value == typed.value
	}
	return false
}

func (v *Variable) Copy() Form {
	return &Variable{v.value}
}
