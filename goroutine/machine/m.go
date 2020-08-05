package machine

import (
	"go_concept_exercise/goroutine/define"
)

func M() {
	for {
		define.SchedObject.Lock.Lock()    //互斥地从就绪G队列中取一个g出来运行
		if len(define.SchedObject.Allg) > 0 {
			g := define.SchedObject.Allg[0]
			define.SchedObject.Allg = define.SchedObject.Allg[1:]
			define.SchedObject.Lock.Unlock()
			g.Run()        //运行它
		} else {
			define.SchedObject.Lock.Unlock()
		}
	}
}
