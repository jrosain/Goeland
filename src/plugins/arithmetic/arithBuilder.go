package arithmetic

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/GoelandProver/Goeland/global"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

const HiGHS_PATH = "./plugins/arithmetic/HiGHS/build/bin/highs"

type CounterExample struct {
	Variables []basictypes.Term
	Values    []int
}

func (ce *CounterExample) ToString() string {
	str := ""
	for i := range ce.Variables {
		str += fmt.Sprintf("%s -> %v, ", ce.Variables[i].ToString(), ce.Values[i])
	}
	return str[:len(str)-2]
}

func IsArithClosure(form basictypes.Form) bool {
	if len(form.GetMetas()) == 0 {
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

func GetCounterExample(formNetworks [][]basictypes.Pred) (example CounterExample, success bool) {
	constraintNetwork, termMap := buildConstraintNetwork(formNetworks)
	allNetworks := getAllNetworks(constraintNetwork)

	for _, network := range allNetworks {
		if example, success = tryConstraintNetwork(network, termMap); success {
			return example, true
		}
	}

	return CounterExample{}, false
}

func buildConstraintNetwork(predNetworks [][]basictypes.Pred) ([]Network, map[string]basictypes.Term) {
	networks := []Network{}
	termMap := make(map[string]basictypes.Term)

	for _, predNetwork := range predNetworks {
		network := Network{}
		for _, pred := range predNetwork {
			form, newMap := convertPred(pred)
			network = append(network, form.Simplify())

			for k, v := range newMap {
				termMap[k] = v
			}
		}
		networks = append(networks, network)
	}

	return networks, termMap
}

func getAllNetworks(networks []Network) []Network {
	if len(networks) == 0 {
		return []Network{{}}
	} else {
		nexts := getAllNetworks(networks[1:])
		result := []Network{}

		for _, constraint := range networks[0] {
			for _, next := range nexts {
				result = append(result, append(next, constraint))
			}
		}

		return result
	}
}

func tryConstraintNetwork(network Network, termMap map[string]basictypes.Term) (example CounterExample, success bool) {
	networkFile := "network.lp"
	solutionFile := "solution.out"
	buildFile(network, networkFile)
	runHiGHS(networkFile, solutionFile)
	return gatherData(solutionFile, termMap)
}

func buildFile(network Network, networkFile string) {
	f, err := os.Create(networkFile)
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	defer f.Close()

	str := "Subject To\n"
	str += network.ToString()
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

func returnVariableValue(line string, termMap map[string]basictypes.Term) (variable basictypes.Term, value int) {
	words := strings.Split(line, " ")
	variable = termMap[words[0]]
	value, err := strconv.Atoi(words[1])
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	return variable, value
}
