package search

import (
	"fmt"

	"github.com/GoelandProver/Goeland/global"
	typing "github.com/GoelandProver/Goeland/polymorphism/typing"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
	complextypes "github.com/GoelandProver/Goeland/types/complex-types"
	visualization "github.com/GoelandProver/Goeland/visualization_exchanges"
)

var ZeqEnabled = false

func EnableZeq() {
	ds.doCorrectApplyRules = ds.zeqApplyRule
}

func equalitySort(fatherId uint64, state complextypes.State, c Communication, newAtomics basictypes.FormAndTermsList, currentNodeId int, originalNodeId int, metaToReintroduce []int) (eqs, neqs basictypes.FormAndTermsList) {
	atoms := state.GetAtomic()
	neqs = basictypes.MakeEmptyFormAndTermsList()
	eqs = basictypes.MakeEmptyFormAndTermsList()

	for _, elem := range atoms {
		switch formTyped := elem.GetForm().(type) {
		case basictypes.Not:
			if typed, ok := formTyped.GetForm().(basictypes.Pred); ok && typed.GetID().Equals(basictypes.Id_eq) {
				neqs = neqs.AppendIfNotContains(elem)
			}
		case basictypes.Pred:
			if formTyped.GetID().Equals(basictypes.Id_eq) {
				eqs = eqs.AppendIfNotContains(elem)
			}
		}
	}

	return eqs, neqs
}

func (ds *destructiveSearch) zeqApplyRule(fatherId uint64, state complextypes.State, c Communication, newAtomics basictypes.FormAndTermsList, currentNodeId int, originalNodeId int, metaToReintroduce []int) {

	eqs, neqs := equalitySort(fatherId, state, c, newAtomics, currentNodeId, originalNodeId, metaToReintroduce)
	state.SetEqs(eqs.ExtractForms())
	state.SetNeqs(neqs.ExtractForms())
	eqPair := CanApplyTs(state)
	predPair := CanApplyPred(state)

	global.PrintDebug("PS", fmt.Sprintf("Equations : %d, Inequations : %d", eqs.Len(), neqs.Len()))

	global.PrintDebug("AR", "ApplyRule")
	switch {
	case len(newAtomics) > 0 && global.IsLoaded("dmt") && len(state.GetSubstsFound()) == 0:
		ds.manageRewriteRules(fatherId, state, c, newAtomics, currentNodeId, originalNodeId, metaToReintroduce)

	case len(state.GetAlpha()) > 0:
		ds.manageAlphaRules(fatherId, state, c, originalNodeId)

	case len(state.GetDelta()) > 0:
		ds.manageDeltaRules(fatherId, state, c, originalNodeId)

	case !isNilPair(predPair):
		ds.applyPredRules(fatherId, state, c, originalNodeId, currentNodeId, metaToReintroduce, predPair)

	case !isNilPair(eqPair):
		ds.applyTsRules(fatherId, state, c, originalNodeId, currentNodeId, metaToReintroduce, eqPair)

	case len(state.GetBeta()) > 0:
		ds.manageBetaRules(fatherId, state, c, currentNodeId, originalNodeId, metaToReintroduce)

	case len(state.GetGamma()) > 0:
		ds.manageGammaRules(fatherId, state, c, originalNodeId)

	case len(state.GetMetaGen()) > 0 && state.CanReintroduce():
		ds.manageReintroductionRules(fatherId, state, c, originalNodeId, metaToReintroduce, newAtomics, currentNodeId, true)

	default:
		visualization.WriteExchanges(fatherId, state, nil, complextypes.MakeEmptySubstAndForm(), "ApplyRules - SAT")
		state.SetCurrentProofRule("Sat")
		state.SetProof(append(state.GetProof(), state.GetCurrentProof()))
		global.PrintDebug("PS", "Nothing found, return sat")
		ds.sendSubToFather(c, false, false, fatherId, state, []complextypes.SubstAndForm{}, currentNodeId, originalNodeId, []int{})
	}
}

func (ds *destructiveSearch) applyTsRules(fatherId uint64, state complextypes.State, c Communication, originalNodeId int, currentNodeId int, metaToReintroduce []int, pair global.BasicPaired[basictypes.Form, basictypes.Form]) {
	global.PrintDebug("PS", "Zeq rule")
	hdfEq := pair.GetFst()
	hdfNeq := pair.GetSnd()
	global.PrintDebug("PS", fmt.Sprintf("Rule applied on : %s %s", hdfEq.ToString(), hdfNeq.ToString()))

	s, t := hdfEq.(basictypes.Pred).GetArgs().Get(0), hdfEq.(basictypes.Pred).GetArgs().Get(1)
	u, v := hdfNeq.(basictypes.Not).GetForm().(basictypes.Pred).GetArgs().Get(0), hdfNeq.(basictypes.Not).GetForm().(basictypes.Pred).GetArgs().Get(1)

	global.PrintDebug("PS", fmt.Sprintf("Found litterals : s = %s t = %s, u = %s, v = %s", s.ToString(), t.ToString(), u.ToString(), v.ToString()))

	vneqs := basictypes.RefuteForm(basictypes.MakerPred(
		basictypes.Id_eq,
		basictypes.NewTermList(v, s),
		[]typing.TypeApp{},
	))

	tnequ := basictypes.RefuteForm(basictypes.MakerPred(
		basictypes.Id_eq,
		basictypes.NewTermList(t, u),
		[]typing.TypeApp{},
	))

	global.PrintDebug("PS", fmt.Sprintf("Generated formulas : %s, %s", vneqs.ToString(), tnequ.ToString()))

	state.AddToAlreadyAppliedZeq(pair)

	var formTs [2]basictypes.Form
	formTs[0] = vneqs
	formTs[1] = tnequ

	childIds := []int{}
	var channels []Communication

	for _, elem := range formTs {
		i := global.IncrCptNode()

		otherState := state.Copy()
		otherFl := basictypes.MakeSingleElementFormAndTermList(basictypes.MakeFormAndTerm(elem, basictypes.NewTermList()))

		otherState.SetLF(otherFl)
		childIds = append(childIds, i)

		if global.IsDestructive() {
			channelChild := Communication{make(chan bool), make(chan Result)}
			channels = append(channels, channelChild)
			go ds.ProofSearch(global.GetGID(), otherState, channelChild, complextypes.MakeEmptySubstAndForm(), i, i, []int{})
		} else {
			go ds.ProofSearch(global.GetGID(), otherState, c, complextypes.MakeEmptySubstAndForm(), i, i, []int{})
		}

		global.IncrGoRoutine(1)
		global.PrintDebug("PS", fmt.Sprintf("GO %v !", i))
	}
	ds.DoEndApplyZeq(fatherId, state, c, channels, currentNodeId, originalNodeId, childIds, metaToReintroduce)
	// ds.ProofSearch(fatherId, state, c, complextypes.MakeEmptySubstAndForm(), childId, originalNodeId, []int{})
}

func (ds *destructiveSearch) applyPredRules(fatherId uint64, state complextypes.State, c Communication, originalNodeId int, currentNodeId int, metaToReintroduce []int, pair global.BasicPaired[basictypes.Form, basictypes.Form]) {
	global.PrintDebug("PS", "Zeq rule")
	hdfEq := pair.GetFst()
	hdfNeq := pair.GetSnd()
	global.PrintDebug("PS", fmt.Sprintf("Rule applied on : %s %s", hdfEq.ToString(), hdfNeq.ToString()))

	s, t := hdfEq.(basictypes.Pred).GetArgs().Get(0), hdfEq.(basictypes.Pred).GetArgs().Get(1)
	u, v := hdfNeq.(basictypes.Not).GetForm().(basictypes.Pred).GetArgs().Get(0), hdfNeq.(basictypes.Not).GetForm().(basictypes.Pred).GetArgs().Get(1)

	global.PrintDebug("PS", fmt.Sprintf("Found litterals : s = %s t = %s, u = %s, v = %s", s.ToString(), t.ToString(), u.ToString(), v.ToString()))

	vneqs := basictypes.RefuteForm(basictypes.MakerPred(
		basictypes.Id_eq,
		basictypes.NewTermList(v, s),
		[]typing.TypeApp{},
	))

	tnequ := basictypes.RefuteForm(basictypes.MakerPred(
		basictypes.Id_eq,
		basictypes.NewTermList(t, u),
		[]typing.TypeApp{},
	))

	global.PrintDebug("PS", fmt.Sprintf("Generated formulas : %s, %s", vneqs.ToString(), tnequ.ToString()))

	state.AddToAlreadyAppliedZeq(pair)

	var formTs [2]basictypes.Form
	formTs[0] = vneqs
	formTs[1] = tnequ

	childIds := []int{}
	var channels []Communication

	for _, elem := range formTs {
		i := global.IncrCptNode()

		otherState := state.Copy()
		otherFl := basictypes.MakeSingleElementFormAndTermList(basictypes.MakeFormAndTerm(elem, basictypes.NewTermList()))

		otherState.SetLF(otherFl)
		childIds = append(childIds, i)

		if global.IsDestructive() {
			channelChild := Communication{make(chan bool), make(chan Result)}
			channels = append(channels, channelChild)
			go ds.ProofSearch(global.GetGID(), otherState, channelChild, complextypes.MakeEmptySubstAndForm(), i, i, []int{})
		} else {
			go ds.ProofSearch(global.GetGID(), otherState, c, complextypes.MakeEmptySubstAndForm(), i, i, []int{})
		}

		global.IncrGoRoutine(1)
		global.PrintDebug("PS", fmt.Sprintf("GO %v !", i))
	}
	ds.DoEndApplyZeq(fatherId, state, c, channels, currentNodeId, originalNodeId, childIds, metaToReintroduce)
	// ds.ProofSearch(fatherId, state, c, complextypes.MakeEmptySubstAndForm(), childId, originalNodeId, []int{})
}

func (ds *destructiveSearch) DoEndApplyZeq(fatherId uint64, state complextypes.State, c Communication, channels []Communication, currentNodeId int, originalNodeId int, childIds []int, metaToReintroduce []int) {
	ds.waitChildren(MakeWcdArgs(fatherId, state, c, channels, []complextypes.SubstAndForm{}, complextypes.MakeEmptySubstAndForm(), []complextypes.SubstAndForm{}, []complextypes.IntSubstAndFormAndTerms{}, currentNodeId, originalNodeId, false, childIds, metaToReintroduce))
}

func CanApplyTs(state complextypes.State) global.BasicPaired[basictypes.Form, basictypes.Form] {
	for _, f1 := range state.GetEqs().Slice() {
		for _, f2 := range state.GetNeqs().Slice() {
			pair := global.NewBasicPair[basictypes.Form, basictypes.Form](f1, f2)
			if !isAlreadyApplied(state, pair) {
				return pair
			}
		}
	}
	return global.NewBasicPair[basictypes.Form, basictypes.Form](nil, nil)
}

func CanApplyPred(state complextypes.State) global.BasicPaired[basictypes.Form, basictypes.Form] {
	predCandidateList := getPosAndNegPreds(state)
	for _, predPair := range predCandidateList.Slice() {
		pair := global.NewBasicPair[basictypes.Form, basictypes.Form](predPair.GetFst(), predPair.GetSnd())
		if !isAlreadyApplied(state, pair) {
			return pair
		}
	}
	return global.NewBasicPair[basictypes.Form, basictypes.Form](nil, nil)
}

func isAlreadyApplied(state complextypes.State, pair global.BasicPaired[basictypes.Form, basictypes.Form]) bool {
	for _, each := range state.GetAlreadyAppliedZeq().Slice() {
		if each.GetFst().Equals(pair.GetFst()) && each.GetSnd().Equals(pair.GetSnd()) {
			return true
		}
	}
	return false
}

func isNilPair(pair global.BasicPaired[basictypes.Form, basictypes.Form]) bool {
	return pair.GetFst() == nil && pair.GetSnd() == nil
}

func getPosAndNegPreds(st complextypes.State) *global.List[global.BasicPaired[basictypes.Pred, basictypes.Pred]] {
	result := global.NewList[global.BasicPaired[basictypes.Pred, basictypes.Pred]]()
	atomics := st.GetAtomic().ExtractForms()

	for i := 0; i < atomics.Len()-1; i++ {
		if typedFirst, ok := atomics.Get(i).(basictypes.Pred); ok {
			for j := 1; j < atomics.Len(); j++ {
				if not, ok := atomics.Get(j).(basictypes.Not); ok {
					if typedSecond, ok := not.GetForm().(basictypes.Pred); ok {
						result.AppendIfNotContains(global.NewBasicPair(typedFirst, typedSecond))
					}
				}
			}
		}
	}

	return result
}
