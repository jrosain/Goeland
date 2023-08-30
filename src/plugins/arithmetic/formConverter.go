package arithmetic

import (
	"strconv"

	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

func convertPred(old basictypes.Pred) (result ComparisonForm, termMap map[string]basictypes.Term) {
	termMap = make(map[string]basictypes.Term)
	args := []Form{}

	for _, term := range old.GetArgs() {
		form, isVariable := convertTerm(term)
		args = append(args, form)

		if isVariable {
			termMap[term.ToMappedString(basictypes.DefaultMap, false)] = term
		}
	}

	switch old.GetID().GetName() {
	case "=":
		return NewEq(args[0], args[1]), termMap
	case "$lesseq":
		return NewLessEq(args[0], args[1]), termMap
	case "$less":
		return NewLess(args[0], args[1]), termMap
	case "$greatereq":
		return NewGreatEq(args[0], args[1]), termMap
	case "$greater":
		return NewGreat(args[0], args[1]), termMap
	default:
		return nil, termMap
	}
}

func convertTerm(old basictypes.Term) (result Form, isVariable bool) {
	value, err := strconv.Atoi(old.GetName())
	if err == nil {
		return NewConstant(value), false
	} else {
		return NewFactor(One, NewVariable(old.ToMappedString(basictypes.DefaultMap, false))), true
	}
}
