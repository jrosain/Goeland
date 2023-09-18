package arithmetic

import (
	"strconv"
	"sync"

	typing "github.com/GoelandProver/Goeland/polymorphism/typing"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

var constantMap map[int]basictypes.Term = make(map[int]basictypes.Term)
var constantMapMutex sync.Mutex

func addToConstantMap(value int) {
	constantMapMutex.Lock()
	defer constantMapMutex.Unlock()

	if _, ok := constantMap[value]; !ok {
		constantMap[value] = basictypes.MakerFun(basictypes.MakerId(strconv.Itoa(value)), basictypes.MakeEmptyTermList(), []typing.TypeApp{}, typing.MkTypeHint("$int"))
	}
}

func setToConstantMap(key int, value basictypes.Term) {
	constantMapMutex.Lock()
	defer constantMapMutex.Unlock()

	constantMap[key] = value
}

func getValueTerm(value int) basictypes.Term {
	constantMapMutex.Lock()
	defer constantMapMutex.Unlock()

	result := constantMap[value]
	return result
}

func convertPred(old basictypes.Pred) (result ComparisonForm, termMap map[string]basictypes.Term) {
	termMap = make(map[string]basictypes.Term)
	args := []Evaluable[int]{}

	for _, term := range old.GetArgs() {
		form, terms := convertTermAndRegisterVariables(term)
		args = append(args, form)

		for _, term := range terms {
			termMap[varPrefix+term.ToMappedString(basictypes.DefaultMap, false)] = term
		}
	}

	switch old.GetID().GetName() {
	case "=":
		return NewEq(args[0], args[1]), termMap
	case "!=":
		return NewNotEq(args[0], args[1]), termMap
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

func convertTermAndRegisterVariables(old basictypes.Term) (result Evaluable[int], terms basictypes.TermList) {
	terms = basictypes.TermList{}
	name := old.GetName()
	switch name {
	case "$sum":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewSum(form1, form2), terms
		}

		return nil, terms
	case "$difference":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewDiff(form1, form2), terms
		}

		return nil, terms
	case "$product":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewProduct(form1, form2), terms
		}

		return nil, terms
	case "$uminus":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewNeg(form), terms
		}

		return nil, terms
	default:
		value, err := strconv.Atoi(name)
		if err == nil {
			setToConstantMap(value, old)
			return NewConstant(value), terms
		} else {
			terms = append(terms, old)
			return NewFactor(One, NewVariable(old.ToMappedString(basictypes.DefaultMap, false))), terms
		}
	}
}
