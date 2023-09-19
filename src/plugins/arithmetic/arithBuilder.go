package arithmetic

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	treetypes "github.com/GoelandProver/Goeland/code-trees/tree-types"
	"github.com/GoelandProver/Goeland/global"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

const HiGHS_PATH = "./plugins/arithmetic/HiGHS/build/bin/highs"

type CounterExample struct {
	Variables []basictypes.Term
	Values    []Numeric
}

func (ce *CounterExample) ToString() string {
	str := ""
	for i := range ce.Variables {
		str += fmt.Sprintf("%s -> %v, ", ce.Variables[i].ToString(), ce.Values[i])
	}
	return str[:len(str)-2]
}

func (ce *CounterExample) convert() treetypes.Substitutions {
	result := treetypes.MakeEmptySubstitution()

	for i := range ce.Variables {
		if meta, ok := ce.Variables[i].(basictypes.Meta); ok {
			addToConstantMap(ce.Values[i])
			result = append(result, treetypes.MakeSubstitution(meta, getValueTerm(ce.Values[i])))
		}
	}

	return result
}

func IsArithClosure(form basictypes.Form) bool {
	metas := form.GetMetas()

	if len(metas) == 0 {
		switch typed := form.(type) {
		case basictypes.Not:
			if predTyped, ok := typed.GetForm().(basictypes.Pred); ok {
				comparaison, _ := convertPred(predTyped)
				return comparaison.Evaluate()
			}
		case basictypes.Pred:
			comparaison, _ := convertPred(typed)
			return !comparaison.Evaluate()
		}
	}

	return false
}

func IsArithmeticable(forms basictypes.FormAndTermsList) bool {
	for _, form := range forms.ExtractForms() {
		switch typed := form.(type) {
		case basictypes.Not:
			if predTyped, ok := typed.GetForm().(basictypes.Pred); ok {
				if converted, _ := convertPred(predTyped); converted != nil {
					return true
				}
			}
		case basictypes.Pred:
			if converted, _ := convertPred(typed); converted != nil {
				return true
			}
		}
	}

	return false
}

func GetArithResult(forms basictypes.FormAndTermsList) (subs []treetypes.Substitutions, form basictypes.FormAndTerms, success bool) {
	channel := make(chan *SubAnswer)
	go Manager.GetArithResult(channel, forms)
	answer := <-channel

	if answer.success {
		success = true
		subs = []treetypes.Substitutions{answer.sub}
		form = answer.form
	}

	return subs, form, success
}

func GetCounterExample(formNetworks []basictypes.FormAndTermsList) (example CounterExample, formsUsed basictypes.FormAndTermsList, success bool) {
	constraintNetwork, termMap := buildConstraintNetwork(formNetworks)
	allNetworks, allForms := getAllNetworks(constraintNetwork, formNetworks)

	for i, network := range allNetworks {
		if example, success = tryConstraintNetwork(network, termMap); success {
			return example, allForms[i], true
		}
	}

	return CounterExample{}, basictypes.FormAndTermsList{}, false
}

func buildConstraintNetwork(formNetworks []basictypes.FormAndTermsList) ([]Network, map[string]basictypes.Term) {
	networks := []Network{}
	termMap := make(map[string]basictypes.Term)

	for _, formNetwork := range formNetworks {
		network := Network{}
		for _, form := range formNetwork {
			var compForm ComparisonForm
			var newMap map[string]basictypes.Term

			switch typed := form.GetForm().(type) {
			case basictypes.Pred:
				compForm, newMap = convertPred(typed)
			case basictypes.Not:
				if pred, ok := typed.GetForm().(basictypes.Pred); ok {
					compForm, newMap = convertPred(pred)
					compForm = compForm.Reverse()
				}
			}

			network = append(network, compForm.Simplify())
			for k, v := range newMap {
				termMap[k] = v
			}
		}
		networks = append(networks, network)
	}

	return networks, termMap
}

func getAllNetworks(networks []Network, forms []basictypes.FormAndTermsList) ([]Network, []basictypes.FormAndTermsList) {
	if len(networks) == 0 {
		return []Network{{}}, []basictypes.FormAndTermsList{{}}
	} else {
		nextsNetworks, nextsForms := getAllNetworks(networks[1:], forms[1:])
		resNetworks := []Network{}
		resForms := []basictypes.FormAndTermsList{}

		for _, constraint := range networks[0] {
			for _, next := range nextsNetworks {
				resNetworks = append(resNetworks, append(next, constraint))
			}
		}

		for _, form := range forms[0] {
			for _, next := range nextsForms {
				resForms = append(resForms, append(next, form))
			}
		}

		return resNetworks, resForms
	}
}

func tryConstraintNetwork(network Network, termMap map[string]basictypes.Term) (example CounterExample, success bool) {
	networkFile := "network.lp"
	solutionFile := "solution.out"

	terms := make([]string, len(termMap))
	i := 0
	for k := range termMap {
		terms[i] = k
		i++
	}

	buildFile(network, networkFile, termMap)
	runHiGHS(networkFile, solutionFile)
	return gatherData(solutionFile, termMap)
}

func buildFile(network Network, networkFile string, terms map[string]basictypes.Term) {
	f, err := os.Create(networkFile)
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	defer f.Close()

	str := "Subject To\n"
	str += network.ToString()
	str += "Bounds\n"
	for key := range terms {
		str += key + " free\n"
	}
	str += "General\n"
	for key, term := range terms {
		if meta, ok := term.(basictypes.Meta); ok {
			if meta.GetTypeHint().ToString() == "$int" {
				str += key + "\n"
			}
		}
	}
	str += "End"

	_, err = f.WriteString(str)
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
}

func runHiGHS(networkFile, solutionFile string) {
	highsWriteSolutionArg := "--solution_file"
	cmd := exec.Command(HiGHS_PATH, highsWriteSolutionArg, solutionFile, networkFile)
	_, err := cmd.Output()

	if err != nil {
		global.PrintFatal("HiGHS", err.Error())
	}

	automaticLogFile := "HiGHS.log"
	err = os.Remove(automaticLogFile)
	if err != nil {
		global.PrintFatal("HiGHS", err.Error())
	}
}

func gatherData(solutionFile string, termMap map[string]basictypes.Term) (example CounterExample, success bool) {
	readFile, err := os.Open(solutionFile)
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	defer readFile.Close()

	return readAndParseFile(readFile, termMap)
}

func readAndParseFile(readFile *os.File, termMap map[string]basictypes.Term) (example CounterExample, success bool) {
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	cpt := 0
	varAmount := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()

		successValue := "Optimal"
		successLine := 1
		if cpt == successLine && line != successValue {
			return CounterExample{}, false
		}

		varAmountLine := 6
		if cpt == varAmountLine {
			varAmount = returnVarAmount(line)
		}

		if cpt > varAmountLine && cpt <= varAmountLine+varAmount {
			variable, value := returnVariableValue(line, termMap)
			example.Variables = append(example.Variables, variable)
			example.Values = append(example.Values, value)
		}

		cpt++
	}

	return example, true
}

func returnVarAmount(line string) int {
	words := strings.Split(line, " ")

	varAmount, err := strconv.Atoi(words[len(words)-1])
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}

	return varAmount
}

func returnVariableValue(line string, termMap map[string]basictypes.Term) (variable basictypes.Term, value Numeric) {
	words := strings.Split(line, " ")
	variable = termMap[words[0]]
	val, err := strconv.ParseFloat(words[1], 64)
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	return variable, Numeric(val)
}
