package task_pool_test

import (
	"runtime"
	"testing"
	"time"

	baseTaskPool "github.com/lkjx82/basego/task_pool"

	"fmt"
)

// -----------------------------------------------------------------------------

func recvFunc(t *baseTaskPool.Task, v []interface{}) {
	fmt.Println(t.GetName())
	for _, va := range v {
		fmt.Println(va)
	}
}

// -----------------------------------------------------------------------------

func TestTaskPoolFixed(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	fmt.Println("task pool test...")
	tp := baseTaskPool.NewTaskPool()
	task := tp.NewTask("recvTask ", recvFunc)
	task.Notify("what th fuck", 1, "你妹")
	<-time.After(time.Second)
	task.Quit()
	<-time.After(time.Second)
}
