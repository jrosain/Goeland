package arithmetic

import (
	"math"
	"strconv"
	"strings"
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

func convertPred(old basictypes.Pred) (result Evaluable[bool], termMap map[string]basictypes.Term, success bool) {
	termMap = make(map[string]basictypes.Term)
	args := []Evaluable[Numeric]{}

	for _, term := range old.GetArgs() {
		form, terms := convertTermAndRegisterVariables(term)
		args = append(args, form)

		for _, term := range terms {
			termMap[varPrefix+term.ToMappedString(basictypes.DefaultMap, false)] = term
		}
	}

	return convertBooleanPred(old, args, termMap)
}

func convertBooleanPred(old basictypes.Pred, args []Evaluable[Numeric], termMap map[string]basictypes.Term) (Evaluable[bool], map[string]basictypes.Term, bool) {
	switch old.GetID().GetName() {
	case "$is_int":
		return NewIsInt(args[0]), termMap, true
	case "$is_rat":
		return NewIsRat(args[0]), termMap, true
	default:
		return convertComparaisonPred(old)
	}
}

func convertComparaisonPred(old basictypes.Pred) (result ComparaisonForm, termMap map[string]basictypes.Term, success bool) {
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
		return NewEq(args[0], args[1]), termMap, true
	case "!=":
		return NewNotEq(args[0], args[1]), termMap, true
	case "$lesseq":
		return NewLessEq(args[0], args[1]), termMap, true
	case "$less":
		return NewLess(args[0], args[1]), termMap, true
	case "$greatereq":
		return NewGreatEq(args[0], args[1]), termMap, true
	case "$greater":
		return NewGreat(args[0], args[1]), termMap, true
	default:
		return nil, termMap, false
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
		value, success := getNumericForm(name)
		if success {
			setToConstantMap(value, old)
			return NewConstant(value), terms
		} else {
			terms = append(terms, old)
			return NewFactor(One, NewVariable(old.ToMappedString(basictypes.DefaultMap, false))), terms
		}
	}
}

func getNumericForm(str string) (result Numeric, success bool) {
	switch {
	case strings.Contains(str, "/"):
		return getRationalForm(str)
	case strings.Contains(str, "e"):
		return getExponentForm(str, "e")
	case strings.Contains(str, "E"):
		return getExponentForm(str, "E")
	case strings.Contains(str, "."):
		return getRealForm(str)
	default:
		return getIntegerForm(str)
	}
}

func getIntegerForm(str string) (result Integer, success bool) {
	if str[0] == '+' {
		str = str[1:]
	}

	res, err := strconv.Atoi(str)
	return Integer(res), err == nil
}

func getRationalForm(str string) (result Rational, success bool) {
	parts := strings.Split(str, "/")
	if len(parts) != 2 {
		return Rational{}, false
	}

	res1, err1 := getIntegerForm(parts[0])
	res2, err2 := getIntegerForm(parts[1])

	return Rational{int(res1), int(res2)}, err1 && err2
}

func getRealForm(str string) (result Real, success bool) {
	if str[0] == '+' {
		str = str[1:]
	}

	res, err := strconv.ParseFloat(str, 64)
	return Real(res), err == nil
}

func getExponentForm(str string, expSymbol string) (result Real, success bool) {
	parts := strings.Split(str, expSymbol)
	if len(parts) != 2 {
		return Real(0), false
	}

	res1, err1 := getRealForm(parts[0])
	res2, err2 := getIntegerForm(parts[1])

	base := 1.0
	for i := 0; i < int(math.Abs(res2.Evaluate())); i++ {
		base *= 10
	}
	if res2 < 0 {
		base = 1 / base
	}

	return Real(res1.Evaluate() * base), err1 && err2
}
