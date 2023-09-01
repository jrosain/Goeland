package global

import (
	"sync"
)

/*
This is a useful task manager that will execute tasks one by one in the same goroutine.

It is safe to try to add another task as part of a task.
*/
type SyncTaskList struct {
	mutex  sync.Mutex
	tasks  []func()
	buffer []func()
}

func (stl *SyncTaskList) AddTask(task func()) {
	stl.mutex.Lock()
	if len(stl.buffer) == 0 {
		defer stl.doTasks()
	}
	defer stl.mutex.Unlock()

	stl.buffer = append(stl.buffer, task)
}

func (stl *SyncTaskList) doTasks() {
	stl.swapBuffer()

	for _, task := range stl.tasks {
		task()
	}
}

func (stl *SyncTaskList) swapBuffer() {
	stl.mutex.Lock()
	defer stl.mutex.Unlock()

	stl.tasks = stl.buffer
	stl.buffer = []func(){}
}
