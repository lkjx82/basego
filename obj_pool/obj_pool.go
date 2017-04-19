package obj_pool

import (
	"container/list"
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------
// obj pool interface
type ObjPool interface {
	Alloc() interface{}
	Free(interface{})
	//	Dump() int
}

// -----------------------------------------------------------------------------
// list node for pool
type objPoolLstNode struct {
	t   time.Time
	obj interface{}
	//	block []byte
}

// -----------------------------------------------------------------------------
// real object pool
type objPoolImp struct {
	AllocC chan interface{}
	freeC  chan interface{}
}

// -----------------------------------------------------------------------------
// alloc a object
func (this *objPoolImp) Alloc() interface{} {
	b := <-this.AllocC
	return b
}

// -----------------------------------------------------------------------------
// free a object to pool
func (this *objPoolImp) Free(v interface{}) {
	this.freeC <- v
}

// -----------------------------------------------------------------------------
// ObjPool call this func new a object, and cache it to the pool
type ObjPoolNewFun func() interface{}

// -----------------------------------------------------------------------------
// create a obj pool,  newFun: the object construct function
func NewObjPool(newFun ObjPoolNewFun) ObjPool {
	rbi := objPoolImp{}
	rbi.AllocC = make(chan interface{})
	rbi.freeC = make(chan interface{})

	go func() {
		lst := list.New()

		for {
			// 先往队列里放1个
			if lst.Len() == 0 {
				lst.PushBack(objPoolLstNode{t: time.Now(), obj: newFun()})
			}

			e := lst.Front()
			// 超时检测器
			timeout := time.NewTimer(time.Minute)

			select {
			// 外部有人还了，就放回到list里去
			case obj := <-rbi.freeC:
				timeout.Stop()
				lst.PushFront(objPoolLstNode{t: time.Now(), obj: obj})

				fmt.Println(lst.Len())

			// 先给get 里推一个，等待外部来取
			case rbi.AllocC <- e.Value.(objPoolLstNode).obj:
				timeout.Stop()
				lst.Remove(e)

				fmt.Println(lst.Len())

			// 1分钟，没有操作过了，就把1分钟以前的干掉，释放内存
			case <-timeout.C:
				e := lst.Front()
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(objPoolLstNode).t) > time.Minute {
						lst.Remove(e)
						e.Value = nil
					}
					e = n
				}

				fmt.Println(lst.Len())
			}
		}
	}()

	return &rbi
}

// -----------------------------------------------------------------------------
