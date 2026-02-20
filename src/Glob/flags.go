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
 * This file contains the boolean-valued options of Goéland.
 * They are stored as integer-valued flags.
 **/

package Glob

import (
	"fmt"
	"github.com/GoelandProver/Goeland/Lib"
)

type OptFlag = int
type conflict = Lib.List[OptFlag]

var (
	None                     OptFlag = 0
	Destructive                      = 1
	NonDestructive                   = 2 << 0
	OutputRocq                       = 2 << 1
	OutputLambdaPi                   = 2 << 2
	OutputTPTP                       = 2 << 3
	OutputSCTPTP                     = 2 << 4
	OutputProof                      = 2 << 5
	InteractiveDebug                 = 2 << 6
	PrettyPrint                      = 2 << 7
	StopAfterOneStep                 = 2 << 8
	TryDMTBeforeEq                   = 2 << 9
	CompleteProofSearch              = 2 << 10
	OutputTypeProof                  = 2 << 11
	UseArithmetic                    = 2 << 12
	UseInnerSkolemization            = 2 << 13
	UsePreInnerSkolemization         = 2 << 14
	PrintVersion                     = 2 << 15
	Flatten                          = 2 << 16
	TypeCheck                        = 2 << 17
	IncrementalEquality              = 2 << 18
	WriteLogsInFile                  = 2 << 19
	SilentMode                       = 2 << 20
)

var flags OptFlag
var conflicts Lib.List[conflict]
var implications map[OptFlag]OptFlag
var stringify map[OptFlag]string
var checked = false

func init_flags() {
	flags = Destructive | PrintVersion | TypeCheck

	conflicts = Lib.MkListV(
		Lib.MkListV(Destructive, NonDestructive),
		Lib.MkListV(SilentMode, WriteLogsInFile),
		Lib.MkListV(UseInnerSkolemization, UsePreInnerSkolemization),
		Lib.MkListV(OutputRocq, OutputLambdaPi, OutputTPTP, OutputSCTPTP),
	)

	// /!\ a cycle in the implications between flags will lead to an infinite loop /!\
	implications = map[OptFlag]OptFlag{
		OutputRocq:     OutputProof,
		OutputLambdaPi: OutputProof,
		OutputTPTP:     OutputProof,
		OutputSCTPTP:   OutputProof,
	}

	stringify = map[OptFlag]string{
		None:                     "None",
		Destructive:              "Destructive",
		NonDestructive:           "NonDestructive",
		OutputRocq:               "OutputRocq",
		OutputLambdaPi:           "OutputLambdaPi",
		OutputTPTP:               "OutputTPTP",
		OutputSCTPTP:             "OutputSCTPTP",
		OutputProof:              "OutputProof",
		InteractiveDebug:         "InteractiveDebug",
		PrettyPrint:              "PrettyPrint",
		StopAfterOneStep:         "StopAfterOneStep",
		TryDMTBeforeEq:           "TryDMTBeforeEq",
		CompleteProofSearch:      "CompleteProofSearch",
		OutputTypeProof:          "OutputTypeProof",
		UseArithmetic:            "UseArithmetic",
		UseInnerSkolemization:    "UseInnerSkolemization",
		UsePreInnerSkolemization: "UsePreInnerSkolemization",
		PrintVersion:             "PrintVersion",
		Flatten:                  "Flatten",
		TypeCheck:                "TypeCheck",
		IncrementalEquality:      "IncrementalEquality",
		WriteLogsInFile:          "WriteLogsInFile",
		SilentMode:               "SilentMode",
	}
}

func GetFlag(flag OptFlag) bool {
	if !checked {
		Fatal(
			"opt",
			"Flags have not been checked. This may lead to an unexpected behaviour of Goéland",
		)
	}
	return flags&flag != 0
}

func SetFlag(flag OptFlag) {
	flags = flags | flag
	infer_flags(flag)
}

func infer_flags(flag OptFlag) {
	for i := 0; i <= 20; i++ {
		if flag&(2<<i) != 0 {
			if impl, found := implications[2<<i]; found {
				SetFlag(impl)
			}
		}
	}
}

func CheckFlags() {
	checked = true
	for _, conflict := range conflicts.GetSlice() {
		number_flags_activated := 0
		for _, flag := range conflict.GetSlice() {
			if GetFlag(flag) {
				number_flags_activated++
			}
		}
		if number_flags_activated > 1 {
			Fatal(
				"opt",
				"Options "+conflict_to_string(conflict)+" cannot be activated simultaneously",
			)
		}
	}
}

func conflict_to_string(c conflict) string {
	return c.ToString(
		func(flag OptFlag) string {
			str, found := stringify[flag]
			if !found {
				Anomaly("opt", fmt.Sprintf("Flag %d cannot be stringified", flag))
			}
			return str
		}, "|", "None",
	)
}
