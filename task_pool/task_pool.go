// -----------------------------------------------------------------------------

package task_pool

// -----------------------------------------------------------------------------

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// -----------------------------------------------------------------------------

type Task struct {
	ch   chan []interface{}                // 接受 notify 参数的chan
	q    chan int                          // 结束的chan
	name string                            // 名字
	f    func(task *Task, v []interface{}) // task 的执行函数
	ref  int32                             // 读写协程的引用基数，为0表示读写协程都关闭了。此时可以释放
}

// -----------------------------------------------------------------------------
// task 结束，退出
func (this *Task) Quit() {
	fmt.Println("quit task ~~~~")
	if atomic.LoadInt32(&this.ref) != 0 {
		this.q <- 1
	}
}

// -----------------------------------------------------------------------------
// 传递参数给task执行，激活task, task 包装的function就会执行
func (this *Task) Notify(v ...interface{}) {
	if atomic.LoadInt32(&this.ref) != 0 {
		this.ch <- v
	}
}

// -----------------------------------------------------------------------------
// 获取 task 的name
func (this *Task) GetName() string {
	return this.name
}

// -----------------------------------------------------------------------------
// 新建一个task， task的名字和函数原型
func (this *TaskPool) NewTask(name string, f func(task *Task, v []interface{})) *Task {
	t := &Task{}
	t.f = f
	t.name = name

	if this.Add(t) == false {
		return nil
	}

	t.ch = make(chan []interface{}, 10)
	t.q = make(chan int, 5)

	atomic.StoreInt32(&t.ref, 1)

	go func() {
		defer this.Del(t.name)
		for {
			select {
			case v := <-t.ch:
				t.f(t, v)
			case <-t.q:
				fmt.Println("task quit")
				atomic.StoreInt32(&t.ref, 0)
				close(t.ch)
				close(t.q)
				return
			}
		}
	}()

	return t
}

// -----------------------------------------------------------------------------
// task pool
type TaskPool struct {
	tasks map[string]*Task
	lock  sync.Mutex
}

// -----------------------------------------------------------------------------

func NewTaskPool() *TaskPool {
	tp := &TaskPool{}
	tp.tasks = make(map[string]*Task)
	return tp
}

// -----------------------------------------------------------------------------

func (this *TaskPool) Add(t *Task) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.tasks[t.name]; ok {
		return false
	}
	this.tasks[t.name] = t
	return true
}

// -----------------------------------------------------------------------------

func (this *TaskPool) Del(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.tasks, name)
}

// -----------------------------------------------------------------------------

func (this *TaskPool) Get(name string) *Task {
	this.lock.Lock()
	defer this.lock.Unlock()
	if t, ok := this.tasks[name]; ok {
		return t
	}
	return nil
}

// -----------------------------------------------------------------------------
// 默认的全局 task 池
var gTp *TaskPool = NewTaskPool()

// -----------------------------------------------------------------------------

func NewTask(name string, f func(task *Task, v []interface{})) *Task {
	return gTp.NewTask(name, f)
}

// -----------------------------------------------------------------------------

func AddTask(t *Task) bool {
	return gTp.Add(t)
}

// -----------------------------------------------------------------------------

func DelTask(name string) *Task {
	t := gTp.Get(name)
	gTp.Del(name)
	return t
}

// -----------------------------------------------------------------------------

func GetTask(name string) *Task {
	return gTp.Get(name)
}

// -----------------------------------------------------------------------------

func NotifyTask(name string, v ...interface{}) bool {
	fmt.Println("Notify Task ", v)
	if t := gTp.Get(name); t != nil {
		t.Notify(v...)
		return true
	}
	return false
}

// -----------------------------------------------------------------------------
