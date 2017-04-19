package obj_pool

// ----------------------------------------------------------------------------
// 线程安全的对象池,池大小定长

type objPoolFixedImp struct {
	objs     chan interface{}
	newFunc  ObjPoolNewFun
	overSize bool // 超出默认size，是否继续new, false: 返回nil
}

// ----------------------------------------------------------------------------

func NewObjPoolFixed(size int, newFun ObjPoolNewFun, overSize bool) ObjPool {
	pool := objPoolFixedImp{}
	pool.objs = make(chan interface{}, size)
	pool.newFunc = newFun
	pool.overSize = overSize
	if !overSize {
		for i := 0; i < size; i++ {
			pool.objs <- newFun()
		}
	}
	return &pool
}

// ----------------------------------------------------------------------------

func (this *objPoolFixedImp) Alloc() interface{} {
	var obj interface{} = nil
	select {
	case obj = <-this.objs:
	default:
		if this.overSize {
			obj = this.newFunc()
		}
	}
	return obj
}

// ----------------------------------------------------------------------------

func (this *objPoolFixedImp) Free(obj interface{}) {
	select {
	case this.objs <- obj:
	default:
	}
}

// ----------------------------------------------------------------------------
