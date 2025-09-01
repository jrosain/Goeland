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

/**
 * This file provides a definition of a GS3 proof: the object & the exported functions.
 **/

package gs3

import (
	"fmt"

	"github.com/GoelandProver/Goeland/AST"
	"github.com/GoelandProver/Goeland/Glob"
	"github.com/GoelandProver/Goeland/Lib"
	"github.com/GoelandProver/Goeland/Mods/dmt"
	"github.com/GoelandProver/Goeland/Search"
)

var debug Glob.Debugger
var gs3_label = "GS3"

// ----------------------------------------------------------------------------
// Structures used to deskolemize
// ----------------------------------------------------------------------------

// A branch is a list of integers, each one indicating which "path" to take from
// the root to reach the current node.
type Branch Lib.List[int]

func makeBranch(src Branch) Branch {
	return Branch(Lib.MkListV(Lib.List[int](src).GetSlice()...))
}

func (b Branch) Append(v int) Branch {
	ls := Lib.List[int](b)
	ls.Append(v)
	return Branch(ls)
}

// FIXME: Branches are related one to another by a prefix relation.
// I'll do it later.

// The deskolemization algorithm runs in a context that records:
//   - the current hypotheses,
//   - the current branch of the context,
//   - the formula & rule applied to get a formula,
//   - the mapping of formulas to their original branch,
//   - anything else?
type GS3Con struct {
	forms       Lib.List[AST.Form]
	branch      Branch
	form_origin Lib.Map[AST.Form, Lib.Pair[AST.Form, Search.TableauxRule]]
	form_branch Lib.Map[AST.Form, Branch]
}

func (con GS3Con) Copy() GS3Con {
	return GS3Con{
		forms:       Lib.ListCpy(con.forms),
		branch:      makeBranch(con.branch),
		form_origin: con.form_origin.Copy(),
		form_branch: con.form_branch.Copy(),
	}
}

// Below are methods that take care of the origin branch of a formula, the rule from which
// it was created, etc.
func (con GS3Con) addFormulas(forms Lib.List[AST.Form], proof Search.TableauxProof) GS3Con {
	for _, form := range forms.GetSlice() {
		con.forms.Append(form)
		con.form_branch = con.form_branch.Set(form, con.branch)
		con.form_origin = con.form_origin.Set(
			form,
			Lib.MkPair(proof.AppliedOn(), proof.RuleApplied()),
		)
	}
	return con
}

func (con GS3Con) replaceFormula(index int, target AST.Form) GS3Con {
	con.forms.Upd(index, target)
	con.form_branch = con.form_branch.Set(target, con.branch)
	con.form_origin = con.form_origin.Set(target, Lib.MkPair(
		con.forms.At(index),
		Search.RuleRew,
	))
	return con
}

// ----------------------------------------------------------------------------
// Public methods: debugger and entry point
// ----------------------------------------------------------------------------

func InitDebugger() {
	debug = Glob.CreateDebugger(gs3_label)
}

var MakeGS3Proof = func(proof Search.TableauxProof) *GS3Sequent {
	context := GS3Con{
		forms:       Lib.MkListV(proof[0].Formula.GetForm()),
		form_origin: Lib.MkMap[AST.Form, Lib.Pair[AST.Form, Search.TableauxRule]](),
		form_branch: Lib.MkMap[AST.Form, Branch](),
	}

	// FIXME: this should be a hook
	if Glob.IsLoaded("dmt") {
		context.forms.Append(dmt.GetRegisteredAxioms().GetSlice()...)
	}

	root := deskolemize(context, proof)
	return root
}

// ----------------------------------------------------------------------------
// Private methods: deskolemization procedure
// ----------------------------------------------------------------------------

func deskolemize(context GS3Con, proof Search.TableauxProof) *GS3Sequent {
	if len(proof) == 0 {
		Glob.Anomaly(gs3_label, "Tried to deskolemize an empty proof")
	}
	contexts, sequent := makeOneProofStep(context, proof)

	for i := 0; i < contexts.Len(); i++ {
		child_sequent := deskolemize(contexts.At(i), proof.Child(i).(Search.TableauxProof))

		if child_sequent == nil {
			Glob.Anomaly(gs3_label, "Deskolemization has yielded an empty sequent")
		}

		sequent.children.Append(child_sequent)
	}

	return sequent
}

// Invariant: the output contexts length should be equal to the number of children the sequent has.
func makeOneProofStep(context GS3Con, proof Search.TableauxProof) (Lib.List[GS3Con], *GS3Sequent) {
	// Makes a new context with the i-th result formulas
	mkNewContext := func(ctx GS3Con, i int) GS3Con {
		return ctx.addFormulas(proof.ResultFormulas().At(i), proof)
	}

	sequent := NewSequent().
		setHypotheses(context.forms).
		setAppliedRule(proof.RuleApplied()).
		setAppliedOn(proof.AppliedOn())

	switch proof.KindOfRule() {
	case Search.KindAlpha:
		return Lib.MkListV(mkNewContext(context, 0)), sequent

	case Search.KindGamma:
		sequent = sequent.setTermGenerated(proof.RewrittenWith())
		return Lib.MkListV(mkNewContext(context, 0)), sequent

	case Search.KindBeta:
		contexts := Lib.MkList[GS3Con](proof.Children().Len())
		for i := 0; i < proof.Children().Len(); i++ {
			local_ctx := context.Copy()
			contexts.Upd(i, mkNewContext(local_ctx, i))
		}
		return contexts, sequent

	case Search.KindDelta:
		switch term := proof.TermGenerated().(type) {
		case Lib.Some[Lib.Either[AST.Ty, AST.Term]]:
			// FIXME: weakening phase I, formulas containing [term] not of this branch
			//        more tricky than expected, maybe the branch should be recorded in the TableauxProof
			// FIXME: weakening phase II, formulas containing [term] of this branch
			// FIXME: apply delta rule
			// FIXME: reapply weakened formulas from phase II
		}
		return Lib.MkListV(mkNewContext(context, 0)), sequent

	case Search.KindClosure:
		return Lib.NewList[GS3Con](), sequent

	case Search.KindRew:
		switch f := proof.RewrittenWith().(type) {
		case Lib.Some[AST.Form]:
			sequent = sequent.setRewrittenWith(f.Val)
			return Lib.MkListV(context.replaceFormula(
				sequent.appliedOn,
				proof.ResultFormulas().At(0).At(0),
			)), sequent
		}
		Glob.Anomaly(gs3_label, "Rewrite rule has been applied without rewrite formula")
		return Lib.NewList[GS3Con](), sequent
	}

	Glob.Fatal(gs3_label, fmt.Sprintf(
		"Rule %s is not yet implemented by the deskolemization algorithm",
		proof.RuleApplied().ToString(),
	))
	return Lib.MkListV(context), nil
}
