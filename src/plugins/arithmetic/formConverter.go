package arithmetic

import (
	"strconv"
	"sync"

	typing "github.com/GoelandProver/Goeland/polymorphism/typing"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

var constantMap map[Numeric]basictypes.Term = make(map[Numeric]basictypes.Term)
var constantMapMutex sync.Mutex

func addToConstantMap(value Numeric) {
	constantMapMutex.Lock()
	defer constantMapMutex.Unlock()

	if _, ok := constantMap[value]; !ok {
		constantMap[value] = basictypes.MakerFun(basictypes.MakerId(value.ToString()), basictypes.MakeEmptyTermList(), []typing.TypeApp{}, typing.MkTypeHint("$int"))
	}
}

func setToConstantMap(key Numeric, value basictypes.Term) {
	constantMapMutex.Lock()
	defer constantMapMutex.Unlock()

	constantMap[key] = value
}

func getValueTerm(value Numeric) basictypes.Term {
	constantMapMutex.Lock()
	defer constantMapMutex.Unlock()

	result := constantMap[value]
	return result
}

func convertPred(old basictypes.Pred) (result Evaluable[bool], termMap map[string]basictypes.Term) {
	termMap = make(map[string]basictypes.Term)
	args := []Evaluable[Numeric]{}

	for _, term := range old.GetArgs() {
		form, terms := convertTermAndRegisterVariables(term)
		args = append(args, form)

		for _, term := range terms {
			termMap[varPrefix+term.ToMappedString(basictypes.DefaultMap, false)] = term
		}
	}

	switch old.GetID().GetName() {
	case "$is_int":
		return NewIsInt(args[0]), termMap
	case "$is_rat":
		return NewIsRat(args[0]), termMap
	default:
		return convertComparaisonPred(old)
	}
}

func convertComparaisonPred(old basictypes.Pred) (result ComparaisonForm, termMap map[string]basictypes.Term) {
	termMap = make(map[string]basictypes.Term)
	args := []Evaluable[Numeric]{}

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

func convertTermAndRegisterVariables(old basictypes.Term) (result Evaluable[Numeric], terms basictypes.TermList) {
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
	case "$floor":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewFloor(form), terms
		}

		return nil, terms
	case "$ceiling":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewCeil(form), terms
		}

		return nil, terms
	case "$truncate":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewTrunc(form), terms
		}

		return nil, terms
	case "$round":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewRound(form), terms
		}

		return nil, terms
	case "$quotient":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewQuotient(form1, form2), terms
		}

		return nil, terms
	case "$quotient_e":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewQuotientE(form1, form2), terms
		}

		return nil, terms
	case "$quotient_t":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewQuotientT(form1, form2), terms
		}

		return nil, terms
	case "$quotient_f":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewQuotientF(form1, form2), terms
		}

		return nil, terms
	case "$remainder_e":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewRemainderE(form1, form2), terms
		}

		return nil, terms
	case "$remainder_t":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewRemainderT(form1, form2), terms
		}

		return nil, terms
	case "$remainder_f":
		if typed, ok := old.(basictypes.Fun); ok {
			form1, newTerms1 := convertTermAndRegisterVariables(typed.GetArgs()[0])
			form2, newTerms2 := convertTermAndRegisterVariables(typed.GetArgs()[1])
			terms = append(terms, newTerms1...)
			terms = append(terms, newTerms2...)
			return NewRemainderF(form1, form2), terms
		}

		return nil, terms
	case "$to_int":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewToInt(form), terms
		}

		return nil, terms
	case "$to_rat":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewToRat(form), terms
		}

		return nil, terms
	case "$to_real":
		if typed, ok := old.(basictypes.Fun); ok {
			form, newTerms := convertTermAndRegisterVariables(typed.GetArgs()[0])
			terms = append(terms, newTerms...)
			return NewToReal(form), terms
		}

		return nil, terms
	default:
		value, err := strconv.Atoi(name)
		if err == nil {
			setToConstantMap(Integer(value), old)
			return NewConstant(Integer(value)), terms
		} else {
			terms = append(terms, old)
			return NewFactor(One, NewVariable(old.ToMappedString(basictypes.DefaultMap, false))), terms
		}
	}
}
