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

// This file implements the environment of the non-destructive proof.

package ND

import (
	"container/heap"
	"fmt"
	"slices"

	"github.com/GoelandProver/Goeland/AST"
	"github.com/GoelandProver/Goeland/Core"
	"github.com/GoelandProver/Goeland/Glob"
	"github.com/GoelandProver/Goeland/Lib"
	"github.com/GoelandProver/Goeland/Search"
	"github.com/GoelandProver/Goeland/Unif"
)

func isEquNEqu(f AST.Form) bool {
	return (Glob.Is[AST.Not](f) && Glob.Is[AST.Equ](Glob.To[AST.Not](f).GetForm())) ||
		Glob.Is[AST.Equ](f)
}

// We use a min-heap for expansion rule application.
// This allows us to simply store this heap in the state and to apply the given rule directly.
// Each expansion rule has a weigth which also makes it possible to insert intermediary rules
// (e.g., dmt, equality, etc).
type rule struct {
	wgt  int
	kind Search.TableauxRuleKind
	tgt  AST.Form
}

func (r rule) ToString() string {
	return fmt.Sprintf("K[%s] on %s --- weight: %d", r.kind.ToString(), r.tgt.ToString(), r.wgt)
}

func (r rule) withWgt(wgt int) rule {
	return rule{wgt, r.kind, r.tgt}
}

func (r rule) apply(metas Lib.Set[AST.Meta]) (Lib.List[Lib.List[AST.Form]], Lib.Set[AST.Meta]) {
	negateList := func(fl Lib.List[AST.Form]) Lib.List[AST.Form] {
		return Lib.ListMap(
			fl,
			func(f AST.Form) AST.Form { return AST.MakeNot(f) },
		)
	}

	// todo: maybe use real rules instead, there are too many special cases to make this factorization worth it
	switch r.kind {
	case Search.KindClosure:
		raise_anomaly("tried to get result formulas out of a closure rule")
	case Search.KindAlpha:
		result_formulas := r.tgt.GetChildFormulas()
		if Glob.Is[AST.Not](r.tgt) {
			tgt := Glob.To[AST.Not](r.tgt).GetForm()
			if Glob.Is[AST.Not](tgt) {
				result_formulas = Lib.MkListV(Glob.To[AST.Not](tgt).GetForm())
			} else if Glob.Is[AST.Imp](tgt) {
				timp := Glob.To[AST.Imp](tgt)
				result_formulas =
					Lib.MkListV[AST.Form](timp.GetF1(), AST.MakeNot(timp.GetF2()))
			} else {
				result_formulas = negateList(tgt.GetChildFormulas())
			}
		}
		return Lib.MkListV(result_formulas), metas
	case Search.KindDelta:
		result := Core.Skolemize(r.tgt, metas)
		return Lib.MkListV(Lib.MkListV(result)), metas
	case Search.KindBeta:
		// Special case for equivalence
		if isEquNEqu(r.tgt) {
			if Glob.Is[AST.Equ](r.tgt) {
				f1 := Glob.To[AST.Equ](r.tgt).GetF1()
				f2 := Glob.To[AST.Equ](r.tgt).GetF2()

				return Lib.MkListV(
					Lib.MkListV(f1, f2),
					Lib.MkListV[AST.Form](AST.MakeNot(f1), AST.MakeNot(f2)),
				), metas
			} else {
				f1 := Glob.To[AST.Equ](Glob.To[AST.Not](r.tgt).GetForm()).GetF1()
				f2 := Glob.To[AST.Equ](Glob.To[AST.Not](r.tgt).GetForm()).GetF2()

				return Lib.MkListV(
					Lib.MkListV[AST.Form](f1, AST.MakeNot(f2)),
					Lib.MkListV[AST.Form](AST.MakeNot(f1), f2),
				), metas
			}
		} else if Glob.Is[AST.Imp](r.tgt) {
			f1 := Glob.To[AST.Imp](r.tgt).GetF1()
			f2 := Glob.To[AST.Imp](r.tgt).GetF2()

			return Lib.MkListV(Lib.MkListV[AST.Form](AST.MakeNot(f1)), Lib.MkListV(f2)), metas

		} else {
			result_formulas := r.tgt.GetChildFormulas()
			if Glob.Is[AST.Not](r.tgt) {
				result_formulas = negateList(Glob.To[AST.Not](r.tgt).GetForm().GetChildFormulas())
			}

			return Lib.ListMap(
				result_formulas,
				func(f AST.Form) Lib.List[AST.Form] {
					return Lib.MkListV(f)
				},
			), metas
		}
	case Search.KindGamma:
		result, new_metas := Core.Instantiate(
			Core.MakeFormAndTerm(r.tgt, Lib.NewList[AST.Term]()),
			-1,
		)
		return Lib.MkListV(Lib.MkListV(result.GetForm())), metas.Union(new_metas)
	}

	Glob.Fatal("nd rules", fmt.Sprintf("Unmanaged rule kind: %s"))
	return Lib.NewList[Lib.List[AST.Form]](), metas
}

func ruleKind(form AST.Form) Lib.Option[Search.TableauxRuleKind] {
	switch form := form.(type) {
	case AST.Not:
		switch form.GetForm().(type) {
		case AST.Pred, AST.Top:
			return Lib.MkSome(Search.KindClosure)
		case AST.Not, AST.Or, AST.Imp:
			return Lib.MkSome(Search.KindAlpha)
		case AST.And, AST.Equ:
			return Lib.MkSome(Search.KindBeta)
		case AST.All:
			return Lib.MkSome(Search.KindDelta)
		case AST.Ex:
			return Lib.MkSome(Search.KindGamma)
		}
	case AST.Pred, AST.Bot:
		return Lib.MkSome(Search.KindClosure)
	case AST.And:
		return Lib.MkSome(Search.KindAlpha)
	case AST.Or, AST.Imp, AST.Equ:
		return Lib.MkSome(Search.KindBeta)
	case AST.All:
		return Lib.MkSome(Search.KindGamma)
	case AST.Ex:
		return Lib.MkSome(Search.KindDelta)
	}

	return Lib.MkNone[Search.TableauxRuleKind]()
}

func weightOf(kind Search.TableauxRuleKind) int {
	switch kind {
	case Search.KindClosure:
		return 10
	case Search.KindAlpha:
		return 100
	case Search.KindDelta:
		return 200
	case Search.KindBeta:
		return 300
	case Search.KindGamma:
		return 1000
	}

	Glob.Fatal("nd rules", fmt.Sprintf("Unmanaged rule kind: %s"))
	return -1
}

type rule_heap []rule

// Heap interface
func (h *rule_heap) Len() int {
	return len(*h)
}

func (h *rule_heap) Less(i, j int) bool {
	return (*h)[i].wgt < (*h)[j].wgt
}

func (h *rule_heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *rule_heap) Push(x any) {
	*h = append(*h, x.(rule))
}

func (h *rule_heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type state struct {
	rules     *rule_heap
	metas     Lib.Set[AST.Meta]
	lit       Lib.Set[AST.Form]
	posLits   Unif.DataStructure
	negLits   Unif.DataStructure
	metaIntro Lib.AssqList[AST.Meta, AST.Form]
}

func emptyState() state {
	return state{
		rules:     &rule_heap{},
		metas:     Lib.EmptySet[AST.Meta](),
		lit:       Lib.EmptySet[AST.Form](),
		posLits:   Unif.MakeUnifProblem(Lib.NewList[AST.Form](), true),
		negLits:   Unif.MakeUnifProblem(Lib.NewList[AST.Form](), false),
		metaIntro: Lib.EmptyAssqList[AST.Meta, AST.Form](),
	}
}

func (s state) copy() state {
	return state{
		rules:     s.rules, // Should be safe not to copy: the heap is "functional"
		metas:     s.metas.Copy(),
		lit:       s.lit.Copy(),
		posLits:   s.posLits,   // Should be safe not to copy: no side effects in unification
		negLits:   s.negLits,   // ^
		metaIntro: s.metaIntro, // Should be safe to copy: we use it functionally
	}
}

func (s state) addMetaGeneration(meta AST.Meta, form AST.Form) state {
	return state{s.rules, s.metas, s.lit, s.posLits, s.negLits, s.metaIntro.Push(meta, form)}
}

func (s state) withMetas(metas Lib.Set[AST.Meta]) state {
	return state{s.rules, metas, s.lit, s.posLits, s.negLits, s.metaIntro}
}

func (s state) addLit(f AST.Form) state {
	if Glob.Is[AST.Not](f) {
		s.negLits = s.negLits.InsertFormulaListToDataStructure(
			Lib.MkListV(Glob.To[AST.Not](f).GetForm()),
		)
	} else {
		s.posLits = s.posLits.InsertFormulaListToDataStructure(Lib.MkListV(f))
	}
	return state{s.rules, s.metas, s.lit.Add(f), s.posLits, s.negLits, s.metaIntro}
}

func (s state) addRule(r rule) state {
	rules := slices.Clone(*s.rules)
	heap.Push(&rules, r)
	return state{&rules, s.metas, s.lit, s.posLits, s.negLits, s.metaIntro}
}

func (s state) nextRule() Lib.Option[rule] {
	if s.rules.Len() == 0 {
		return Lib.MkNone[rule]()
	} else {
		rules := slices.Clone(*s.rules)
		return Lib.MkSome(heap.Pop(&rules).(rule))
	}
}

func (s state) removeNextRule() state {
	rules := slices.Clone(*s.rules)
	heap.Pop(&rules)
	return state{&rules, s.metas, s.lit, s.posLits, s.negLits, s.metaIntro}
}

func (s state) dispatch(form AST.Form) state {
	debug(
		Lib.MkLazy(func() string {
			return fmt.Sprintf("Dispatching %s.", form.ToString())
		}),
	)

	switch kind := ruleKind(form).(type) {
	case Lib.None[Search.TableauxRuleKind]:
		debug(
			Lib.MkLazy(func() string {
				return fmt.Sprintf("No implemented rule to dispatch %s.", form.ToString())
			}),
		)
		return s
	case Lib.Some[Search.TableauxRuleKind]:
		debug(
			Lib.MkLazy(func() string {
				return fmt.Sprintf("Can dispatch %s as K[%s] with priority %d.",
					form.ToString(), kind.Val.ToString(), weightOf(kind.Val))
			}),
		)

		return s.addRule(
			rule{
				wgt:  weightOf(kind.Val),
				kind: kind.Val,
				tgt:  form,
			},
		)
	}

	raise_anomaly("Reached an unreachable case")
	return s
}
