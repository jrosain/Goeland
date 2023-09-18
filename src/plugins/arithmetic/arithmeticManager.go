package arithmetic

import (
	"sync"

	treetypes "github.com/GoelandProver/Goeland/code-trees/tree-types"
	"github.com/GoelandProver/Goeland/global"
	basictypes "github.com/GoelandProver/Goeland/types/basic-types"
)

type ArithManager struct {
	counter *global.SyncCounter

	communications []chan *SubAnswer
	forms          []basictypes.FormAndTermsList

	mutex sync.Mutex
}

type SubAnswer struct {
	sub     treetypes.Substitutions
	form    basictypes.FormAndTerms
	success bool
}

var Manager *ArithManager = NewArithManager()

func NewArithManager() *ArithManager {
	manager := &ArithManager{communications: []chan *SubAnswer{}, forms: []basictypes.FormAndTermsList{}}
	counter := global.NewSyncCounter(manager.SendSubstitutions)
	manager.counter = counter
	manager.counter.Increment()

	return manager
}

func (am *ArithManager) OpenBranch() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.counter.Increment()
}

func (am *ArithManager) BranchClosure() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.counter.Decrement()
}

func (am *ArithManager) GetArithResult(channel chan *SubAnswer, forms basictypes.FormAndTermsList) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.communications = append(am.communications, channel)
	am.forms = append(am.forms, forms)
	am.counter.Decrement()
}

func (am *ArithManager) SendSubstitutions() {
	if len(am.forms) > 0 {
		example, forms, success := GetCounterExample(am.forms)
		var result treetypes.Substitutions

		if success {
			result = example.convert()
			for i, channel := range am.communications {
				channel <- (&SubAnswer{result, forms[i], true})
			}
		} else {
			for _, channel := range am.communications {
				channel <- (&SubAnswer{nil, basictypes.FormAndTerms{}, false})
			}
		}

		am.communications = []chan *SubAnswer{}
		am.forms = []basictypes.FormAndTermsList{}
	}
}
