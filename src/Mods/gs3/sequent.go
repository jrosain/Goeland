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
package gs3

import (
	"strings"

	"fmt"

	"github.com/GoelandProver/Goeland/AST"
	"github.com/GoelandProver/Goeland/Glob"
	"github.com/GoelandProver/Goeland/Lib"
	"github.com/GoelandProver/Goeland/Search"
)

type GS3Sequent struct {
	hypotheses    Lib.List[AST.Form]
	rule          Search.TableauxRule
	appliedOn     int
	rewriteWith   Lib.Option[int]
	termGenerated Lib.Option[Lib.Either[AST.Ty, AST.Term]]
	children      Lib.List[*GS3Sequent]
}

func NewSequent() *GS3Sequent {
	seq := new(GS3Sequent)
	seq.hypotheses = Lib.NewList[AST.Form]()
	seq.rewriteWith = Lib.MkNone[int]()
	seq.termGenerated = Lib.MkNone[AST.Term]()
	seq.children = Lib.NewList[*GS3Sequent]()
	return seq
}

// ----------------------------------------------------------------------------
// Public methods
// ----------------------------------------------------------------------------

func (seq *GS3Sequent) GetTargetForm() AST.Form {
	return seq.hypotheses.At(seq.appliedOn)
}

func (seq *GS3Sequent) Child(i int) *GS3Sequent {
	return seq.children.At(i)
}

func (seq *GS3Sequent) Children() Lib.List[*GS3Sequent] {
	return seq.children
}

func (seq *GS3Sequent) Rule() Search.TableauxRule {
	return seq.rule
}

func (seq *GS3Sequent) TermGenerated() Lib.Option[Lib.Either[AST.Ty, AST.Term]] {
	return seq.termGenerated
}

func (seq *GS3Sequent) Empty() bool {
	return seq.hypotheses.Empty()
}

func (seq *GS3Sequent) ToString() string {
	return seq.toStringAux(0)
}

func (seq *GS3Sequent) GetRewriteWith() Lib.Option[AST.Form] {
	return Lib.OptBind(seq.rewriteWith, func(i int) Lib.Option[AST.Form] {
		return Lib.MkSome(seq.hypotheses.At(i))
	})
}

func (seq *GS3Sequent) SetChildren(c Lib.List[*GS3Sequent]) *GS3Sequent {
	seq.children = c
	return seq
}

func (seq *GS3Sequent) SetTargetForm(f AST.Form) *GS3Sequent {
	seq.hypotheses.Upd(seq.appliedOn, f)
	return seq
}

// ----------------------------------------------------------------------------
// Private methods & functions
// ----------------------------------------------------------------------------

func (seq *GS3Sequent) setHypotheses(forms Lib.List[AST.Form]) *GS3Sequent {
	seq.hypotheses = Lib.ListCpy(forms)
	// If equality reasoning has been used to terminate the proof, then an empty predicate is expected
	// (see search_destructive, manageClosureRule on eq reasoning).
	// As such, add an hypothesis with the empty =
	seq.hypotheses.Append(AST.EmptyPredEq)

	return seq
}

func (seq *GS3Sequent) setAppliedRule(rule Search.TableauxRule) *GS3Sequent {
	seq.rule = rule
	return seq
}

func (seq *GS3Sequent) setAppliedOn(hypothesis AST.Form) *GS3Sequent {
	seq.appliedOn = seq.getIndexOf(hypothesis)
	return seq
}

func (seq *GS3Sequent) setTermGenerated(t Lib.Option[Lib.Either[AST.Ty, AST.Term]]) *GS3Sequent {
	seq.termGenerated = t
	return seq
}

func (seq *GS3Sequent) addChild(oth ...*GS3Sequent) *GS3Sequent {
	seq.children.Append(oth...)
	return seq
}

func (seq *GS3Sequent) toStringAux(i int) string {
	identation := strings.Repeat("  ", i)
	status := seq.rule.ToString() + " on " + seq.hypotheses.At(seq.appliedOn).ToString()
	if seq.Empty() {
		status = "EMPTY"
	}
	childrenStrings := make([]string, seq.children.Len())
	for j, child := range seq.children.GetSlice() {
		if child != nil {
			childrenStrings[j] = child.toStringAux(i + 1)
		}
	}
	return "\n" + identation + status + strings.Join(childrenStrings, "")
}

func (seq *GS3Sequent) setRewrittenWith(form AST.Form) *GS3Sequent {
	rewrite_index := seq.getIndexOf(form)
	seq.rewriteWith = Lib.MkSome(rewrite_index)
	return seq
}

func (seq *GS3Sequent) getIndexOf(target AST.Form) int {
	index_opt := Lib.ListIndexOf(target, seq.hypotheses)

	switch index := index_opt.(type) {
	case Lib.Some[int]:
		return int(index.Val)
	case Lib.None[int]:
		debug(
			Lib.MkLazy(func() string {
				return fmt.Sprintf(
					"Tried to get the index of %s in a context composed of the following hypotheses: \n%s",
					target.ToString(),
					Lib.ListToString(seq.hypotheses, "\n", "(empty context)"),
				)
			}),
		)
		Glob.Anomaly(gs3_label, "Failure: tried to get a missing hypothesis")
	}

	return -1
}
