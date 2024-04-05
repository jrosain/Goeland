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
	state.SetEqs(eqs)
	state.SetNeqs(neqs)
	pair := CanApplyTs(state)

	global.PrintDebug("PS", fmt.Sprintf("Equations : %d, Inequations : %d", eqs.Len(), neqs.Len()))

	global.PrintDebug("AR", "ApplyRule")
	switch {
	case len(newAtomics) > 0 && global.IsLoaded("dmt") && len(state.GetSubstsFound()) == 0:
		ds.manageRewriteRules(fatherId, state, c, newAtomics, currentNodeId, originalNodeId, metaToReintroduce)

	case len(state.GetAlpha()) > 0:
		ds.manageAlphaRules(fatherId, state, c, originalNodeId)

	case len(state.GetDelta()) > 0:
		ds.manageDeltaRules(fatherId, state, c, originalNodeId)

	case !isNilPair(pair):
		ds.applyZeqRules(fatherId, state, c, originalNodeId, pair)

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

func (ds *destructiveSearch) applyZeqRules(fatherId uint64, state complextypes.State, c Communication, originalNodeId int, pair global.Pair[int, int]) {
	global.PrintDebug("PS", "Zeq rule")
	hdfEq := state.GetEqs()[pair.Fst]
	hdfNeq := state.GetNeqs()[pair.Snd]
	global.PrintDebug("PS", fmt.Sprintf("Rule applied on : %s %s", hdfEq.ToString(), hdfNeq.ToString()))

	s, t := hdfEq.GetForm().(basictypes.Pred).GetArgs().Get(0), hdfEq.GetForm().(basictypes.Pred).GetArgs().Get(1)
	u, v := hdfNeq.GetForm().(basictypes.Not).GetForm().(basictypes.Pred).GetArgs().Get(0), hdfNeq.GetForm().(basictypes.Not).GetForm().(basictypes.Pred).GetArgs().Get(1)

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

	global.PrintDebug("PS", fmt.Sprintf("Found litterals : s = %s t = %s, u = %s, v = %s", s.ToString(), t.ToString(), u.ToString(), v.ToString()))
	atomicList := state.GetAtomic()
	fat := basictypes.MakeFormAndTerm(vneqs, atomicList[0].GetTerms())
	atomicList = atomicList.AppendIfNotContains(fat)
	fat = basictypes.MakeFormAndTerm(tnequ, atomicList[0].GetTerms())
	atomicList = atomicList.AppendIfNotContains(fat)

	global.PrintDebug("PS", fmt.Sprintf("New atomic formulae : %s", atomicList.ToString()))
	state.SetAtomic(atomicList)

	childId := global.IncrCptNode()

	ds.ProofSearch(fatherId, state, c, complextypes.MakeEmptySubstAndForm(), childId, originalNodeId, []int{})
}

func CanApplyTs(state complextypes.State) global.Pair[int, int] {
	for i := range state.GetEqs() {
		for j := range state.GetNeqs() {
			pair := global.MakePair[int, int](i, j)
			if !isAlreadyApplied(state, pair) {
				return pair
			}
		}
	}
	return global.MakePair[int, int](-1, -1)
}

func isAlreadyApplied(state complextypes.State, pair global.Pair[int, int]) bool {
	for _, each := range state.GetAlreadyAppliedZeq() {
		if each.Fst == pair.Fst && each.Snd == pair.Snd {
			return true
		}
	}
	return false
}

func isNilPair(pair global.Pair[int, int]) bool {
	return pair.Fst == -1 && pair.Snd == -1
}
