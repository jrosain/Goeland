/**
* Copyright 2022 by the authors (see AUTHORS).
*
* Goéland is an automated theorem prover for first order logic.
*
* This software is governed by the CeCILL license under French law and
* abiding by the rules of distribution of free software.  You can  use,
* modify and/ or redistribute the software under the terms of the CeCILL
* license as circulated by CEA, CNRS and INRIA at the following URL
* "http://www.cecill.info".
*
* As a counterpart to the access to the source code and  rights to copy,
* modify and redistribute granted by the license, users are provided only
* with a limited warranty  and the software's author,  the holder of the
* economic rights,  and the successive licensors  have only  limited
* liability.
*
* In this respect, the user's attention is drawn to the risks associated
* with loading,  using,  modifying and/or developing or reproducing the
* software by the user in light of its specific status of free software,
* that may mean  that it is complicated to manipulate,  and  that  also
* therefore means  that it is reserved for developers  and  experienced
* professionals having in-depth computer knowledge. Users are therefore
* encouraged to load and test the software's suitability as regards their
* requirements in conditions enabling the security of their systems and/or
* data to be ensured and,  more generally, to use and operate it in the
* same conditions as regards security.
*
* The fact that you are presently reading this means that you have had
* knowledge of the CeCILL license and that you accept its terms.
**/

package ND

import (
	"fmt"

	"github.com/GoelandProver/Goeland/AST"
	"github.com/GoelandProver/Goeland/Glob"
	"github.com/GoelandProver/Goeland/Lib"
	"github.com/GoelandProver/Goeland/Search"
	"github.com/GoelandProver/Goeland/Unif"
)

type algo struct {
	Search.SearchAlgorithm
}

func Algo() Search.SearchAlgorithm {
	return algo{}
}

func (a algo) Search(form AST.Form, _ int) bool {
	debug(
		Lib.MkLazy(func() string {
			return "Launching proof-search."
		}),
	)
	res := proofSearch(emptyState().dispatch(form))

	if res {
		Glob.PrintInfo("nd", "VALID")
	} else {
		Glob.PrintInfo("nd", "NOT VALID")
	}

	return res
}

func proofSearch(s state) bool {
	switch next_rule_opt := s.nextRule().(type) {
	case Lib.None[rule]:
		debug(
			Lib.MkLazy(func() string {
				return "No more rules to apply"
			}),
		)
		return false

	case Lib.Some[rule]:
		next_rule := next_rule_opt.Val

		debug(
			Lib.MkLazy(func() string {
				return fmt.Sprintf("Applying rule: %s", next_rule.ToString())
			}),
		)

		if next_rule.kind == Search.KindClosure {
			return tryClosure(s, next_rule)
		} else {
			return applyExpansion(s, next_rule)
		}
	}

	raise_anomaly("Reached an unreachable case")
	return false
}

func tryClosure(s state, r rule) bool {
	if findTrivialContradiction(s.lit, r.tgt) {
		debug(Lib.MkLazy(func() string {
			return fmt.Sprintf(
				"Found trivial contradiction with %s in context.",
				r.ToString(),
			)
		}))
		return true
	} else if found, _ := findUnification(r.tgt, s.posLits, s.negLits); found {
		// todo:
		//   we can probably discard the formulas that get returned by the unification
		//   then, we can try the substitutions one by one --- as a backtrack point.
		//   we don't need to instantiate anything in the state, we should simply introduce
		//   gamma instantiation rules. The question is the following: what priority should
		//   such rules have? We probably want to do them directly.
		//
		//   We should !! however !! keep track of which formula meta comes from which formula.
		//   We should instantiate this as we apply the gamma-inst rule.
		//   This implies that we need an order on the metas, like which one depends on which one.
		//   This is probably given by the order in which they have been introduced, so it's probably
		//   easy to manage. We should nevertheless specify that invariant in a comment.

		return false // not yet implemented
	} else {
		s = s.removeNextRule().addLit(r.tgt)

		debug(Lib.MkLazy(func() string {
			return fmt.Sprintf(
				"No contradiction found, adding %s in the litterals --- litterals: {%s}.",
				r.tgt.ToString(),
				Lib.ListToString(s.lit.Elements(), ", ", ""),
			)
		}))

		return proofSearch(s)
	}
}

func applyExpansion(s state, r rule) bool {
	result_formulas, metas := r.apply(
		s.metas,
	) // todo: special case for gamma-inst. Add an (optional) term to instantiate with in the rule?

	diff_metas := metas.Diff(s.metas)
	if r.kind == Search.KindGamma && !diff_metas.IsEmpty() {
		if diff_metas.Cardinal() != 1 {
			raise_anomaly("Gamma rule generated more than one meta")
		}

		meta := diff_metas.Elements().At(0)
		s = s.addMetaGeneration(meta, r.tgt)
	}

	s = s.removeNextRule().withMetas(metas)

	// For gamma rules, we _have to_ make them available again, so we add it back in the state.
	// Indeed, see e.g. the proof of the drinker formula in outer Skolemization to see a formula
	// that needs to reintroduce a free variable (even in non-destructive mode).
	//
	// However, we add it back with an increased weight in order to make the rule application
	// _fair_: we need to do all the other gamma rules before trying to do this one again.
	// If we don't do that and the min heap keeps selecting the same rule, then it can lead to
	// incompleteness (see e.g., Julie's PhD thesis for an example of this)
	if r.kind == Search.KindGamma {
		s = s.addRule(r.withWgt(r.wgt + 1))
	}

	debug(
		Lib.MkLazy(func() string {
			return fmt.Sprintf("Result formulas: [%s]", result_formulas.ToString(
				func(l Lib.List[AST.Form]) string {
					return "[" + Lib.ListToString(l, ", ", "") + "]"
				}, ", ", ""),
			)
		}),
	)

	// FIXME: hide the following code behind a flag to be able to make the proof search sequential (useful for debugging purposes)
	// for _, ls := range result_formulas.GetSlice() {
	// 	state := s.copy()
	// 	for _, form := range ls.GetSlice() {
	// 		state = state.dispatch(form)
	// 	}
	// 	if !proofSearch(state) {
	// 		return false
	// 	}
	// }

	// return true

	calls := []func(chan bool){}

	// todo: debug why it says not valid when searched in parallel
	for _, ls := range result_formulas.GetSlice() {
		loop_list := ls
		loop_state := s.copy()
		calls = append(calls, func(outchan chan bool) {
			for _, form := range loop_list.GetSlice() {
				loop_state = loop_state.dispatch(form)
			}
			outchan <- proofSearch(loop_state)
		})
	}

	res, err := Lib.GenericParallel(calls,
		func(x, y bool) bool { return x && y },
		true,
	)

	if err != nil {
		raise_anomaly(
			fmt.Sprintf("Encountered an error on parallel launch of proof search: %s", err.Error()),
		)
	}

	return res
}

func findTrivialContradiction(lit Lib.Set[AST.Form], tgt AST.Form) bool {
	return tgt.Equals(AST.MakeBot()) || tgt.Equals(AST.MakeNot(AST.MakeTop())) ||
		(Glob.Is[AST.Not](tgt) && lit.Contains(Glob.To[AST.Not](tgt).GetForm())) ||
		(Glob.Is[AST.Pred](tgt) && lit.Contains(AST.MakeNot(tgt)))
}

func findUnification(
	tgt AST.Form,
	pos_tree, neg_tree Unif.DataStructure,
) (bool, []Unif.MixedSubstitutions) {
	if Glob.Is[AST.Not](tgt) {
		return pos_tree.Unify(Glob.To[AST.Not](tgt).GetForm())
	} else {
		return neg_tree.Unify(tgt)
	}
}
