package id_test

import (
	"container/list"

	"runtime"
	"sync"
	"testing"

	baseId "github.com/lkjx82/basego/id"
)

// -----------------------------------------------------------------------------

func TestIdGen(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	idMap := make(map[int64]int64)

	lst := list.New()

	var m sync.Mutex
	var g sync.WaitGroup

	svrCnt := int16(1025)
	idGenCntPerSvr := 30000

	g.Add(int(svrCnt))

	for i := int16(0); i < svrCnt; i++ {
		ig := baseId.NewIdGen(i)

		if ig == nil {
			if i < 1024 {
				t.Fail()
			}
			g.Done()
			continue
		}

		tmpI := i
		go func() {
			for j := 0; j < idGenCntPerSvr; j++ {
				id := ig.GenId()
				m.Lock()
				lst.PushBack(id)
				m.Unlock()

				if tmpI != baseId.Id2ServId(id) {
					t.Logf("id:%d, i:%d", id, tmpI)
					t.Fail()
				}
			}
			g.Done()
		}()
	}

	t.Log("all done")
	g.Wait()

	for lst.Front() != nil {
		e := lst.Front()
		id := e.Value.(int64)
		if _, has := idMap[id]; has {
			t.Log(id)
			t.Fail()
		} else {
			idMap[id] = id
		}
		lst.Remove(e)
	}
}

// -----------------------------------------------------------------------------

func TestIdGenUnsafe(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
	idMap := make(map[int64]int64)

	lst := list.New()

	svrCnt := int16(1025)
	idGenCntPerSvr := 30000

	for i := int16(0); i < svrCnt; i++ {
		ig := baseId.NewIdGen(i)

		if ig == nil {
			if i < 1024 {
				t.Fail()
			}
			continue
		}

		tmpI := i
		for j := 0; j < idGenCntPerSvr; j++ {
			id := ig.GenIdUnsafe()
			lst.PushBack(id)

			if tmpI != baseId.Id2ServId(id) {
				t.Logf("id:%d, i:%d", id, tmpI)
				t.Fail()
			}
		}
	}

	for lst.Front() != nil {
		e := lst.Front()
		id := e.Value.(int64)
		if _, has := idMap[id]; has {
			t.Fail()
		} else {
			idMap[id] = id
		}
		lst.Remove(e)
	}
}

// -----------------------------------------------------------------------------

func BenchmarkIdGen(b *testing.B) {
	g := sync.WaitGroup{}
	g.Add(b.N)
	idg := baseId.NewIdGen(int16(1))
	for i := 0; i < b.N; i++ {
		go func() {
			idg.GenId()
			g.Done()
		}()
	}
	g.Wait()
}

// -----------------------------------------------------------------------------
