package define

import "sync"

type G interface {
	Run()
}

type Sched struct {
	Allg  []G
	Lock    *sync.Mutex
}

var SchedObject = Sched{
	Lock: new(sync.Mutex),
}
