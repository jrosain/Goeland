package arithmetic

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/GoelandProver/Goeland/global"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

const HiGHS_PATH = "./plugins/arithmetic/HiGHS/build/bin/highs"

type CounterExample struct {
	Variables []string
	Values    []int
}

func GetCounterExample(formNetworks []basictypes.FormList) (example CounterExample, success bool) {
	constraintNetwork := buildConstraintNetwork(formNetworks)
	allNetworks := getAllNetworks(constraintNetwork)

	for _, network := range allNetworks {
		if example, success = tryConstraintNetwork(network); success {
			return example, true
		}
	}

	return CounterExample{}, false
}

func buildConstraintNetwork(formNetworks []basictypes.FormList) []Network {
	a := MakeSimpleConstraint(&Variable{"X"}, GreaterEq, &Constant{0})
	b := MakeSimpleConstraint(&Variable{"X"}, GreaterEq, &Constant{1})
	c := MakeSimpleConstraint(&Variable{"X"}, GreaterEq, &Constant{-5})
	d := MakeSimpleConstraint(&Variable{"X"}, LesserEq, &Constant{0})

	return []Network{{&a, &b}, {&c}, {&d}}
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

func tryConstraintNetwork(network Network) (example CounterExample, success bool) {
	networkFile := "network.lp"
	solutionFile := "solution.out"
	buildFile(network, networkFile)
	runHiGHS(networkFile, solutionFile)
	return gatherData(solutionFile)
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

func gatherData(solutionFile string) (example CounterExample, success bool) {
	readFile, err := os.Open(solutionFile)
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	defer readFile.Close()

	return readAndParseFile(readFile)
}

func readAndParseFile(readFile *os.File) (example CounterExample, success bool) {
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
			variable, value := returnVariableValue(line)
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

func returnVariableValue(line string) (variable string, value int) {
	words := strings.Split(line, " ")
	variable = words[0]
	value, err := strconv.Atoi(words[1])
	if err != nil {
		global.PrintFatal("ARI", err.Error())
	}
	return variable, value
}
