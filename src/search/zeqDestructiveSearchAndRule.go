package search

import (
	"fmt"

	"github.com/GoelandProver/Goeland/global"
	typing "github.com/GoelandProver/Goeland/polymorphism/typing"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
	complextypes "github.com/GoelandProver/Goeland/types/complex-types"
	visualization "github.com/GoelandProver/Goeland/visualization_exchanges"
)

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
	fmt.Printf("Equations : %d, Inequations : %d\n", eqs.Len(), neqs.Len())

	global.PrintDebug("AR", "ApplyRule")
	switch {
	case len(newAtomics) > 0 && global.IsLoaded("dmt") && len(state.GetSubstsFound()) == 0:
		ds.manageRewriteRules(fatherId, state, c, newAtomics, currentNodeId, originalNodeId, metaToReintroduce)

	case len(state.GetAlpha()) > 0:
		ds.manageAlphaRules(fatherId, state, c, originalNodeId)

	case len(state.GetDelta()) > 0:
		ds.manageDeltaRules(fatherId, state, c, originalNodeId)

	case (len(eqs) > 0) && (len(neqs) > 0):
		fmt.Printf("TRANSSYM\n")

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

func (ds *destructiveSearch) applyZeqRules(fatherId uint64, state complextypes.State, c Communication, originalNodeId int, eqs, neqs basictypes.FormAndTermsList) {
	global.PrintDebug("PS", "Zeq rule")
	hdfEq := eqs[0]
	hdfNeq := neqs[0]
	global.PrintDebug("PS", fmt.Sprintf("Rule applied on : %s %s", hdfEq.ToString(), hdfNeq.ToString()))

	resultForms := basictypes.MakeEmptyFormAndTermsList()

	if typedEq, ok := hdfEq.GetForm().(basictypes.Pred); ok {
		if typedNeq, ok := hdfNeq.GetForm().(basictypes.Pred); ok {
			s, t := typedEq.GetArgs().Get(0), typedEq.GetArgs().Get(1)
			u, v := typedNeq.GetArgs().Get(0), typedNeq.GetArgs().Get(1)

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
		}

	}

}
