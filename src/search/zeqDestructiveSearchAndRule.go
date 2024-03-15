package search

import (
	"fmt"

	"github.com/GoelandProver/Goeland/global"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
	complextypes "github.com/GoelandProver/Goeland/types/complex-types"
	visualization "github.com/GoelandProver/Goeland/visualization_exchanges"
)

var zeqEnable = false

func EnableZeqDestructiveSearch() {
	global.PrintInfo("ZEQ", "ZEQ plugin enabled")
	zeqEnable = true
}

func (ds *destructiveSearch) zeqApplyRule(fatherId uint64, state complextypes.State, c Communication, newAtomics basictypes.FormAndTermsList, currentNodeId int, originalNodeId int, metaToReintroduce []int) {

	var eqs basictypes.FormAndTermsList
	var neqs basictypes.FormAndTermsList

	global.PrintDebug("AR", "ApplyRule")
	switch {
	case len(newAtomics) > 0 && global.IsLoaded("dmt") && len(state.GetSubstsFound()) == 0:
		ds.manageRewriteRules(fatherId, state, c, newAtomics, currentNodeId, originalNodeId, metaToReintroduce)

	case len(state.GetAlpha()) > 0:
		ds.manageAlphaRules(fatherId, state, c, originalNodeId)

	// [TEMP] the case for zeq rules
	case len(eqs) > 0 && len(neqs) > 0:
	//	ds.manageZeqRules(fatherId, state, c, originalNodeId)

	case len(state.GetDelta()) > 0:
		ds.manageDeltaRules(fatherId, state, c, originalNodeId)

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

func (ds *destructiveSearch) manageZeqRules(fatherId uint64, state complextypes.State, c Communication, originalNodeId int, eqs, neqs basictypes.FormAndTermsList) {
	global.PrintDebug("PS", "Zeq rule")
	hdfEq := eqs[0]
	hdfNeq := neqs[0]
	global.PrintDebug("PS", fmt.Sprintf("Rule applied on : %s %s", hdfEq.ToString(), hdfNeq.ToString()))

}
