package global

import (
	"sync"
)

/*
This is a useful counter that executes a given function when the counter gets to 0.

The SyncCounter can be used concurrently without issue. Although care should be taken when incrementing and decrementing, as the counter will execute the given function any time Decrement() is called and the counter is equals to 0 after.

The lock on the counter will only unlock after each Decrement() call, thus if the given function is called, it will only unlock after that function returns.
*/
type SyncCounter struct {
	counter int
	mutex   sync.Mutex

	WhenFinished func()
}

func NewSyncCounter(WhenFinished func()) *SyncCounter {
	sc := new(SyncCounter)
	sc.WhenFinished = WhenFinished
	return sc
}

func (sc *SyncCounter) Increment() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.counter++
}

func (sc *SyncCounter) Decrement() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.counter--

	if sc.counter == 0 {
		sc.WhenFinished()
	}
}
