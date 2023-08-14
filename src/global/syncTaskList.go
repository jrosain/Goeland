package global

import (
	"sync"
)

/*
This is a useful task manager that will execute tasks one by one in the same goroutine.
*/
type SyncTaskList struct {
	mutex sync.Mutex
	tasks []func()
}

func (stl *SyncTaskList) AddTask(task func()) {
	stl.mutex.Lock()
	if len(stl.tasks) == 0 {
		defer stl.doTasks()
	}
	defer stl.mutex.Unlock()

	stl.tasks = append(stl.tasks, task)
}

func (stl *SyncTaskList) doTasks() {
	for len(stl.tasks) > 0 {
		firstTask := stl.getFirstTask()
		firstTask()
	}
}

func (stl *SyncTaskList) getFirstTask() func() {
	stl.mutex.Lock()
	defer stl.mutex.Unlock()

	current := stl.tasks[0]
	stl.tasks = stl.tasks[1:]

	return current
}
