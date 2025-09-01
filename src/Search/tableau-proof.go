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
 * This file provides an interface to ease communication between proof structures.
 **/

package Search

import (
	"fmt"

	"github.com/GoelandProver/Goeland/AST"
	"github.com/GoelandProver/Goeland/Glob"
	"github.com/GoelandProver/Goeland/Lib"
)

type TableauxRule int
type TableauxRuleKind int

const (
	RuleClosure TableauxRule = iota
	RuleImp
	RuleAnd
	RuleOr
	RuleEqu
	RuleEx
	RuleAll
	RuleNotNot
	RuleNotImp
	RuleNotAnd
	RuleNotOr
	RuleNotEqu
	RuleNotEx
	RuleNotAll
	RuleReintro
	RuleRew
)

const (
	KindAlpha TableauxRuleKind = iota
	KindBeta
	KindDelta
	KindGamma
	KindRew
	KindClosure
)

type IProof interface {
	AppliedOn() AST.Form
	RuleApplied() TableauxRule
	KindOfRule() TableauxRuleKind
	ResultFormulas() Lib.List[Lib.List[AST.Form]]
	Children() Lib.List[IProof]
	Child(int) IProof
	RewrittenWith() Lib.Option[AST.Form]
	TermGenerated() Lib.Option[Lib.Either[AST.Ty, AST.Term]]
}

// ----------------------------------------------------------------------------
// Public methods
// ----------------------------------------------------------------------------

func (r TableauxRule) ToString() string {
	switch r {
	case RuleClosure:
		return "closure"
	case RuleAnd:
		return "alpha_and"
	case RuleNotNot:
		return "alpha_not_not"
	case RuleNotImp:
		return "alpha_not_imp"
	case RuleNotOr:
		return "alpha_not_or"
	case RuleImp:
		return "beta_imply"
	case RuleOr:
		return "beta_or"
	case RuleEqu:
		return "beta_equ"
	case RuleNotAnd:
		return "beta_not_and"
	case RuleNotEqu:
		return "beta_not_equ"
	case RuleEx:
		return "delta_ex"
	case RuleNotAll:
		return "delta_not_all"
	case RuleAll:
		return "gamma_all"
	case RuleNotEx:
		return "gamma_not_ex"
	case RuleReintro:
		return "reintroduction"
	case RuleRew:
		return "rewrite"
	}

	Glob.Anomaly("IProof", fmt.Sprintf("Unknown rule %d", r))
	return ""
}

func (r TableauxRule) KindOfRule() TableauxRuleKind {
	switch r {
	case RuleNotNot, RuleNotOr, RuleNotImp, RuleAnd:
		return KindAlpha
	case RuleNotAnd, RuleNotEqu, RuleOr, RuleImp, RuleEqu:
		return KindBeta
	case RuleNotAll, RuleEx:
		return KindDelta
	case RuleNotEx, RuleAll, RuleReintro:
		return KindGamma
	case RuleRew:
		return KindRew
	case RuleClosure:
		return KindClosure
	}

	Glob.Anomaly(label, "Unknown kind of rule")
	return 0
}
