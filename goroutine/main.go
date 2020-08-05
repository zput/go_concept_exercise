package main

import (
	"go_concept_exercise/goroutine/define"
	"go_concept_exercise/goroutine/goroutine"
	"go_concept_exercise/goroutine/machine"
	"time"
)

const GOMAXPROCS = 2

func main(){
	for i:=0; i<GOMAXPROCS; i++ {
		go machine.M()
	}

	for {

		goroutineObjectTemp := new(goroutine.SelfGoroutine)

		define.SchedObject.Lock.Lock()
		define.SchedObject.Allg = append(define.SchedObject.Allg, goroutineObjectTemp)

		define.SchedObject.Lock.Unlock()

		time.Sleep(time.Second)
	}

}
