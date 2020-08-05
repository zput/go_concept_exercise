package goroutine

import (
	"fmt"
)

type SelfGoroutine struct{
}

func (this *SelfGoroutine)Run(){
	fmt.Println("working now")
}
